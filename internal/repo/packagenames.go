package repo

// Port of ComposerRepository::getPackageNames (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"regexp"
	"strings"
)

// PackageNames is ComposerRepository::getPackageNames($packageFilter).
// Repositories that only expose their names through a "list" API or through
// v1 providers are not supported and fail explicitly.
func (r *Composer) PackageNames(filter string) ([]string, error) {
	data, err := r.loadRootServerFile()
	if err != nil {
		return nil, err
	}
	apply := func(names []string) []string { return names }
	if filter != "" {
		// BasePackage::packageNameToRegexp
		re := regexp.MustCompile(`(?i)^` + strings.ReplaceAll(regexp.QuoteMeta(filter), `\*`, `.*`) + `$`)
		apply = func(names []string) []string {
			var out []string
			for _, n := range names {
				if re.MatchString(n) {
					out = append(out, n)
				}
			}
			return out
		}
	}
	if r.lazyProvidersURL != "" {
		if r.availablePackages != nil {
			return apply(r.availableNames), nil
		}
		if data.Has("list") {
			return nil, &UnsupportedError{"Listing package names through the repository list API is not supported by gocomposer yet."}
		}
		hasPartial, err := r.partialPackages()
		if err != nil {
			return nil, err
		}
		if hasPartial && r.partialByName != nil {
			return apply(append([]string(nil), r.partialOrder...)), nil
		}
		return nil, nil
	}
	if r.hasProviders {
		return nil, &UnsupportedError{"Listing package names of provider-based (v1 providers) Composer repositories is not supported by gocomposer yet."}
	}
	ps, err := r.Packages()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(ps))
	for _, p := range ps {
		names = append(names, p.PrettyName)
	}
	return apply(names), nil
}
