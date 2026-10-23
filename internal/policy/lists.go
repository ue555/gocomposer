// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Policy/ListPolicyConfig.php, AdvisoriesPolicyConfig.php,
// MalwarePolicyConfig.php, AbandonedPolicyConfig.php,
// CustomListPolicyConfig.php and
// src/Composer/FilterList/Source/SourceValidator.php, UrlSource.php.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"fmt"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// ListPolicyConfig::AUDIT_* and BLOCK_SCOPE_*.
const (
	AuditIgnore = "ignore"
	AuditReport = "report"
	AuditFail   = "fail"

	BlockScopeUpdate  = "update"
	BlockScopeInstall = "install"
	BlockScopeAll     = "all"
)

// Audits is ListPolicyConfig::AUDITS.
var Audits = []string{AuditIgnore, AuditReport, AuditFail}

// AvailableBlockScopes is ListPolicyConfig::AVAILABLE_BLOCK_SCOPES.
var AvailableBlockScopes = []string{BlockScopeUpdate, BlockScopeInstall, BlockScopeAll}

// Built-in list names.
const (
	AdvisoriesName = "advisories"
	MalwareName    = "malware"
	AbandonedName  = "abandoned"
	CooldownName   = "cooldown"
)

// ListPolicy is the abstract ListPolicyConfig.
type ListPolicy interface {
	ListName() string
	Blocks() bool
	AuditMode() string
	IgnoreRules() *IgnoreMap
	ShouldBlock(blockScope string) bool
	GetIgnoreForOperation(op string) *IgnoreMap
	GetFlatIgnoreForOperation(op string) *OMap[*string]
}

// ListBase holds the fields shared by every list (name/block/audit/ignore).
type ListBase struct {
	Name   string
	Block  bool
	Audit  string
	Ignore *IgnoreMap
}

func (l *ListBase) ListName() string        { return l.Name }
func (l *ListBase) Blocks() bool            { return l.Block }
func (l *ListBase) AuditMode() string       { return l.Audit }
func (l *ListBase) IgnoreRules() *IgnoreMap { return l.Ignore }

func (l *ListBase) baseShouldBlock(blockScope string, supportsInstall bool) bool {
	if !l.Block {
		return false
	}
	if blockScope == BlockScopeInstall && !supportsInstall {
		return false
	}
	return true
}

// GetIgnoreForOperation is ListPolicyConfig::getIgnoreForOperation.
func (l *ListBase) GetIgnoreForOperation(op string) *IgnoreMap {
	return FilterByOperation(l.Ignore, op)
}

// GetFlatIgnoreForOperation is ListPolicyConfig::getFlatIgnoreForOperation.
func (l *ListBase) GetFlatIgnoreForOperation(op string) *OMap[*string] {
	result := NewOMap[*string]()
	l.Ignore.Each(func(name string, rules []*IgnorePackageRule) {
		for _, r := range rules {
			if appliesTo(op, r.OnBlock, r.OnAudit) {
				result.Set(name, mergeReason(result, name, r.Reason))
			}
		}
	})
	return result
}

// mergeReason is ListPolicyConfig::mergeReason.
func mergeReason(m *OMap[*string], key string, newReason *string) *string {
	existing, ok := m.Get(key)
	if !ok {
		return newReason
	}
	if newReason == nil {
		return existing
	}
	if existing == nil {
		return newReason
	}
	if *existing == *newReason {
		return existing
	}
	for _, part := range strings.Split(*existing, ";") {
		if strings.Trim(part, " \t\n\r\x00\x0B") == *newReason {
			return existing
		}
	}
	s := *existing + "; " + *newReason
	return &s
}

func emptyIgnore() *IgnoreMap { return NewOMap[[]*IgnorePackageRule]() }

func orEmpty(m *IgnoreMap) *IgnoreMap {
	if m == nil {
		return emptyIgnore()
	}
	return m
}

func auditValue(v phpjson.Value, def string) string {
	if v == nil {
		return def
	}
	return toStr(v)
}

// legacySingle is ListPolicyConfig::parseLegacySingleIgnore.
func parseLegacySingleIgnore(key string, value phpjson.Value) (reason *string, onBlock, onAudit bool, err error) {
	onBlock, onAudit = true, true
	_, valueIsString := value.(string)
	switch {
	case isIntKey(key) && valueIsString:
	case valueIsString:
		reason = strPtr(value.(string))
	case phpjson.IsArray(value):
		apply := coalesce(value, "apply", "all")
		reason = nullableStr(coalesce(value, "reason", nil))
		as, isStr := apply.(string)
		if !isStr || (as != "audit" && as != "block" && as != "all") {
			shown := debugType(apply)
			if isStr {
				shown = as
			}
			return nil, false, false, newErr(InvalidArgumentException, fmt.Sprintf("Invalid 'apply' value for '%s': %s. Expected 'audit', 'block', or 'all'.", key, shown))
		}
		onBlock = as == "block" || as == "all"
		onAudit = as == "audit" || as == "all"
	}
	return reason, onBlock, onAudit, nil
}

