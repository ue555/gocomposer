package validate

// ValidEmail ports PHP's filter_var($v, FILTER_VALIDATE_EMAIL)
// (php-src ext/filter/logical_filters.c, php_filter_validate_email, PHP
// License 3.01). The PCRE look-aheads of the original pattern are evaluated
// by hand; the rest is matched byte-wise like PCRE without /u.

import (
	"regexp"
	"strings"
)

const (
	emailAtom     = `[\x21\x23-\x27\x2A\x2B\x2D\x2F-\x39\x3D\x3F\x5E-\x7E]+`
	emailQuoted   = `\x22(?:[\x01-\x08\x0B\x0C\x0E-\x1F\x21\x23-\x5B\x5D-\x7F]|\x5C[\x00-\x7F])*\x22`
	emailWord     = `(?:` + emailAtom + `|` + emailQuoted + `)`
	emailIPv4     = `(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])(?:\.(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])){3}`
	emailHex      = `[a-f0-9]{1,4}`
	emailIPv6Full = emailHex + `(?::` + emailHex + `){7}`
	emailIPv6Comp = `(?:` + emailHex + `(?::` + emailHex + `){0,5})?::(?:` + emailHex + `(?::` + emailHex + `){0,5})?`
	emailIPv4v6F  = emailHex + `(?::` + emailHex + `){5}:`
	emailIPv4v6C  = `(?:` + emailHex + `(?::` + emailHex + `){0,3})?::(?:` + emailHex + `(?::` + emailHex + `){0,3}:)?`
)

var (
	reEmailLocal  = regexp.MustCompile(`(?i)^` + emailWord + `(?:\.` + emailWord + `)*$`)
	reEmailDomain = regexp.MustCompile(`(?i)^(?:(?:xn--)?[a-z0-9]+(?:-+[a-z0-9]+)*\.)+(?:[a-z][a-z0-9]*|xn--[a-z0-9]+)(?:-+[a-z0-9]+)*$`)
	reIPv6Full    = regexp.MustCompile(`(?i)^IPv6:` + emailIPv6Full + `$`)
	reIPv6Comp    = regexp.MustCompile(`(?i)^IPv6:` + emailIPv6Comp + `$`)
	reIPv4Mapped  = regexp.MustCompile(`(?i)^(?:IPv6:(?:` + emailIPv4v6F + `|` + emailIPv4v6C + `))?` + emailIPv4 + `$`)
	reHexColon    = regexp.MustCompile(`(?i)[a-f0-9][:\]]`)
	reHexColon2   = regexp.MustCompile(`(?i)[a-f0-9]:`)
)

// unitCounts returns, for every offset, the maximum number of look-ahead
// "units" ((?:"?\\[\x00-\x7E]"?)|(?:"?[^\\"]"?)) that can end there.
func unitCounts(s string) []int {
	best := make([]int, len(s)+1)
	for i := range best {
		best[i] = -1
	}
	best[0] = 0
	for i := 0; i <= len(s); i++ {
		if best[i] < 0 {
			continue
		}
		starts := []int{i}
		if i < len(s) && s[i] == '"' {
			starts = append(starts, i+1)
		}
		for _, st := range starts {
			var ends []int
			if st+1 < len(s) && s[st] == '\\' && s[st+1] <= 0x7E {
				ends = append(ends, st+2)
			}
			if st < len(s) && s[st] != '\\' && s[st] != '"' {
				ends = append(ends, st+1)
			}
			for _, e := range ends {
				for _, e2 := range []int{e, e + 1} {
					if e2 > e && !(e < len(s) && s[e] == '"') {
						continue
					}
					if best[i]+1 > best[e2] {
						best[e2] = best[i] + 1
					}
				}
			}
		}
	}
	return best
}

// ValidEmail is FILTER_VALIDATE_EMAIL.
func ValidEmail(s string) bool {
	if len(s) > 320 {
		return false
	}
	counts := unitCounts(s)
	for i, c := range counts {
		if c >= 255 {
			return false
		}
		if c >= 65 && i < len(s) && s[i] == '@' {
			return false
		}
	}
	// the local part may itself contain "@" inside quotes: try every "@"
	for at := strings.LastIndexByte(s, '@'); at > 0; at = strings.LastIndexByte(s[:at], '@') {
		local, domain := latin1(s[:at]), s[at+1:]
		if !reEmailLocal.MatchString(local) {
			continue
		}
		if validEmailDomain(domain) {
			return true
		}
	}
	return false
}

func validEmailDomain(d string) bool {
	if strings.HasPrefix(d, "[") && strings.HasSuffix(d, "]") {
		inner := d[1 : len(d)-1]
		rest := d[1:]
		if reIPv6Full.MatchString(inner) {
			return true
		}
		if reIPv6Comp.MatchString(inner) && len(reHexColon.FindAllStringIndex(strings.TrimPrefix(rest, "IPv6:"), -1)) < 7 {
			return true
		}
		if reIPv4Mapped.MatchString(inner) {
			if strings.Contains(inner, "::") && len(reHexColon2.FindAllStringIndex(strings.TrimPrefix(rest, "IPv6:"), -1)) >= 5 {
				return false
			}
			return true
		}
		return false
	}
	// (?!.*[^.]{64,}): no run of 64 non-dot characters
	run := 0
	for i := 0; i < len(d); i++ {
		if d[i] == '\n' {
			break
		}
		if d[i] == '.' {
			run = 0
			continue
		}
		run++
		if run >= 64 {
			return false
		}
	}
	return reEmailDomain.MatchString(latin1(d))
}

func latin1(s string) string {
	r := make([]rune, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = rune(s[i])
	}
	return string(r)
}
