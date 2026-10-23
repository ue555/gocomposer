// Ported for gocomposer (MIT) from justinrainbow/json-schema 6.13.1
// (MIT, Copyright (c) 2016 justinrainbow/json-schema contributors):
// src/JsonSchema/Constraints/FormatConstraint.php, src/JsonSchema/Rfc3339.php,
// src/JsonSchema/Tool/Validator/UriValidator.php and
// src/JsonSchema/Tool/Validator/RelativeReferenceValidator.php.
//
// The PHP functions those files call are reimplemented with PHP 8.4
// semantics: DateTime::createFromFormat() round trips for 'Y-m-d', 'H:i:s'
// and 'U'; filter_var() FILTER_VALIDATE_EMAIL (FILTER_FLAG_EMAIL_UNICODE;
// the pattern is copied from PHP 8.4's ext/filter/logical_filters.c, where it
// carries the notice "Copyright (c) Michael Rushton 2009-10
// http://squiloople.com/ Feel free to use and redistribute this code. But
// please keep this copyright notice.") and
// FILTER_VALIDATE_IP (_php_filter_validate_ipv4/_ipv6). Unknown format names
// are ignored, exactly like the PHP switch's default branch.
package jsonschema

import (
	"strconv"
	"strings"

	"gocomposer/internal/phpjson"
)

func (v *validator) formatCheck(in *inst, element Value, schema *phpjson.Map, path jsonPointer) {
	if !isset(schema, "format") || v.disableFormat {
		return
	}
	format := get(schema, "format")
	str, isStr := element.(string)
	is := func(name string) bool { return looseEqual(format, name) }
	switch {
	case is("date"):
		if isStr && !validateDateYmd(str) {
			v.addError(in, "dateFormat", path, param{"date", element}, param{"format", format})
		}
	case is("time"):
		if isStr && !validateTimeHis(str) {
			v.addError(in, "timeFormat", path, param{"time", jsonEncode(element)}, param{"format", format})
		}
	case is("date-time"):
		if isStr && !rfc3339Valid(str) {
			v.addError(in, "dateTimeFormat", path, param{"dateTime", jsonEncode(element)}, param{"format", format})
		}
	case is("utc-millisec"):
		if !(isStr && validateUnixTimestamp(str)) {
			v.addError(in, "dateUtcFormat", path, param{"value", element}, param{"format", format})
		}
	case is("regex"):
		if isStr {
			if _, err := compileDelimited("~" + strings.ReplaceAll(str, "~", `\~`) + "~u"); err != nil {
				v.addError(in, "regexFormat", path, param{"value", element}, param{"format", format})
			}
		}
	case is("color"):
		if isStr && !validateColor(str) {
			v.addError(in, "colorFormat", path, param{"format", format})
		}
	case is("style"):
		if !isStr {
			throw("TypeError", "rtrim(): Argument #1 ($string) must be of type string, "+debugType(element)+" given")
		}
		if !validateStyle(str) {
			v.addError(in, "styleFormat", path, param{"format", format})
		}
	case is("phone"):
		if !isStr {
			throw("TypeError", "preg_match(): Argument #2 ($subject) must be of type string, "+debugType(element)+" given")
		}
		if pregMatch(phoneRe, str) != 1 {
			v.addError(in, "phoneFormat", path, param{"format", format})
		}
	case is("uri"):
		if isStr && !uriValid(str) {
			v.addError(in, "urlFormat", path, param{"format", format})
		}
	case is("uriref"), is("uri-reference"):
		if isStr && !(uriValid(str) || relativeReferenceValid(str)) {
			v.addError(in, "urlFormat", path, param{"format", format})
		}
	case is("email"):
		if isStr && !emailValid(str) {
			v.addError(in, "emailFormat", path, param{"format", format})
		}
	case is("ip-address"), is("ipv4"):
		if isStr && !ipValid(str, true) {
			v.addError(in, "ipFormat", path, param{"format", format})
		}
	case is("ipv6"):
		if isStr && !ipValid(str, false) {
			v.addError(in, "ipFormat", path, param{"format", format})
		}
	case is("host-name"), is("hostname"):
		if isStr && pregMatch(hostnameRe, str) != 1 {
			v.addError(in, "styleHostName", path, param{"format", format})
		}
	}
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return s != ""
}

