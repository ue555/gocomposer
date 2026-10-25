package main

// show --all / --available without composer.json: Composer lists the default
// repositories (RepositoryFactory::defaultReposWithDefaultManager), here the
// loopback packagist.org stub of the cli3-packagist-* scenarios.
func finalShowScenarios() []Scenario {
	fx := []string{"packagiststub"}
	env := map[string]string{"SSL_CERT_FILE": "/fixtures/git/_api/ca.pem"}
	return []Scenario{
		{Name: "show3-packagist-available", NoJSON: true, GitFixtures: fx, Env: env, Command: "show --available"},
		{Name: "show3-packagist-available-name", NoJSON: true, GitFixtures: fx, Env: env, Command: "show --available stubvendor/alpha"},
		{Name: "show3-packagist-all", NoJSON: true, GitFixtures: fx, Env: env, Command: "show --all"},
		{Name: "show3-packagist-all-name", NoJSON: true, GitFixtures: fx, Env: env, Command: "show --all stubvendor/beta"},
		{Name: "show3-packagist-all-names-only", NoJSON: true, GitFixtures: fx, Env: env, Command: "show -a -N"},
		{Name: "show3-packagist-all-json", NoJSON: true, GitFixtures: fx, Env: env, Command: "show -a --format=json stubvendor/alpha"},
		{Name: "show3-packagist-all-missing", NoJSON: true, GitFixtures: fx, Env: env, Command: "show -a stubvendor/none"},
	}
}

func init() { extraScenarios = append(extraScenarios, finalShowScenarios) }
