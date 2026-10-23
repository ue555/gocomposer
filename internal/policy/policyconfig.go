// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Policy/PolicyConfig.php.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"fmt"
	"strings"

	"gocomposer/internal/phpjson"
)

// PolicyConfig::RESERVED_NAMES, BUILTIN_LIST_NAMES, FUTURE_RESERVED_*, NON_LIST_KEYS.
var (
	ReservedNames             = []string{AdvisoriesName, AbandonedName, CooldownName}
	BuiltinListNames          = []string{AdvisoriesName, AbandonedName, MalwareName, CooldownName}
	FutureReservedPrefixes    = []string{"ignore"}
	FutureReservedNames       = []string{"package", "packages", "license", "licence", "licenses", "licences", "support", "maintenance", "security", "minimum-release-age"}
	NonListKeys               = []string{"ignore-unreachable"}
	builtinAndNonListKeyNames = append(append([]string(nil), BuiltinListNames...), NonListKeys...)
)

// ConfigSource is the part of Composer\Config used here (*config.Config).
type ConfigSource interface {
	Get(key string) phpjson.Value
}

// PolicyConfig is Composer\Policy\PolicyConfig.
type PolicyConfig struct {
	Enabled           bool
	Advisories        *AdvisoriesPolicyConfig
	Malware           *MalwarePolicyConfig
	Abandoned         *AbandonedPolicyConfig
	Cooldown          *CooldownPolicyConfig
	CustomLists       *OMap[*CustomListPolicyConfig]
	IgnoreUnreachable IgnoreUnreachable
}

// NewPolicyConfig mirrors the constructor.
func NewPolicyConfig(enabled bool, advisories *AdvisoriesPolicyConfig, malware *MalwarePolicyConfig, abandoned *AbandonedPolicyConfig, cooldown *CooldownPolicyConfig, customLists *OMap[*CustomListPolicyConfig], ignoreUnreachable IgnoreUnreachable) *PolicyConfig {
	if customLists == nil {
		customLists = NewOMap[*CustomListPolicyConfig]()
	}
	return &PolicyConfig{enabled, advisories, malware, abandoned, cooldown, customLists, ignoreUnreachable}
}

// GetFutureReservedListNameError is PolicyConfig::getFutureReservedListNameError ("" = none).
func GetFutureReservedListNameError(listName string) string {
	for _, prefix := range FutureReservedPrefixes {
		if strings.HasPrefix(listName, prefix) {
			return fmt.Sprintf(`"%s" starts with reserved prefix "%s".`, listName, prefix)
		}
	}
	if inStrings(listName, FutureReservedNames) {
		return fmt.Sprintf(`"%s" is reserved for future use.`, listName)
	}
	return ""
}

func assertCustomListNameAllowed(listName string) error {
	if inStrings(listName, ReservedNames) {
		return newErr(UnexpectedValueException, fmt.Sprintf(`Invalid custom dependency policy name "%s": this name is reserved for a built-in dependency policy.`, listName))
	}
	if e := GetFutureReservedListNameError(listName); e != "" {
		return newErr(UnexpectedValueException, "Invalid custom dependency policy name: "+e)
	}
	return nil
}

