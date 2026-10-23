package cli

import (
	"regexp"
	"sort"
	"strings"

	"gocomposer/internal/semver"
)

func regexpQuote(s string) string { return regexp.QuoteMeta(s) }

func regexpCompile(s string) (*regexp.Regexp, error) { return regexp.Compile(s) }

// sortNaturalFold is ksort(SORT_NATURAL | SORT_FLAG_CASE).
func sortNaturalFold(list []string) {
	sort.SliceStable(list, func(i, j int) bool {
		return semver.NatCompare(strings.ToLower(list[i]), strings.ToLower(list[j])) < 0
	})
}