// parseLegacyAuditIgnore is ListPolicyConfig::parseLegacyAuditIgnore.
func parseLegacyAuditIgnore(config phpjson.Value) (*IgnoreMap, *OMap[*IgnoreIdRule], error) {
	packages := emptyIgnore()
	ids := NewOMap[*IgnoreIdRule]()
	var err error
	eachArr(config, func(key string, value phpjson.Value) {
		if err != nil {
			return
		}
		id := key
		if isIntKey(key) {
			id = toStr(value)
		}
		isPackageName := strings.Contains(id, "/")
		reason, onBlock, onAudit, e := parseLegacySingleIgnore(key, value)
		if e != nil {
			err = e
			return
		}
		if isPackageName {
			appendTo(packages, id, NewIgnorePackageRule(id, &semver.MatchAll{}, reason, onBlock, onAudit))
		} else {
			ids.Set(id, NewIgnoreIdRule(id, reason, onBlock, onAudit))
		}
	})
	return packages, ids, err
}

// parseLegacyIgnoreWithApply is ListPolicyConfig::parseLegacyIgnoreWithApply.
func parseLegacyIgnoreWithApply(config phpjson.Value) (*IgnoreMap, error) {
	result := emptyIgnore()
	var err error
	eachArr(config, func(key string, value phpjson.Value) {
		if err != nil {
			return
		}
		name := key
		if isIntKey(key) {
			name = toStr(value)
		}
		reason, onBlock, onAudit, e := parseLegacySingleIgnore(key, value)
		if e != nil {
			err = e
			return
		}
		result.Set(name, []*IgnorePackageRule{NewIgnorePackageRule(name, &semver.MatchAll{}, reason, onBlock, onAudit)})
	})
	return result, err
}

// ---------------------------------------------------------------- advisories

// AdvisoriesPolicyConfig is Composer\Policy\AdvisoriesPolicyConfig.
type AdvisoriesPolicyConfig struct {
	ListBase
	IgnoreID       *OMap[*IgnoreIdRule]
	IgnoreSeverity *OMap[*IgnoreSeverityRule]
}

func NewAdvisoriesPolicyConfig(block bool, audit string, ignore *IgnoreMap, ignoreID *OMap[*IgnoreIdRule], ignoreSeverity *OMap[*IgnoreSeverityRule]) *AdvisoriesPolicyConfig {
	if ignoreID == nil {
		ignoreID = NewOMap[*IgnoreIdRule]()
	}
	if ignoreSeverity == nil {
		ignoreSeverity = NewOMap[*IgnoreSeverityRule]()
	}
	return &AdvisoriesPolicyConfig{ListBase{AdvisoriesName, block, audit, orEmpty(ignore)}, ignoreID, ignoreSeverity}
}

func (c *AdvisoriesPolicyConfig) ShouldBlock(scope string) bool {
	return c.baseShouldBlock(scope, false)
}

// GetIgnoreIdForOperation is AdvisoriesPolicyConfig::getIgnoreIdForOperation.
func (c *AdvisoriesPolicyConfig) GetIgnoreIdForOperation(op string) *OMap[*string] {
	result := NewOMap[*string]()
	c.IgnoreID.Each(func(id string, r *IgnoreIdRule) {
		if appliesTo(op, r.OnBlock, r.OnAudit) {
			result.Set(id, r.Reason)
		}
	})
	return result
}

// GetIgnoreListForOperation is AdvisoriesPolicyConfig::getIgnoreListForOperation.
func (c *AdvisoriesPolicyConfig) GetIgnoreListForOperation(op string) *OMap[*string] {
	result := c.GetIgnoreIdForOperation(op)
	c.GetFlatIgnoreForOperation(op).Each(func(name string, reason *string) {
		result.Set(name, mergeReason(result, name, reason))
	})
	return result
}

// GetIgnoreSeverityForOperation is AdvisoriesPolicyConfig::getIgnoreSeverityForOperation.
func (c *AdvisoriesPolicyConfig) GetIgnoreSeverityForOperation(op string) *OMap[*string] {
	result := NewOMap[*string]()
	c.IgnoreSeverity.Each(func(sev string, r *IgnoreSeverityRule) {
		if appliesTo(op, r.OnBlock, r.OnAudit) {
			result.Set(sev, r.Reason)
		}
	})
	return result
}