// FromConfig is PolicyConfig::fromConfig: reads config.policy with BC
// fallback to config.audit and applies the COMPOSER_POLICY_* /
// COMPOSER_SECURITY_BLOCKING_ABANDONED / COMPOSER_AUDIT_ABANDONED /
// COMPOSER_POLICY_COOLDOWN_PERIOD overrides.
func FromConfig(config ConfigSource) (*PolicyConfig, error) {
	policyRaw := config.Get("policy")
	auditRaw := config.Get("audit")

	if policyRaw == false {
		return NewPolicyConfig(false, AdvisoriesPolicyConfigDisabled(), MalwarePolicyConfigDisabled(), AbandonedPolicyConfigDisabled(), CooldownPolicyConfigDisabled(), nil, IgnoreUnreachableAll()), nil
	}
	var policyConfig phpjson.Value = phpjson.NewMap()
	if phpjson.IsArray(policyRaw) {
		policyConfig = policyRaw
	}
	var auditConfig phpjson.Value = phpjson.NewMap()
	if phpjson.IsArray(auditRaw) {
		auditConfig = auditRaw
	}

	advisories, err := AdvisoriesPolicyConfigFromRawConfig(policyConfig, auditConfig)
	if err != nil {
		return nil, err
	}
	malware, err := MalwarePolicyConfigFromRawConfig(policyConfig)
	if err != nil {
		return nil, err
	}
	abandoned, err := AbandonedPolicyConfigFromRawConfig(policyConfig, auditConfig)
	if err != nil {
		return nil, err
	}
	cooldown, err := CooldownPolicyConfigFromRawConfig(policyConfig)
	if err != nil {
		return nil, err
	}

	customLists := NewOMap[*CustomListPolicyConfig]()
	eachArr(policyConfig, func(listName string, listConfig phpjson.Value) {
		if err != nil {
			return
		}
		if !isIntKey(listName) && inStrings(listName, builtinAndNonListKeyNames) {
			return
		}
		if err = assertCustomListNameAllowed(listName); err != nil {
			return
		}
		var c *CustomListPolicyConfig
		c, err = CustomListPolicyConfigFromRawConfig(listName, listConfig)
		if err == nil {
			customLists.Set(listName, c)
		}
	})
	if err != nil {
		return nil, err
	}

	ignoreUnreachable := IgnoreUnreachableDefault()
	if isset(policyConfig, "ignore-unreachable") {
		ignoreUnreachable = IgnoreUnreachableFromRawPolicyConfig(policyConfig)
	} else if isset(auditConfig, "ignore-unreachable") {
		ignoreUnreachable = IgnoreUnreachableFromRawAuditConfig(auditConfig)
	}

	if o, err := GetBoolEnv("COMPOSER_POLICY_ADVISORIES_BLOCK"); err != nil {
		return nil, err
	} else if o != nil {
		advisories = NewAdvisoriesPolicyConfig(*o, advisories.Audit, advisories.Ignore, advisories.IgnoreID, advisories.IgnoreSeverity)
	}

	if o, err := GetBoolEnv("COMPOSER_POLICY_MALWARE_BLOCK"); err != nil {
		return nil, err
	} else if o != nil {
		malware = NewMalwarePolicyConfig(*o, malware.Audit, malware.BlockScope, malware.Ignore, malware.IgnoreSource)
	}

	// PHP evaluates "getBoolEnv(canonical) ?? getBoolEnv(legacy)", so the legacy
	// variable is only read (and validated) when the canonical one is unset
	o, err := GetBoolEnv("COMPOSER_POLICY_ABANDONED_BLOCK")
	if err != nil {
		return nil, err
	}
	if o == nil {
		if o, err = GetBoolEnv("COMPOSER_SECURITY_BLOCKING_ABANDONED"); err != nil {
			return nil, err
		}
	}
	if o != nil {
		abandoned = NewAbandonedPolicyConfig(*o, abandoned.Audit, abandoned.Ignore)
	}

	if v, ok := getEnv("COMPOSER_AUDIT_ABANDONED"); ok {
		if !inStrings(v, Audits) {
			return nil, newErr(RuntimeException, "Invalid value for COMPOSER_AUDIT_ABANDONED: "+v+". Expected one of "+strings.Join(Audits, ", ")+".")
		}
		abandoned = NewAbandonedPolicyConfig(abandoned.Block, v, abandoned.Ignore)
	}

	if o, err := GetBoolEnv("COMPOSER_POLICY_COOLDOWN_BLOCK"); err != nil {
		return nil, err
	} else if o != nil {
		cooldown = NewCooldownPolicyConfig(*o, cooldown.Audit, cooldown.Ignore, cooldown.Period)
	}

	return NewPolicyConfig(true, advisories, malware, abandoned, cooldown, customLists, ignoreUnreachable), nil
}

// GetAllLists is PolicyConfig::getAllLists (built-in lists then custom lists).
func (c *PolicyConfig) GetAllLists() *OMap[ListPolicy] {
	all := NewOMap[ListPolicy]()
	all.Set(AdvisoriesName, c.Advisories)
	all.Set(MalwareName, c.Malware)
	all.Set(AbandonedName, c.Abandoned)
	all.Set(CooldownName, c.Cooldown)
	c.CustomLists.Each(func(k string, v *CustomListPolicyConfig) { all.Set(k, v) })
	return all
}

