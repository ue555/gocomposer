// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Repository/RepositoryFactory.php (defaultRepos keys)
// (MIT, (c) Nils Adermann, Jordi Boggiano).

package composer

import "gocomposer/internal/phpjson"

// DefaultRepositoryNames are the keys of RepositoryFactory::defaultRepos
// (the names generated like createRepositories does), in order.
func (c *Composer) DefaultRepositoryNames() []string {
	var out []string
	names := map[string]bool{}
	repos := c.Config.Repositories()
	for _, index := range repos.Keys() {
		rv, _ := repos.Get(index)
		rc := phpjson.ToMap(rv)
		name := index
		if _, isInt := phpjson.IntKey(index); isInt && rc != nil && rc.Isset("url") {
			name = reHTTPPrefix.ReplaceAllString(rc.Str("url"), "")
		}
		for names[name] {
			name += "2"
		}
		names[name] = true
		out = append(out, name)
	}
	return out
}
