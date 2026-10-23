// Package semver ports composer/semver 3.5.1 (the version locked by Composer
// f83588392df43d4fd027d84d26de459946ef4e80) together with PHP's
// version_compare(). Go module semver is deliberately not used: Composer
// versions have four components, stability suffixes and dev branches.
package semver

import (
	"strconv"
	"strings"
)

// canonicalize mirrors php_canonicalize_version() from ext/standard/versioning.c.
func canonicalize(v string) string {
	if v == "" {
		return ""
	}
	isdig := func(c byte) bool { return c >= '0' && c <= '9' && c != '.' }
	isndig := func(c byte) bool { return !(c >= '0' && c <= '9') && c != '.' }
	isalnum := func(c byte) bool {
		return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
	}
	buf := make([]byte, 0, len(v)*2)
	lp := v[0]
	buf = append(buf, v[0])
	for i := 1; i < len(v); i++ {
		c := v[i]
		last := buf[len(buf)-1]
		switch {
		case c == '-' || c == '_' || c == '+':
			if last != '.' {
				buf = append(buf, '.')
			}
		case (isndig(lp) && isdig(c)) || (isdig(lp) && isndig(c)):
			if last != '.' {
				buf = append(buf, '.')
			}
			buf = append(buf, c)
		case !isalnum(c):
			if last != '.' {
				buf = append(buf, '.')
			}
		default:
			buf = append(buf, c)
		}
		lp = c
	}
	return string(buf)
}

var specialForms = []struct {
	name  string
	order int
}{
	{"dev", 0}, {"alpha", 1}, {"a", 1}, {"beta", 2}, {"b", 2}, {"RC", 3}, {"rc", 3}, {"#", 4}, {"pl", 5}, {"p", 5},
}

func specialOrder(form string) int {
	for _, f := range specialForms {
		if strings.HasPrefix(form, f.name) {
			return f.order
		}
	}
	return -1
}

func compareSpecial(a, b string) int {
	return sign(specialOrder(a) - specialOrder(b))
}

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	}
	return 0
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// strtol parses the leading decimal digits like C strtol(p, NULL, 10).
func strtol(s string) int64 {
	end := 0
	for end < len(s) && isDigit(s[end]) {
		end++
	}
	n, err := strconv.ParseInt(s[:end], 10, 64)
	if err != nil {
		// overflow saturates like strtol
		return 1<<63 - 1
	}
	return n
}

// VersionCompare mirrors PHP version_compare($a, $b) returning -1, 0 or 1.
func VersionCompare(orig1, orig2 string) int {
	if orig1 == "" || orig2 == "" {
		if orig1 == "" && orig2 == "" {
			return 0
		}
		if orig1 == "" {
			return -1
		}
		return 1
	}
	// PHP only canonicalizes strings that do not start with '#' (the "#N#" marker)
	v1, v2 := orig1, orig2
	if orig1[0] != '#' {
		v1 = canonicalize(orig1)
	}
	if orig2[0] != '#' {
		v2 = canonicalize(orig2)
	}
	// Observed on PHP 8.4: a trailing separator ("1.0-", "1.0.") compares as
	// if absent; see testdata probes in TestVersionCompareTrailing.
	v1 = strings.TrimRight(v1, ".")
	v2 = strings.TrimRight(v2, ".")
	if v1 == "" || v2 == "" {
		return VersionCompare(v1, v2)
	}
	p1, p2 := v1, v2
	n1, n2 := -2, -2 // -2 = not yet searched, -1 = NULL
	compare := 0
	for p1 != "" && p2 != "" && n1 != -1 && n2 != -1 {
		n1 = strings.IndexByte(p1, '.')
		n2 = strings.IndexByte(p2, '.')
		part1, part2 := p1, p2
		if n1 >= 0 {
			part1 = p1[:n1]
		}
		if n2 >= 0 {
			part2 = p2[:n2]
		}
		switch {
		case part1 != "" && part2 != "" && isDigit(part1[0]) && isDigit(part2[0]):
			l1, l2 := strtol(part1), strtol(part2)
			switch {
			case l1 < l2:
				compare = -1
			case l1 > l2:
				compare = 1
			default:
				compare = 0
			}
		case !(part1 != "" && isDigit(part1[0])) && !(part2 != "" && isDigit(part2[0])):
			compare = compareSpecial(part1, part2)
		default:
			if part1 != "" && isDigit(part1[0]) {
				compare = compareSpecial("#N#", part2)
			} else {
				compare = compareSpecial(part1, "#N#")
			}
		}
		if compare != 0 {
			break
		}
		if n1 >= 0 {
			p1 = p1[n1+1:]
		}
		if n2 >= 0 {
			p2 = p2[n2+1:]
		}
	}
	if compare == 0 {
		if n1 >= 0 {
			if p1 != "" && isDigit(p1[0]) {
				compare = 1
			} else {
				compare = VersionCompare(p1, "#N#")
			}
		} else if n2 >= 0 {
			if p2 != "" && isDigit(p2[0]) {
				compare = -1
			} else {
				compare = VersionCompare("#N#", p2)
			}
		}
	}
	return compare
}

// VersionCompareOp mirrors version_compare($a, $b, $op).
func VersionCompareOp(a, b, op string) bool {
	c := VersionCompare(a, b)
	switch op {
	case "<", "lt":
		return c == -1
	case "<=", "le":
		return c != 1
	case ">", "gt":
		return c == 1
	case ">=", "ge":
		return c != -1
	case "==", "=", "eq":
		return c == 0
	case "!=", "<>", "ne":
		return c != 0
	}
	return false
}
