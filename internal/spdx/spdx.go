// Package spdx ports composer/spdx-licenses 1.6.0 (src/SpdxLicenses.php and
// its res/*.json data, copied verbatim) as locked by composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80.
// composer/spdx-licenses is MIT licensed, Copyright (C) 2015 Composer (see LICENSE).
package spdx

import (
	_ "embed"
	"encoding/json"
	"sort"
	"strings"
	"unicode"
)

//go:embed spdx-licenses.json
var licensesJSON []byte

//go:embed spdx-exceptions.json
var exceptionsJSON []byte

type license struct {
	Identifier, Name string
	OsiApproved      bool
	Deprecated       bool
}

var (
	licenses   map[string]license
	exceptions map[string][2]string
	// lower-cased ids, in rsort() order like the PHP alternation
	licenseIDs   []string
	exceptionIDs []string
)

func init() {
	var l map[string][3]interface{}
	if err := json.Unmarshal(licensesJSON, &l); err != nil {
		panic(err)
	}
	licenses = map[string]license{}
	for id, v := range l {
		name, _ := v[0].(string)
		osi, _ := v[1].(bool)
		dep, _ := v[2].(bool)
		licenses[strings.ToLower(id)] = license{id, name, osi, dep}
		licenseIDs = append(licenseIDs, id)
	}
	var e map[string][1]string
	if err := json.Unmarshal(exceptionsJSON, &e); err != nil {
		panic(err)
	}
	exceptions = map[string][2]string{}
	for id, v := range e {
		exceptions[strings.ToLower(id)] = [2]string{id, v[0]}
		exceptionIDs = append(exceptionIDs, id)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(licenseIDs)))
	sort.Sort(sort.Reverse(sort.StringSlice(exceptionIDs)))
}

// LicenseInfo is getLicenseByIdentifier's [name, osiApproved, url, deprecated].
type LicenseInfo struct {
	Name        string
	OsiApproved bool
	URL         string
	Deprecated  bool
}

// LicenseByIdentifier is SpdxLicenses::getLicenseByIdentifier.
func LicenseByIdentifier(id string) *LicenseInfo {
	l, ok := licenses[strings.ToLower(id)]
	if !ok {
		return nil
	}
	return &LicenseInfo{l.Name, l.OsiApproved, "https://spdx.org/licenses/" + l.Identifier + ".html#licenseText", l.Deprecated}
}

// Validate is SpdxLicenses::validate for a string or list of strings.
func Validate(licenseList []string) bool {
	s := ""
	switch len(licenseList) {
	case 0:
		s = ""
	case 1:
		s = licenseList[0]
	default:
		s = "(" + strings.Join(licenseList, " OR ") + ")"
	}
	return IsValidLicenseString(s)
}

// IsValidLicenseString is SpdxLicenses::isValidLicenseString: the PCRE
// grammar is evaluated as a parser returning every possible end offset,
// which equals PCRE2's backtracking into non-atomic recursions.
func IsValidLicenseString(s string) bool {
	if _, ok := licenses[strings.ToLower(s)]; ok {
		return true
	}
	p := &parser{s: s}
	ends := map[int]bool{}
	for _, lit := range []string{"NONE", "NOASSERTION"} {
		if p.literal(0, lit) {
			ends[len(lit)] = true
		}
	}
	for e := range p.licenseExpression(0) {
		ends[e] = true
	}
	// "$" matches at the end or before a final newline
	return ends[len(s)] || (strings.HasSuffix(s, "\n") && ends[len(s)-1])
}

type set map[int]bool

type parser struct {
	s     string
	memo  map[[2]int]set
	depth int
}

func (p *parser) literal(pos int, lit string) bool {
	return pos+len(lit) <= len(p.s) && strings.EqualFold(p.s[pos:pos+len(lit)], lit)
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

// spaces returns all end positions of \s+ from pos.
func (p *parser) spaces(pos int) []int {
	var out []int
	for i := pos; i < len(p.s) && isSpace(p.s[i]); i++ {
		out = append(out, i+1)
	}
	return out
}

// spacesStar returns all end positions of \s* from pos.
func (p *parser) spacesStar(pos int) []int {
	return append([]int{pos}, p.spaces(pos)...)
}

func isIDChar(b byte) bool {
	// [\pL\pN.-] without the /u modifier: bytes as Latin-1 code points
	r := rune(b)
	return b == '.' || b == '-' || unicode.IsLetter(r) || unicode.IsNumber(r)
}

func (p *parser) idstring(pos int) []int {
	var out []int
	for i := pos; i < len(p.s) && isIDChar(p.s[i]); i++ {
		out = append(out, i+1)
	}
	return out
}

func (p *parser) fromList(pos int, ids []string) []int {
	var out []int
	for _, id := range ids {
		if p.literal(pos, id) {
			out = append(out, pos+len(id))
		}
	}
	return out
}

func (p *parser) licenseRef(pos int) []int {
	var starts []int
	starts = append(starts, pos)
	if p.literal(pos, "DocumentRef-") {
		for _, e := range p.idstring(pos + len("DocumentRef-")) {
			if e < len(p.s) && p.s[e] == ':' {
				starts = append(starts, e+1)
			}
		}
	}
	var out []int
	for _, st := range starts {
		if p.literal(st, "LicenseRef-") {
			out = append(out, p.idstring(st+len("LicenseRef-"))...)
		}
	}
	return out
}

func (p *parser) simpleExpression(pos int) set {
	out := set{}
	for _, e := range p.fromList(pos, licenseIDs) {
		out[e] = true
		if e < len(p.s) && p.s[e] == '+' {
			out[e+1] = true
		}
	}
	for _, e := range p.licenseRef(pos) {
		out[e] = true
	}
	return out
}

func (p *parser) compoundHead(pos int) set {
	out := set{}
	for e := range p.simpleExpression(pos) {
		out[e] = true
		for _, s1 := range p.spaces(e) {
			if !p.literal(s1, "WITH") {
				continue
			}
			for _, s2 := range p.spaces(s1 + 4) {
				for _, x := range p.fromList(s2, exceptionIDs) {
					out[x] = true
				}
			}
		}
	}
	if pos < len(p.s) && p.s[pos] == '(' {
		for _, s1 := range p.spacesStar(pos + 1) {
			for e := range p.compoundExpression(s1) {
				for _, s2 := range p.spacesStar(e) {
					if s2 < len(p.s) && p.s[s2] == ')' {
						out[s2+1] = true
					}
				}
			}
		}
	}
	return out
}

func (p *parser) compoundExpression(pos int) set {
	if p.memo == nil {
		p.memo = map[[2]int]set{}
	}
	key := [2]int{pos, 0}
	if r, ok := p.memo[key]; ok {
		return r
	}
	p.memo[key] = set{} // guards left recursion (not present in the grammar)
	out := set{}
	for e := range p.compoundHead(pos) {
		out[e] = true
		for _, s1 := range p.spaces(e) {
			for _, op := range []string{"AND", "OR"} {
				if !p.literal(s1, op) {
					continue
				}
				for _, s2 := range p.spaces(s1 + len(op)) {
					for x := range p.compoundExpression(s2) {
						out[x] = true
					}
				}
			}
		}
	}
	p.memo[key] = out
	return out
}

func (p *parser) licenseExpression(pos int) set {
	out := set{}
	for e := range p.compoundExpression(pos) {
		out[e] = true
	}
	for e := range p.simpleExpression(pos) {
		out[e] = true
	}
	return out
}