// filterableLists excludes advisories/abandoned/cooldown (dedicated filters).
func (c *PolicyConfig) filterableLists() *OMap[ListPolicy] {
	all := c.GetAllLists()
	all.Delete(AbandonedName)
	all.Delete(AdvisoriesName)
	all.Delete(CooldownName)
	return all
}

// GetActiveAuditFilterLists is PolicyConfig::getActiveAuditFilterLists.
func (c *PolicyConfig) GetActiveAuditFilterLists() *OMap[ListPolicy] {
	lists := NewOMap[ListPolicy]()
	c.filterableLists().Each(func(name string, l ListPolicy) {
		if l.AuditMode() != AuditIgnore {
			lists.Set(name, l)
		}
	})
	return lists
}

// GetActiveBlockFilterLists is PolicyConfig::getActiveBlockFilterLists.
func (c *PolicyConfig) GetActiveBlockFilterLists(blockScope string) *OMap[ListPolicy] {
	lists := NewOMap[ListPolicy]()
	c.filterableLists().Each(func(name string, l ListPolicy) {
		if l.ShouldBlock(blockScope) {
			lists.Set(name, l)
		}
	})
	return lists
}

func (c *PolicyConfig) GetActiveAuditFilterListNames() []string {
	return c.GetActiveAuditFilterLists().Keys()
}

func (c *PolicyConfig) GetActiveBlockFilterListNames(blockScope string) []string {
	return c.GetActiveBlockFilterLists(blockScope).Keys()
}

// GetCustomListsWithSources is PolicyConfig::getCustomListsWithSources.
func (c *PolicyConfig) GetCustomListsWithSources() *OMap[*CustomListPolicyConfig] {
	out := NewOMap[*CustomListPolicyConfig]()
	c.CustomLists.Each(func(k string, l *CustomListPolicyConfig) {
		if len(l.Sources) > 0 {
			out.Set(k, l)
		}
	})
	return out
}

// WithBlockingDisabled is PolicyConfig::withBlockingDisabled (--no-blocking,
// --no-security-blocking, COMPOSER_NO_SECURITY_BLOCKING, COMPOSER_NO_BLOCKING).
func (c *PolicyConfig) WithBlockingDisabled() *PolicyConfig {
	custom := NewOMap[*CustomListPolicyConfig]()
	c.CustomLists.Each(func(k string, l *CustomListPolicyConfig) { custom.Set(k, l.WithBlockingDisabled()) })
	return NewPolicyConfig(c.Enabled, c.Advisories.WithBlockingDisabled(), c.Malware.WithBlockingDisabled(), c.Abandoned.WithBlockingDisabled(), c.Cooldown.WithBlockingDisabled(), custom, c.IgnoreUnreachable)
}

// WithCooldownBlockingDisabled is PolicyConfig::withCooldownBlockingDisabled (--no-cooldown).
func (c *PolicyConfig) WithCooldownBlockingDisabled() *PolicyConfig {
	return NewPolicyConfig(c.Enabled, c.Advisories, c.Malware, c.Abandoned, c.Cooldown.WithBlockingDisabled(), c.CustomLists, c.IgnoreUnreachable)
}

// WithIgnoreUnreachable is PolicyConfig::withIgnoreUnreachable.
func (c *PolicyConfig) WithIgnoreUnreachable(scopes ...string) (*PolicyConfig, error) {
	iu, err := c.IgnoreUnreachable.With(scopes...)
	if err != nil {
		return nil, err
	}
	return NewPolicyConfig(c.Enabled, c.Advisories, c.Malware, c.Abandoned, c.Cooldown, c.CustomLists, iu), nil
}

// WithIgnoreSeverity is PolicyConfig::withIgnoreSeverity.
func (c *PolicyConfig) WithIgnoreSeverity(severities []string) *PolicyConfig {
	return NewPolicyConfig(c.Enabled, c.Advisories.WithIgnoreSeverity(severities), c.Malware, c.Abandoned, c.Cooldown, c.CustomLists, c.IgnoreUnreachable)
}

// WithAudit is PolicyConfig::withAudit; abandoned "" keeps the current value.
func (c *PolicyConfig) WithAudit(abandoned string) *PolicyConfig {
	a := c.Abandoned.Audit
	if abandoned != "" {
		a = abandoned
	}
	return NewPolicyConfig(c.Enabled, c.Advisories, c.Malware, c.Abandoned.WithAudit(a), c.Cooldown, c.CustomLists, c.IgnoreUnreachable)
}
