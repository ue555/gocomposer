package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Ported from composer/semver src/VersionParser.php. PCRE possessive
// quantifiers are rewritten as greedy ones (no pattern below depends on
// possessiveness for its result) and look-behinds are implemented by hand.

const modifierRegex = `[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)((?:[.-]?\d+)*)?)?([.-]?dev)?`

const DefaultBranchAlias = "9999999-dev"

var (
	reStabilityHash  = regexp.MustCompile(`#.+$`)
	reParseStability = regexp.MustCompile(`(?i)` + modifierRegex + `(?:\+.*)?$`)
	reAliasStrip     = regexp.MustCompile(`^([^,\s]+) +as +([^,\s]+)$`)
	reStabilityFlag  = regexp.MustCompile(`(?i)@(?:stable|RC|beta|alpha|dev)$`)
	reBuildMeta      = regexp.MustCompile(`^([^,\s+]+)\+[^\s]+$`)
	reClassical      = regexp.MustCompile(`(?i)^v?(\d{1,5})(\.\d+)?(\.\d+)?(\.\d+)?` + modifierRegex + `$`)
	reDate           = regexp.MustCompile(`(?i)^v?(\d{4}(?:[.:-]?\d{2}){1,6}(?:[.:-]?\d{1,3}){0,2})` + modifierRegex + `$`)
	reNonDigit       = regexp.MustCompile(`\D`)
	reDevBranch      = regexp.MustCompile(`(?i)^(.*?)[.-]?dev$`)
	reNumericAlias   = regexp.MustCompile(`(?i)^((?:\d+\.)*\d+)(?:\.x)?-dev$`)
	reBranch         = regexp.MustCompile(`(?i)^v?(\d+)(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?$`)
	reOrSplit        = regexp.MustCompile(`\s*\|\|?\s*`)
	reConstraintFlag = regexp.MustCompile(`(?i)^([^,\s]*?)@(stable|RC|beta|alpha|dev)$`)
	reRefStrip       = regexp.MustCompile(`(?i)^(dev-[^,\s@]+?|[^,\s@]+?\.x-dev)#.+$`)
	reWildcardAll    = regexp.MustCompile(`(?i)^(v)?[xX*](\.[xX*])*$`)
	versionRegex     = `v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?(?:` + modifierRegex + `|\.([xX*][.-]?dev))(?:\+[^\s]+)?`
	reTilde          = regexp.MustCompile(`(?i)^~>?` + versionRegex + `$`)
	reCaret          = regexp.MustCompile(`(?i)^\^` + versionRegex + `($)`)
	reXRange         = regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.[xX*])+$`)
	reHyphen         = regexp.MustCompile(`(?i)^(?P<from>` + versionRegex + `) +- +(?P<to>` + versionRegex + `)($)`)
	reBasic          = regexp.MustCompile(`^(<>|!=|>=?|<=?|==?)?\s*(.*)`)
	reModifierSuffix = regexp.MustCompile(`(?i)-` + modifierRegex + `$`)
	reDevSuffixChars = regexp.MustCompile(`^[0-9a-zA-Z-./]+$`)
)

// ParseError mirrors \UnexpectedValueException from the parser.
type ParseError struct{ Msg string }

func (e *ParseError) Error() string { return e.Msg }

// phpTrim mirrors PHP trim() default character list.
func phpTrim(s string) string { return strings.Trim(s, " \t\n\r\x00\x0B") }

// submatches returns PHP-like matches with "" for non-participating groups.
func submatches(re *regexp.Regexp, s string) []string {
	m := re.FindStringSubmatch(s)
	if m == nil {
		return nil
	}
	return m
}

// ParseStability ports VersionParser::parseStability.
func ParseStability(version string) string {
	version = reStabilityHash.ReplaceAllString(version, "")
	if strings.HasPrefix(version, "dev-") || strings.HasSuffix(version, "-dev") {
		return "dev"
	}
	m := submatches(reParseStability, strings.ToLower(version))
	if m == nil {
		return "stable"
	}
	if m[3] != "" {
		return "dev"
	}
	switch m[1] {
	case "beta", "b":
		return "beta"
	case "alpha", "a":
		return "alpha"
	case "rc":
		return "RC"
	}
	return "stable"
}

// NormalizeStability ports VersionParser::normalizeStability.
func NormalizeStability(s string) (string, error) {
	s = strings.ToLower(s)
	switch s {
	case "stable", "beta", "alpha", "dev":
		return s, nil
	case "rc":
		return "RC", nil
	}
	return "", &ParseError{fmt.Sprintf(`Invalid stability string "%s", expected one of stable, RC, beta, alpha or dev`, s)}
}

func expandStability(s string) string {
	switch strings.ToLower(s) {
	case "a":
		return "alpha"
	case "b":
		return "beta"
	case "p", "pl":
		return "patch"
	case "rc":
		return "RC"
	}
	return strings.ToLower(s)
}

// IsValid ports VersionParser::isValid.
func IsValid(version string) bool {
	_, err := Normalize(version)
	return err == nil
}

// Normalize ports VersionParser::normalize.
func Normalize(version string) (string, error) { return NormalizeFull(version, "") }

// NormalizeFull is normalize($version, $fullVersion).
func NormalizeFull(version, fullVersion string) (string, error) {
	version = phpTrim(version)
	origVersion := version
	if fullVersion == "" {
		fullVersion = version
	}
	if m := reAliasStrip.FindStringSubmatch(version); m != nil {
		version = m[1]
	}
	if loc := reStabilityFlag.FindStringIndex(version); loc != nil {
		version = version[:loc[0]]
	}
	switch version {
	case "master", "trunk", "default":
		version = "dev-" + version
	}
	if len(version) >= 4 && strings.EqualFold(version[:4], "dev-") {
		return "dev-" + version[4:], nil
	}
	if m := reBuildMeta.FindStringSubmatch(version); m != nil {
		version = m[1]
	}
	var matches []string
	index := -1
	if m := submatches(reClassical, version); m != nil {
		version = m[1]
		for _, g := range m[2:5] {
			if g != "" {
				version += g
			} else {
				version += ".0"
			}
		}
		matches = m
		index = 5
	} else if m := submatches(reDate, version); m != nil {
		version = reNonDigit.ReplaceAllString(m[1], ".")
		matches = m
		index = 2
	}
	if index >= 0 {
		if matches[index] != "" {
			if matches[index] == "stable" {
				return version, nil
			}
			version += "-" + expandStability(matches[index]) + strings.TrimLeft(matches[index+1], ".-")
		}
		if matches[index+2] != "" {
			version += "-dev"
		}
		return version, nil
	}
	if m := reDevBranch.FindStringSubmatch(version); m != nil {
		normalized := NormalizeBranch(m[1])
		if !strings.Contains(normalized, "dev-") {
			return normalized, nil
		}
	}
	extra := ""
	q := regexp.QuoteMeta(version)
	if regexp.MustCompile(` +as +` + q + `(?:@(?:stable|RC|beta|alpha|dev))?$`).MatchString(fullVersion) {
		extra = ` in "` + fullVersion + `", the alias must be an exact version`
	} else if regexp.MustCompile(`^` + q + `(?:@(?:stable|RC|beta|alpha|dev))? +as +`).MatchString(fullVersion) {
		extra = ` in "` + fullVersion + `", the alias source must be an exact version, if it is a branch name you should prefix it with dev-`
	}
	return "", &ParseError{`Invalid version string "` + origVersion + `"` + extra}
}

// MustNormalize panics when the version is invalid; for internal constants.
func MustNormalize(v string) string {
	n, err := Normalize(v)
	if err != nil {
		panic(err)
	}
	return n
}

// ParseNumericAliasPrefix ports VersionParser::parseNumericAliasPrefix.
func ParseNumericAliasPrefix(branch string) (string, bool) {
	if m := reNumericAlias.FindStringSubmatch(branch); m != nil {
		return m[1] + ".", true
	}
	return "", false
}

// NormalizeBranch ports VersionParser::normalizeBranch.
func NormalizeBranch(name string) string {
	name = phpTrim(name)
	if m := reBranch.FindStringSubmatch(name); m != nil {
		version := ""
		for i := 1; i < 5; i++ {
			if m[i] != "" {
				version += strings.NewReplacer("*", "x", "X", "x").Replace(m[i])
			} else {
				version += ".x"
			}
		}
		return strings.ReplaceAll(version, "x", "9999999") + "-dev"
	}
	return "dev-" + name
}

// NormalizeDefaultBranch ports VersionParser::normalizeDefaultBranch.
func NormalizeDefaultBranch(name string) string {
	if name == "dev-master" || name == "dev-default" || name == "dev-trunk" {
		return DefaultBranchAlias
	}
	return name
}

// ParseConstraints ports VersionParser::parseConstraints.
func ParseConstraints(constraints string) (Constraint, error) {
	pretty := constraints
	orConstraints := reOrSplit.Split(phpTrim(constraints), -1)
	var orGroups []Constraint
	for _, orConstraint := range orConstraints {
		andConstraints := splitAnd(orConstraint)
		var objs []Constraint
		if len(andConstraints) > 1 {
			for _, and := range andConstraints {
				parsed, err := parseConstraint(and)
				if err != nil {
					return nil, err
				}
				objs = append(objs, parsed...)
			}
		} else {
			parsed, err := parseConstraint(andConstraints[0])
			if err != nil {
				return nil, err
			}
			objs = parsed
		}
		if len(objs) == 1 {
			orGroups = append(orGroups, objs[0])
		} else {
			orGroups = append(orGroups, NewMulti(objs, true))
		}
	}
	c := CreateMulti(orGroups, false)
	c.SetPrettyString(pretty)
	return c, nil
}

// MustParseConstraints is for literals in code and tests.
func MustParseConstraints(s string) Constraint {
	c, err := ParseConstraints(s)
	if err != nil {
		panic(err)
	}
	return c
}

// splitAnd emulates preg_split('{(?<!^|as|[=>< ,]) *(?<!-)[, ](?!-) *(?!,|as|$)}', $s)
// with PCRE's leftmost, greedy-then-backtrack search order.
func splitAnd(s string) []string {
	var parts []string
	last := 0
	start := 0
	for start <= len(s) {
		ms, me := -1, -1
		for p := start; p < len(s) && ms < 0; p++ {
			if p == 0 || (p >= 2 && s[p-2:p] == "as") || strings.IndexByte("=>< ,", s[p-1]) >= 0 {
				continue
			}
			maxK := 0
			for p+maxK < len(s) && s[p+maxK] == ' ' {
				maxK++
			}
			for k := maxK; k >= 0 && ms < 0; k-- {
				q := p + k
				if q >= 1 && s[q-1] == '-' {
					continue
				}
				if q >= len(s) || (s[q] != ',' && s[q] != ' ') {
					continue
				}
				r := q + 1
				if r < len(s) && s[r] == '-' {
					continue
				}
				maxM := 0
				for r+maxM < len(s) && s[r+maxM] == ' ' {
					maxM++
				}
				for m := maxM; m >= 0; m-- {
					e := r + m
					if e == len(s) || s[e] == ',' || strings.HasPrefix(s[e:], "as") {
						continue
					}
					ms, me = p, e
					break
				}
			}
		}
		if ms < 0 {
			break
		}
		parts = append(parts, s[last:ms])
		last = me
		start = me
	}
	return append(parts, s[last:])
}

func isEmptyPHP(s string) bool { return s == "" || s == "0" }

func parseConstraint(constraint string) ([]Constraint, error) {
	if m := reAliasStrip.FindStringSubmatch(constraint); m != nil {
		constraint = m[1]
	}
	stabilityModifier := ""
	if m := reConstraintFlag.FindStringSubmatch(constraint); m != nil {
		if m[1] != "" {
			constraint = m[1]
		} else {
			constraint = "*"
		}
		if m[2] != "stable" {
			stabilityModifier = m[2]
		}
	}
	if m := reRefStrip.FindStringSubmatch(constraint); m != nil {
		constraint = m[1]
	}
	if m := reWildcardAll.FindStringSubmatch(constraint); m != nil {
		if m[1] != "" || m[2] != "" {
			return []Constraint{NewConstraint(">=", "0.0.0.0-dev")}, nil
		}
		return []Constraint{&MatchAll{}}, nil
	}

	if m := reTilde.FindStringSubmatch(constraint); m != nil {
		if strings.HasPrefix(constraint, "~>") {
			return nil, &ParseError{"Could not parse version constraint " + constraint + `: Invalid operator "~>", you probably meant to use the "~" operator`}
		}
		position := 1
		switch {
		case m[4] != "":
			position = 4
		case m[3] != "":
			position = 3
		case m[2] != "":
			position = 2
		}
		if m[8] != "" {
			position++
		}
		suffix := ""
		if isEmptyPHP(m[5]) && isEmptyPHP(m[7]) && isEmptyPHP(m[8]) {
			suffix = "-dev"
		}
		low, err := Normalize((constraint + suffix)[1:])
		if err != nil {
			return nil, err
		}
		highPosition := position - 1
		if highPosition < 1 {
			highPosition = 1
		}
		high := manipulateVersionString(m, highPosition, 1) + "-dev"
		return []Constraint{NewConstraint(">=", low), NewConstraint("<", high)}, nil
	}

	if m := reCaret.FindStringSubmatch(constraint); m != nil {
		position := 3
		if m[1] != "0" || m[2] == "" {
			position = 1
		} else if m[2] != "0" || m[3] == "" {
			position = 2
		}
		suffix := ""
		if isEmptyPHP(m[5]) && isEmptyPHP(m[7]) && isEmptyPHP(m[8]) {
			suffix = "-dev"
		}
		low, err := Normalize((constraint + suffix)[1:])
		if err != nil {
			return nil, err
		}
		high := manipulateVersionString(m, position, 1) + "-dev"
		return []Constraint{NewConstraint(">=", low), NewConstraint("<", high)}, nil
	}

	if m := reXRange.FindStringSubmatch(constraint); m != nil {
		position := 1
		if m[3] != "" {
			position = 3
		} else if m[2] != "" {
			position = 2
		}
		mm := []string{m[0], m[1], m[2], m[3], ""}
		low := manipulateVersionString(mm, position, 0) + "-dev"
		high := manipulateVersionString(mm, position, 1) + "-dev"
		if low == "0.0.0.0-dev" {
			return []Constraint{NewConstraint("<", high)}, nil
		}
		return []Constraint{NewConstraint(">=", low), NewConstraint("<", high)}, nil
	}

	if m := reHyphen.FindStringSubmatch(constraint); m != nil {
		lowSuffix := ""
		if isEmptyPHP(m[6]) && isEmptyPHP(m[8]) && isEmptyPHP(m[9]) {
			lowSuffix = "-dev"
		}
		low, err := Normalize(m[1])
		if err != nil {
			return nil, err
		}
		lower := NewConstraint(">=", low+lowSuffix)
		empty := func(x string) bool {
			if x == "0" {
				return false
			}
			return x == ""
		}
		var upper *Single
		if (!empty(m[12]) && !empty(m[13])) || m[15] != "" || m[17] != "" || m[18] != "" {
			high, err := Normalize(m[10])
			if err != nil {
				return nil, err
			}
			upper = NewConstraint("<=", high)
		} else {
			highMatch := []string{"", m[11], m[12], m[13], m[14]}
			if _, err := Normalize(m[10]); err != nil {
				return nil, err
			}
			pos := 2
			if empty(m[12]) {
				pos = 1
			}
			upper = NewConstraint("<", manipulateVersionString(highMatch, pos, 1)+"-dev")
		}
		return []Constraint{lower, upper}, nil
	}

	var lastErr error
	if m := reBasic.FindStringSubmatch(constraint); m != nil {
		version, err := Normalize(m[2])
		if err != nil {
			if strings.HasSuffix(m[2], "-dev") && reDevSuffixChars.MatchString(m[2]) {
				version, err = Normalize("dev-" + m[2][:len(m[2])-4])
			}
		}
		if err == nil {
			op := m[1]
			if op == "" {
				op = "="
			}
			if op != "==" && op != "=" && stabilityModifier != "" && ParseStability(version) == "stable" {
				version += "-" + stabilityModifier
			} else if op == "<" || op == ">=" {
				if !reModifierSuffix.MatchString(m[2]) && !strings.HasPrefix(m[2], "dev-") {
					version += "-dev"
				}
			}
			return []Constraint{NewConstraint(op, version)}, nil
		}
		lastErr = err
	}
	msg := "Could not parse version constraint " + constraint
	if lastErr != nil {
		msg += ": " + lastErr.Error()
	}
	return nil, &ParseError{msg}
}

// manipulateVersionString ports the private helper; matches[1..4] are used.
func manipulateVersionString(matches []string, position, increment int) string {
	m := make([]string, 5)
	copy(m, matches)
	for i := 4; i > 0; i-- {
		if i > position {
			m[i] = "0"
		} else if i == position && increment != 0 {
			n, err := strconv.ParseInt(m[i], 10, 64)
			if err != nil {
				n = 0
			}
			m[i] = strconv.FormatInt(n+int64(increment), 10)
		}
	}
	return m[1] + "." + m[2] + "." + m[3] + "." + m[4]
}

// SplitAnd exposes the AND-constraint splitter used by parseConstraints.
func SplitAnd(s string) []string { return splitAnd(s) }