func (c *AdvisoriesPolicyConfig) WithBlockingDisabled() *AdvisoriesPolicyConfig {
	return NewAdvisoriesPolicyConfig(false, c.Audit, c.Ignore, c.IgnoreID, c.IgnoreSeverity)
}

func (c *AdvisoriesPolicyConfig) WithAudit(audit string) *AdvisoriesPolicyConfig {
	return NewAdvisoriesPolicyConfig(c.Block, audit, c.Ignore, c.IgnoreID, c.IgnoreSeverity)
}

// WithIgnoreSeverity is AdvisoriesPolicyConfig::withIgnoreSeverity.
func (c *AdvisoriesPolicyConfig) WithIgnoreSeverity(severities []string) *AdvisoriesPolicyConfig {
	sev := c.IgnoreSeverity.Clone()
	for _, s := range severities {
		if r, ok := sev.Get(s); !ok || r == nil {
			sev.Set(s, NewIgnoreSeverityRule(s, nil, false, true))
		}
	}
	return NewAdvisoriesPolicyConfig(c.Block, c.Audit, c.Ignore, c.IgnoreID, sev)
}

// AdvisoriesPolicyConfigFromRawConfig is AdvisoriesPolicyConfig::fromRawConfig.
func AdvisoriesPolicyConfigFromRawConfig(policyConfig, auditConfig phpjson.Value) (*AdvisoriesPolicyConfig, error) {
	if !isset(policyConfig, AdvisoriesName) && phpjson.ArrayLen(auditConfig) != 0 {
		pkgs, ids, err := parseLegacyAuditIgnore(coalesce(auditConfig, "ignore", nil))
		if err != nil {
			return nil, err
		}
		sev, err := parseLegacySeverityWithApply(coalesce(auditConfig, "ignore-severity", nil))
		if err != nil {
			return nil, err
		}
		return NewAdvisoriesPolicyConfig(toBool(coalesce(auditConfig, "block-insecure", true)), AuditFail, pkgs, ids, sev), nil
	}
	advisoryConfig := coalesce(policyConfig, AdvisoriesName, phpjson.NewMap())
	if advisoryConfig == false {
		return AdvisoriesPolicyConfigDisabled(), nil
	}
	if !phpjson.IsArray(advisoryConfig) {
		advisoryConfig = phpjson.NewMap()
	}
	ignore, err := ParseIgnoreMap(coalesce(advisoryConfig, "ignore", nil))
	if err != nil {
		return nil, err
	}
	ignoreID, err := ParseIgnoreIdMap(coalesce(advisoryConfig, "ignore-id", nil))
	if err != nil {
		return nil, err
	}
	ignoreSev, err := ParseIgnoreSeverityMap(coalesce(advisoryConfig, "ignore-severity", nil))
	if err != nil {
		return nil, err
	}
	return NewAdvisoriesPolicyConfig(toBool(coalesce(advisoryConfig, "block", true)), auditValue(coalesce(advisoryConfig, "audit", nil), AuditFail), ignore, ignoreID, ignoreSev), nil
}

// AdvisoriesPolicyConfigDisabled is AdvisoriesPolicyConfig::disabled.
func AdvisoriesPolicyConfigDisabled() *AdvisoriesPolicyConfig {
	return NewAdvisoriesPolicyConfig(false, AuditIgnore, nil, nil, nil)
}

func parseLegacySeverityWithApply(config phpjson.Value) (*OMap[*IgnoreSeverityRule], error) {
	result := NewOMap[*IgnoreSeverityRule]()
	var err error
	eachArr(config, func(key string, value phpjson.Value) {
		if err != nil {
			return
		}
		sev := key
		if isIntKey(key) {
			sev = toStr(value)
		}
		reason, onBlock, onAudit, e := parseLegacySingleIgnore(key, value)
		if e != nil {
			err = e
			return
		}
		result.Set(sev, NewIgnoreSeverityRule(sev, reason, onBlock, onAudit))
	})
	return result, err
}

// ------------------------------------------------------------------- malware

// MalwarePolicyConfig is Composer\Policy\MalwarePolicyConfig.
type MalwarePolicyConfig struct {
	ListBase
	BlockScope   string
	IgnoreSource []string
}

