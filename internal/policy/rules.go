// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Policy/IgnoreIdRule.php, IgnorePackageRule.php,
// IgnoreSeverityRule.php, IgnoreUnreachable.php.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package policy

import (
	"fmt"
	"regexp"
	"strings"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// Operations accepted by the *ForOperation helpers.
const (
	OperationBlock = "block"
	OperationAudit = "audit"
)

func appliesTo(op string, onBlock, onAudit bool) bool {
	return (op == OperationBlock && onBlock) || (op == OperationAudit && onAudit)
}

// IgnoreIdRule is Composer\Policy\IgnoreIdRule.
type IgnoreIdRule struct {
	ID      string
	Reason  *string
	OnBlock bool
	OnAudit bool
}

// NewIgnoreIdRule mirrors the constructor defaults (reason null, both true).
func NewIgnoreIdRule(id string, reason *string, onBlock, onAudit bool) *IgnoreIdRule {
	return &IgnoreIdRule{ID: id, Reason: reason, OnBlock: onBlock, OnAudit: onAudit}
}

// ParseIgnoreIdMap is IgnoreIdRule::parseIgnoreIdMap.
func ParseIgnoreIdMap(config phpjson.Value) (*OMap[*IgnoreIdRule], error) {
	rules := NewOMap[*IgnoreIdRule]()
	var err error
	eachArr(config, func(key string, value phpjson.Value) {
		if err != nil {
			return
		}
		if isIntKey(key) {
			s, ok := value.(string)
			if !ok {
				err = newErr(UnexpectedValueException, fmt.Sprintf("Invalid ignore-id entry at index %s: expected an advisory ID string, got %s.", key, debugType(value)))
				return
			}
			rules.Set(s, NewIgnoreIdRule(s, nil, true, true))
			return
		}
		switch t := value.(type) {
		case nil:
			rules.Set(key, NewIgnoreIdRule(key, nil, true, true))
		case string:
			rules.Set(key, NewIgnoreIdRule(key, strPtr(t), true, true))
		case []phpjson.Value, *phpjson.Map:
			rules.Set(key, NewIgnoreIdRule(key, nullableStr(coalesce(t, "reason", nil)), toBool(coalesce(t, "on-block", true)), toBool(coalesce(t, "on-audit", true))))
		default:
			err = newErr(UnexpectedValueException, fmt.Sprintf("Invalid ignore-id entry for \"%s\": value of type %s is not a supported shape. Expected null, a reason string, or a rule object.", key, debugType(value)))
		}
	})
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// IgnoreSeverityRule is Composer\Policy\IgnoreSeverityRule.
type IgnoreSeverityRule struct {
	Severity string
	Reason   *string
	OnBlock  bool
	OnAudit  bool
}

func NewIgnoreSeverityRule(severity string, reason *string, onBlock, onAudit bool) *IgnoreSeverityRule {
	return &IgnoreSeverityRule{Severity: severity, Reason: reason, OnBlock: onBlock, OnAudit: onAudit}
}

// ParseIgnoreSeverityMap is IgnoreSeverityRule::parseIgnoreSeverityMap.
func ParseIgnoreSeverityMap(config phpjson.Value) (*OMap[*IgnoreSeverityRule], error) {
	rules := NewOMap[*IgnoreSeverityRule]()
	var err error
	eachArr(config, func(key string, value phpjson.Value) {
		if err != nil {
			return
		}
		if isIntKey(key) {
			s, ok := value.(string)
			if !ok {
				err = newErr(UnexpectedValueException, fmt.Sprintf("Invalid ignore-severity entry at index %s: expected a severity string, got %s.", key, debugType(value)))
				return
			}
			rules.Set(s, NewIgnoreSeverityRule(s, nil, true, true))
			return
		}
		switch t := value.(type) {
		case nil:
			rules.Set(key, NewIgnoreSeverityRule(key, nil, true, true))
		case string:
			rules.Set(key, NewIgnoreSeverityRule(key, strPtr(t), true, true))
		case []phpjson.Value, *phpjson.Map:
			rules.Set(key, NewIgnoreSeverityRule(key, nullableStr(coalesce(t, "reason", nil)), toBool(coalesce(t, "on-block", true)), toBool(coalesce(t, "on-audit", true))))
		default:
			err = newErr(UnexpectedValueException, fmt.Sprintf("Invalid ignore-severity entry for \"%s\": value of type %s is not a supported shape. Expected null, a reason string, or a rule object.", key, debugType(value)))
		}
	})
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// IgnorePackageRule is Composer\Policy\IgnorePackageRule.
type IgnorePackageRule struct {
	PackageName      string
	Constraint       semver.Constraint
	Reason           *string
	OnBlock          bool
	OnAudit          bool
	PackageNameRegex *regexp.Regexp
}

// NewIgnorePackageRule mirrors the constructor.
func NewIgnorePackageRule(packageName string, constraint semver.Constraint, reason *string, onBlock, onAudit bool) *IgnorePackageRule {
	return &IgnorePackageRule{PackageName: packageName, Constraint: constraint, Reason: reason, OnBlock: onBlock, OnAudit: onAudit, PackageNameRegex: PackageNameToRegexp(packageName)}
}

// IgnoreMap is array<string, list<IgnorePackageRule>> keyed by package name pattern.
type IgnoreMap = OMap[[]*IgnorePackageRule]

// ParseIgnoreMap is IgnorePackageRule::parseIgnoreMap.
func ParseIgnoreMap(config phpjson.Value) (*IgnoreMap, error) {
	rules := NewOMap[[]*IgnorePackageRule]()
	var err error
	eachArr(config, func(key string, value phpjson.Value) {
		if err != nil {
			return
		}
		intKey := isIntKey(key)
		if value == nil {
			appendTo(rules, key, NewIgnorePackageRule(key, &semver.MatchAll{}, nil, true, true))
			return
		}
		if s, ok := value.(string); ok && !intKey {
			appendTo(rules, key, NewIgnorePackageRule(key, &semver.MatchAll{}, strPtr(s), true, true))
			return
		}
		if s, ok := value.(string); ok && intKey {
			appendTo(rules, s, NewIgnorePackageRule(s, &semver.MatchAll{}, nil, true, true))
			return
		}
		if phpjson.IsArray(value) && !intKey && isset(value, "0") {
			eachArr(value, func(_ string, ruleConfig phpjson.Value) {
				if err != nil {
					return
				}
				if !phpjson.IsArray(ruleConfig) {
					err = newErr(UnexpectedValueException, fmt.Sprintf("Invalid ignore rule for \"%s\": expected an object, got %s.", key, debugType(ruleConfig)))
					return
				}
				r, e := ignoreRuleFromObject(key, ruleConfig)
				if e != nil {
					err = e
					return
				}
				appendTo(rules, key, r)
			})
			return
		}
		if phpjson.IsArray(value) && !intKey {
			r, e := ignoreRuleFromObject(key, value)
			if e != nil {
				err = e
				return
			}
			appendTo(rules, key, r)
			return
		}
		err = newErr(UnexpectedValueException, fmt.Sprintf("Invalid ignore entry at key \"%s\": value of type %s is not a supported shape. Expected null, a reason string, a rule object, or a list of rule objects.", key, debugType(value)))
	})
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func ignoreRuleFromObject(packageName string, config phpjson.Value) (*IgnorePackageRule, error) {
	var constraint semver.Constraint = &semver.MatchAll{}
	if isset(config, "constraint") {
		c, err := semver.ParseConstraints(toStr(coalesce(config, "constraint", nil)))
		if err != nil {
			return nil, err
		}
		constraint = c
	}
	return NewIgnorePackageRule(packageName, constraint, nullableStr(coalesce(config, "reason", nil)), toBool(coalesce(config, "on-block", true)), toBool(coalesce(config, "on-audit", true))), nil
}

// FilterByOperation is IgnorePackageRule::filterByOperation.
func FilterByOperation(rules *IgnoreMap, op string) *IgnoreMap {
	filtered := NewOMap[[]*IgnorePackageRule]()
	rules.Each(func(name string, list []*IgnorePackageRule) {
		for _, r := range list {
			if appliesTo(op, r.OnBlock, r.OnAudit) {
				appendTo(filtered, name, r)
			}
		}
	})
	return filtered
}

// Scopes accepted by IgnoreUnreachable.With.
var IgnoreUnreachableScopes = []string{"audit", "install", "update"}

// IgnoreUnreachable is Composer\Policy\IgnoreUnreachable.
type IgnoreUnreachable struct {
	Audit   bool
	Install bool
	Update  bool
}

// ForBlockScope is IgnoreUnreachable::forBlockScope.
func (i IgnoreUnreachable) ForBlockScope(blockScope string) bool {
	if blockScope == BlockScopeInstall {
		return i.Install
	}
	return i.Update
}

func IgnoreUnreachableDefault() IgnoreUnreachable { return IgnoreUnreachable{false, true, true} }
func IgnoreUnreachableAll() IgnoreUnreachable     { return IgnoreUnreachable{true, true, true} }
func IgnoreUnreachableNone() IgnoreUnreachable    { return IgnoreUnreachable{false, false, false} }

// With is IgnoreUnreachable::with.
func (i IgnoreUnreachable) With(scopes ...string) (IgnoreUnreachable, error) {
	if len(scopes) == 0 {
		return i, newErr(InvalidArgumentException, "At least one scope is required.")
	}
	out := i
	for _, s := range scopes {
		switch s {
		case "audit":
			out.Audit = true
		case "install":
			out.Install = true
		case "update":
			out.Update = true
		default:
			return i, newErr(InvalidArgumentException, fmt.Sprintf("Unknown scope \"%s\". Expected one of %s.", s, strings.Join(IgnoreUnreachableScopes, ", ")))
		}
	}
	return out, nil
}

// IgnoreUnreachableFromRawPolicyConfig is IgnoreUnreachable::fromRawPolicyConfig.
func IgnoreUnreachableFromRawPolicyConfig(config phpjson.Value) IgnoreUnreachable {
	if !isset(config, "ignore-unreachable") {
		return IgnoreUnreachableDefault()
	}
	v := coalesce(config, "ignore-unreachable", nil)
	if phpjson.IsArray(v) {
		return IgnoreUnreachable{inArrayStrict("audit", v), inArrayStrict("install", v), inArrayStrict("update", v)}
	}
	if toBool(v) {
		return IgnoreUnreachableAll()
	}
	return IgnoreUnreachableNone()
}

// IgnoreUnreachableFromRawAuditConfig is IgnoreUnreachable::fromRawAuditConfig.
func IgnoreUnreachableFromRawAuditConfig(auditConfig phpjson.Value) IgnoreUnreachable {
	if isset(auditConfig, "ignore-unreachable") && toBool(coalesce(auditConfig, "ignore-unreachable", nil)) {
		return IgnoreUnreachable{true, false, false}
	}
	return IgnoreUnreachableDefault()
}