func atoi2(s string) int { n, _ := strconv.Atoi(s); return n }

func isLeap(y int) bool { return y%4 == 0 && (y%100 != 0 || y%400 == 0) }

func daysIn(y, m int) int {
	switch m {
	case 2:
		if isLeap(y) {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	}
	return 31
}

func validDate(y, m, d int) bool { return m >= 1 && m <= 12 && d >= 1 && d <= daysIn(y, m) }

// validateDateYmd mirrors FormatConstraint::validateDateTime($s, 'Y-m-d'):
// createFromFormat must succeed and format() must reproduce the input.
func validateDateYmd(s string) bool {
	if len(s) != 10 || s[4] != '-' || s[7] != '-' || !allDigits(s[:4]) || !allDigits(s[5:7]) || !allDigits(s[8:]) {
		return false
	}
	return validDate(atoi2(s[:4]), atoi2(s[5:7]), atoi2(s[8:]))
}

// validateTimeHis mirrors validateDateTime($s, 'H:i:s').
func validateTimeHis(s string) bool {
	if len(s) != 8 || s[2] != ':' || s[5] != ':' || !allDigits(s[:2]) || !allDigits(s[3:5]) || !allDigits(s[6:]) {
		return false
	}
	return atoi2(s[:2]) < 24 && atoi2(s[3:5]) < 60 && atoi2(s[6:]) < 60
}

// validateUnixTimestamp mirrors validateDateTime($s, 'U') for strings: only
// canonical decimal integers round trip through format('U').
func validateUnixTimestamp(s string) bool {
	t := strings.TrimPrefix(s, "-")
	if !allDigits(t) || (len(t) > 1 && t[0] == '0') || s == "-0" {
		return false
	}
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

var rfc3339Re = mustRegex(`^(\d{4}-\d{2}-\d{2}[T ](0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):((?:[0-5][0-9]|60)))(\.\d+)?(Z|([+-](0[0-9]|1[0-9]|2[0-3]))(:)?([0-5][0-9]))$`, "D")

// rfc3339Valid mirrors Rfc3339::createFromString() !== null.
func rfc3339Valid(input string) bool {
	upper := asciiOnlyUpper(input)
	m, err := rfc3339Re.exec(upper)
	if err != nil || m == nil {
		return false
	}
	dt := m.group(1)
	y, mo, d := atoi2(dt[0:4]), atoi2(dt[5:7]), atoi2(dt[8:10])
	if !validDate(y, mo, d) {
		return false
	}
	h, mi, sec := atoi2(m.group(2)), atoi2(m.group(3)), atoi2(m.group(4))
	if sec < 60 {
		return true
	}
	offset := 0
	if tz := m.group(6); tz != "Z" {
		offset = atoi2(m.group(8))*60 + atoi2(m.group(10))
		if tz[0] == '-' {
			offset = -offset
		}
	}
	utc := ((h*60+mi-offset)%1440 + 1440) % 1440
	return utc == 23*60+59
}

func asciiOnlyUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}

var (
	colorRe    = mustRegex(`^#([a-f0-9]{3}|[a-f0-9]{6})$`, "i")
	styleRe    = mustRegex(`^\s*[-a-z]+\s*:\s*.+$`, "i")
	phoneRe    = mustRegex(`^\+?(\(\d{3}\)|\d{3}) \d{3} \d{4}$`, "")
	hostnameRe = mustRegex(`^(?!-)(?!.*?[^A-Za-z0-9\-\.])(?:(?!-)[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?\.)*(?!-)[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?$`, "")
)

var colorNames = []string{"aqua", "black", "blue", "fuchsia", "gray", "green", "lime", "maroon", "navy", "olive", "orange", "purple", "red", "silver", "teal", "white", "yellow"}

func validateColor(s string) bool {
	low := strings.ToLower(s)
	for _, c := range colorNames {
		if c == low {
			return true
		}
	}
	return pregMatch(colorRe, s) == 1
}

func validateStyle(s string) bool {
	for _, p := range strings.Split(strings.TrimRight(s, ";"), ";") {
		if pregMatch(styleRe, p) != 1 {
			return false
		}
	}
	return true
}

// ---- UriValidator / RelativeReferenceValidator ----

var (
	uriHierRe     = mustRegex(`^([a-z][a-z0-9+\-.]*):\/\/(?:([^:@\/?#]+)(?::([^@\/?#]*))?@)?([a-z0-9._~-]+|\[[a-f0-9:.]+\])(?::(\d{1,5}))?(\/[a-zA-Z0-9._~!$&'()*+,;=:@\/%-]*)*(\?([^#]*))?(\#(.*))?$`, "i")
	uriNonHierRe  = mustRegex(`^(mailto|data|urn|news|tel|file):(.+)$`, "i")
	uriNewsRe     = mustRegex(`^[a-z0-9]+(\.[a-z0-9]+)*$`, "i")
	uriTelRe      = mustRegex(`^\+?[0-9.\-() ]+$`, "")
	uriEmailRe    = mustRegex(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "")
	uriBadPathRe  = mustRegex(`[<>{}|\\^`+"`"+`]`, "")
	relRefRe      = mustRegex(`^(([^\/?#]+):)?(\/\/([^\/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?$`, "")
	relHTTPRe     = mustRegex(`^(http|https):\/\/`, "")
	relNoScheme   = mustRegex(`^:\/\/`, "")
	relBadSep     = mustRegex(`^:\/`, "")
	relEmptyAuth  = mustRegex(`^\/\/$`, "")
	relThreeSlash = mustRegex(`^\/\/\/[^\/]`, "")
	relSpace      = mustRegex(`\s`, "")
	relPercent    = mustRegex(`%(?![0-9A-Fa-f]{2})`, "")
	dotDotRe      = mustRegex(`\.\.`, "")
)

func phpEmpty(s string) bool { return s == "" || s == "0" }

func uriValid(uri string) bool {
	if m, err := uriHierRe.exec(uri); err == nil && m != nil {
		if host := m.group(4); !phpEmpty(host) && pregMatch(dotDotRe, host) == 1 {
			return false
		}
		if port := m.group(5); !phpEmpty(port) {
			n, _ := strconv.Atoi(port)
			if n < 1 || n > 65535 {
				return false
			}
		}
		if p := m.group(6); !phpEmpty(p) && pregMatch(uriBadPathRe, p) == 1 {
			return false
		}
		return true
	}
	if m, err := uriNonHierRe.exec(uri); err == nil && m != nil {
		scheme := strings.ToLower(m.group(1))
		rest := m.group(2)
		switch scheme {
		case "mailto":
			return pregMatch(uriEmailRe, rest) == 1
		case "news":
			return pregMatch(uriNewsRe, rest) == 1
		case "tel":
			return pregMatch(uriTelRe, rest) == 1
		}
		return true
	}
	return false
}

func relativeReferenceValid(ref string) bool {
	if pregMatch(relRefRe, ref) != 1 {
		return false
	}
	if strings.Contains(ref, `\`) {
		return false
	}
	for _, re := range []*reProg{relHTTPRe, relNoScheme, relBadSep, relEmptyAuth, relThreeSlash, relSpace, relPercent} {
		if pregMatch(re, ref) == 1 {
			return false
		}
	}
	return true
}

// ---- filter_var ----

var emailUnicodeRe = mustRegex(`^(?!(?:(?:\x22?\x5C[\x00-\x7E]\x22?)|(?:\x22?[^\x5C\x22]\x22?)){255,})(?!(?:(?:\x22?\x5C[\x00-\x7E]\x22?)|(?:\x22?[^\x5C\x22]\x22?)){65,}@)(?:(?:[\x21\x23-\x27\x2A\x2B\x2D\x2F-\x39\x3D\x3F\x5E-\x7E\pL\pN]+)|(?:\x22(?:[\x01-\x08\x0B\x0C\x0E-\x1F\x21\x23-\x5B\x5D-\x7F\pL\pN]|(?:\x5C[\x00-\x7F]))*\x22))(?:\.(?:(?:[\x21\x23-\x27\x2A\x2B\x2D\x2F-\x39\x3D\x3F\x5E-\x7E\pL\pN]+)|(?:\x22(?:[\x01-\x08\x0B\x0C\x0E-\x1F\x21\x23-\x5B\x5D-\x7F\pL\pN]|(?:\x5C[\x00-\x7F]))*\x22)))*@(?:(?:(?!.*[^.]{64,})(?:(?:(?:xn--)?[a-z0-9]+(?:-+[a-z0-9]+)*\.){1,126}){1,}(?:(?:[a-z][a-z0-9]*)|(?:(?:xn--)[a-z0-9]+))(?:-+[a-z0-9]+)*)|(?:\[(?:(?:IPv6:(?:(?:[a-f0-9]{1,4}(?::[a-f0-9]{1,4}){7})|(?:(?!(?:.*[a-f0-9][:\]]){7,})(?:[a-f0-9]{1,4}(?::[a-f0-9]{1,4}){0,5})?::(?:[a-f0-9]{1,4}(?::[a-f0-9]{1,4}){0,5})?)))|(?:(?:IPv6:(?:(?:[a-f0-9]{1,4}(?::[a-f0-9]{1,4}){5}:)|(?:(?!(?:.*[a-f0-9]:){5,})(?:[a-f0-9]{1,4}(?::[a-f0-9]{1,4}){0,3})?::(?:[a-f0-9]{1,4}(?::[a-f0-9]{1,4}){0,3}:)?)))?(?:(?:25[0-5])|(?:2[0-4][0-9])|(?:1[0-9]{2})|(?:[1-9]?[0-9]))(?:\.(?:(?:25[0-5])|(?:2[0-4][0-9])|(?:1[0-9]{2})|(?:[1-9]?[0-9]))){3}))\]))$`, "iDu")

func emailValid(s string) bool {
	if len(s) > 320 {
		return false
	}
	return pregMatch(emailUnicodeRe, s) == 1
}

func validateIPv4(s string) bool {
	n := 0
	i := 0
	for i < len(s) {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
		leadingZero := s[i] == '0'
		m := 1
		num := int(s[i] - '0')
		i++
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			num = num*10 + int(s[i]-'0')
			i++
			m++
			if num > 255 || m > 3 {
				return false
			}
		}
		if leadingZero && (num != 0 || m > 1) {
			return false
		}
		n++
		if n == 4 {
			return i == len(s)
		}
		if i >= len(s) || s[i] != '.' {
			return false
		}
		i++
	}
	return false
}

func validateIPv6(str string) bool {
	if !strings.Contains(str, ":") {
		return false
	}
	compressedPos := -1
	blocks := 0
	if dot := strings.IndexByte(str, '.'); dot >= 0 {
		ip4 := dot
		for ip4 > 0 && str[ip4-1] != ':' {
			ip4--
		}
		if !validateIPv4(str[ip4:]) {
			return false
		}
		l := ip4
		if l < 2 {
			return false
		}
		if str[ip4-2] != ':' {
			l--
		}
		str = str[:l]
		blocks = 2
	}
	i := 0
	end := len(str)
	for i < end {
		if str[i] == ':' {
			i++
			if i >= end {
				return false
			}
			if str[i] == ':' {
				if compressedPos >= 0 {
					return false
				}
				compressedPos = blocks
				blocks++
				i++
				if i == end {
					return blocks <= 8
				}
			} else if i-1 == 0 {
				return false
			}
		}
		n := 0
		for i < end && isHex(rune(str[i])) {
			n++
			i++
		}
		if n < 1 || n > 4 {
			return false
		}
		blocks++
		if blocks > 8 {
			return false
		}
	}
	return (compressedPos >= 0 && blocks <= 8) || blocks == 8
}

// ipValid mirrors filter_var($s, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4 or
// FILTER_FLAG_IPV6).
func ipValid(s string, v4 bool) bool {
	isV6 := strings.Contains(s, ":")
	if !isV6 && !strings.Contains(s, ".") {
		return false
	}
	if v4 {
		return !isV6 && validateIPv4(s)
	}
	return isV6 && validateIPv6(s)
}
