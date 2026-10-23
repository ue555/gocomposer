package archiver

import "strings"

// GlobToRegex is Symfony\Component\Finder\Glob::toRegex (symfony/finder
// v5.4, MIT, (c) Fabien Potencier).
func GlobToRegex(glob string, strictLeadingDot, strictWildcardSlash bool, delimiter string) string {
	firstByte := true
	escaping := false
	inCurlies := 0
	var regex strings.Builder
	n := len(glob)
	for i := 0; i < n; i++ {
		car := string(glob[i])
		if firstByte && strictLeadingDot && car != "." {
			regex.WriteString(`(?=[^\.])`)
		}
		firstByte = car == "/"
		if firstByte && strictWildcardSlash && i+2 < n && glob[i+1:i+3] == "**" && (i+3 >= n || glob[i+3] == '/') {
			car = `[^/]++/`
			if i+3 >= n {
				car += "?"
			}
			if strictLeadingDot {
				car = `(?=[^\.])` + car
			}
			car = "/(?:" + car + ")*"
			i += 2
			if i+1 < n {
				i++
			}
			if delimiter == "/" {
				car = strings.ReplaceAll(car, "/", `\/`)
			}
		}
		switch {
		case car == delimiter || car == "." || car == "(" || car == ")" || car == "|" || car == "+" || car == "^" || car == "$":
			regex.WriteString(`\` + car)
		case car == "*":
			if escaping {
				regex.WriteString(`\*`)
			} else if strictWildcardSlash {
				regex.WriteString(`[^/]*`)
			} else {
				regex.WriteString(`.*`)
			}
		case car == "?":
			if escaping {
				regex.WriteString(`\?`)
			} else if strictWildcardSlash {
				regex.WriteString(`[^/]`)
			} else {
				regex.WriteString(`.`)
			}
		case car == "{":
			if escaping {
				regex.WriteString(`\{`)
			} else {
				regex.WriteString("(")
				inCurlies++
			}
		case car == "}" && inCurlies > 0:
			if escaping {
				regex.WriteString("}")
			} else {
				regex.WriteString(")")
				inCurlies--
			}
		case car == "," && inCurlies > 0:
			if escaping {
				regex.WriteString(",")
			} else {
				regex.WriteString("|")
			}
		case car == `\`:
			if escaping {
				regex.WriteString(`\\`)
				escaping = false
			} else {
				escaping = true
			}
			continue
		default:
			regex.WriteString(car)
		}
		escaping = false
	}
	return delimiter + "^" + regex.String() + "$" + delimiter
}