func NewMalwarePolicyConfig(block bool, audit, blockScope string, ignore *IgnoreMap, ignoreSource []string) *MalwarePolicyConfig {
	return &MalwarePolicyConfig{ListBase{MalwareName, block, audit, orEmpty(ignore)}, blockScope, ignoreSource}
}

// ShouldBlock is MalwarePolicyConfig::shouldBlock (supports install scope).
func (c *MalwarePolicyConfig) ShouldBlock(scope string) bool {
	if !c.baseShouldBlock(scope, true) {
		return false
	}
	return c.BlockScope == BlockScopeAll || c.BlockScope == scope
}

func (c *MalwarePolicyConfig) WithBlockingDisabled() *MalwarePolicyConfig {
	return NewMalwarePolicyConfig(false, c.Audit, c.BlockScope, c.Ignore, c.IgnoreSource)
}

func (c *MalwarePolicyConfig) WithAudit(audit string) *MalwarePolicyConfig {
	return NewMalwarePolicyConfig(c.Block, audit, c.BlockScope, c.Ignore, c.IgnoreSource)
}

// MalwarePolicyConfigFromRawConfig is MalwarePolicyConfig::fromRawConfig.
func MalwarePolicyConfigFromRawConfig(policyConfig phpjson.Value) (*MalwarePolicyConfig, error) {
	malwareConfig := coalesce(policyConfig, MalwareName, phpjson.NewMap())
	if malwareConfig == false {
		return MalwarePolicyConfigDisabled(), nil
	}
	if !phpjson.IsArray(malwareConfig) {
		malwareConfig = phpjson.NewMap()
	}
	ignore, err := ParseIgnoreMap(coalesce(malwareConfig, "ignore", nil))
	if err != nil {
		return nil, err
	}
	var ignoreSource []string
	for _, v := range arrayValues(coalesce(malwareConfig, "ignore-source", nil)) {
		ignoreSource = append(ignoreSource, toStr(v))
	}
	return NewMalwarePolicyConfig(toBool(coalesce(malwareConfig, "block", true)), auditValue(coalesce(malwareConfig, "audit", nil), AuditFail),
		auditValue(coalesce(malwareConfig, "block-scope", nil), BlockScopeAll), ignore, ignoreSource), nil
}

func MalwarePolicyConfigDisabled() *MalwarePolicyConfig {
	return NewMalwarePolicyConfig(false, AuditIgnore, BlockScopeAll, nil, nil)
}

// ----------------------------------------------------------------- abandoned

// AbandonedPolicyConfig is Composer\Policy\AbandonedPolicyConfig.
type AbandonedPolicyConfig struct {
	ListBase
}

func NewAbandonedPolicyConfig(block bool, audit string, ignore *IgnoreMap) *AbandonedPolicyConfig {
	return &AbandonedPolicyConfig{ListBase{AbandonedName, block, audit, orEmpty(ignore)}}
}

func (c *AbandonedPolicyConfig) ShouldBlock(scope string) bool {
	return c.baseShouldBlock(scope, false)
}

func (c *AbandonedPolicyConfig) WithBlockingDisabled() *AbandonedPolicyConfig {
	return NewAbandonedPolicyConfig(false, c.Audit, c.Ignore)
}

func (c *AbandonedPolicyConfig) WithAudit(audit string) *AbandonedPolicyConfig {
	return NewAbandonedPolicyConfig(c.Block, audit, c.Ignore)
}

// AbandonedPolicyConfigFromRawConfig is AbandonedPolicyConfig::fromRawConfig.
func AbandonedPolicyConfigFromRawConfig(policyConfig, auditConfig phpjson.Value) (*AbandonedPolicyConfig, error) {
	if !isset(policyConfig, AbandonedName) && phpjson.ArrayLen(auditConfig) != 0 {
		ignore, err := parseLegacyIgnoreWithApply(coalesce(auditConfig, "ignore-abandoned", nil))
		if err != nil {
			return nil, err
		}
		return NewAbandonedPolicyConfig(toBool(coalesce(auditConfig, "block-abandoned", false)), auditValue(coalesce(auditConfig, "abandoned", nil), AuditFail), ignore), nil
	}
	abandonedConfig := coalesce(policyConfig, AbandonedName, phpjson.NewMap())
	if abandonedConfig == false {
		return AbandonedPolicyConfigDisabled(), nil
	}
	if !phpjson.IsArray(abandonedConfig) {
		abandonedConfig = phpjson.NewMap()
	}
	ignore, err := ParseIgnoreMap(coalesce(abandonedConfig, "ignore", nil))
	if err != nil {
		return nil, err
	}
	return NewAbandonedPolicyConfig(toBool(coalesce(abandonedConfig, "block", false)), auditValue(coalesce(abandonedConfig, "audit", nil), AuditFail), ignore), nil
}

func AbandonedPolicyConfigDisabled() *AbandonedPolicyConfig {
	return NewAbandonedPolicyConfig(false, AuditIgnore, nil)
}

// --------------------------------------------------------------- custom list

// UrlSource is Composer\FilterList\Source\UrlSource.
type UrlSource struct {
	ListName string
	URL      string
}

// ValidateSource is SourceValidator::validate.
func ValidateSource(listName string, source phpjson.Value) (*UrlSource, error) {
	if !isset(source, "type") {
		return nil, newErr(RuntimeException, `Source configuration is missing the "type" field.`)
	}
	typ := coalesce(source, "type", nil)
	if typ == "url" {
		u, ok := coalesce(source, "url", nil).(string)
		if !ok {
			return nil, newErr(RuntimeException, `Source configuration is missing a string "url" field.`)
		}
		if !strings.HasPrefix(u, "https://") {
			return nil, newErr(RuntimeException, `Source URL for policy list "`+listName+`" must start with "https://"; got "`+u+`".`)
		}
		return &UrlSource{ListName: listName, URL: u}, nil
	}
	return nil, newErr(RuntimeException, `Unsupported source type "`+toStr(typ)+`". Only "url" is currently supported.`)
}

// CustomListPolicyConfig is Composer\Policy\CustomListPolicyConfig.
type CustomListPolicyConfig struct {
	ListBase
	BlockScope string
	Sources    []*UrlSource
}

func NewCustomListPolicyConfig(name string, block bool, audit, blockScope string, ignore *IgnoreMap, sources []*UrlSource) *CustomListPolicyConfig {
	return &CustomListPolicyConfig{ListBase{name, block, audit, orEmpty(ignore)}, blockScope, sources}
}

// ShouldBlock is CustomListPolicyConfig::shouldBlock (supports install scope).
func (c *CustomListPolicyConfig) ShouldBlock(scope string) bool {
	if !c.baseShouldBlock(scope, true) {
		return false
	}
	return c.BlockScope == BlockScopeAll || c.BlockScope == scope
}

func (c *CustomListPolicyConfig) WithBlockingDisabled() *CustomListPolicyConfig {
	return NewCustomListPolicyConfig(c.Name, false, c.Audit, c.BlockScope, c.Ignore, c.Sources)
}

func (c *CustomListPolicyConfig) WithAudit(audit string) *CustomListPolicyConfig {
	return NewCustomListPolicyConfig(c.Name, c.Block, audit, c.BlockScope, c.Ignore, c.Sources)
}

// CustomListPolicyConfigFromRawConfig is CustomListPolicyConfig::fromRawConfig.
func CustomListPolicyConfigFromRawConfig(listName string, listConfig phpjson.Value) (*CustomListPolicyConfig, error) {
	if listConfig == false {
		return CustomListPolicyConfigDisabled(listName), nil
	}
	if listConfig == true {
		listConfig = phpjson.NewMap()
	}
	if !phpjson.IsArray(listConfig) {
		return CustomListPolicyConfigDisabled(listName), nil
	}
	var sources []*UrlSource
	var err error
	eachArr(coalesce(listConfig, "sources", nil), func(_ string, sc phpjson.Value) {
		if err != nil || !phpjson.IsArray(sc) {
			return
		}
		s, e := ValidateSource(listName, sc)
		if e != nil {
			err = e
			return
		}
		sources = append(sources, s)
	})
	if err != nil {
		return nil, err
	}
	bs := coalesce(listConfig, "block-scope", BlockScopeUpdate)
	bss, isStr := bs.(string)
	if !isStr || !inStrings(bss, AvailableBlockScopes) {
		enc, _ := phpjson.Encode(bs, 0)
		return nil, newErr(UnexpectedValueException, fmt.Sprintf(`Invalid block-scope for policy list "%s": expected one of "%s", got %s.`, listName, strings.Join(AvailableBlockScopes, `", "`), enc))
	}
	ignore, err := ParseIgnoreMap(coalesce(listConfig, "ignore", nil))
	if err != nil {
		return nil, err
	}
	return NewCustomListPolicyConfig(listName, toBool(coalesce(listConfig, "block", true)), auditValue(coalesce(listConfig, "audit", nil), AuditFail), bss, ignore, sources), nil
}

func CustomListPolicyConfigDisabled(listName string) *CustomListPolicyConfig {
	return NewCustomListPolicyConfig(listName, false, AuditIgnore, BlockScopeUpdate, nil, nil)
}
