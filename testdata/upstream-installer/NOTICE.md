# Upstream installer fixtures

Copied unchanged from Composer `tests/Composer/Test/Fixtures/installer` at commit
f83588392df43d4fd027d84d26de459946ef4e80 (https://github.com/composer/composer).
Copyright (c) Nils Adermann, Jordi Boggiano. MIT license, see licenses/composer-MIT.txt.
Expectations are not edited; fixtures needing unsupported features are reported as skipped with a reason.

Also copied unchanged, as the fixtures reference them through relative path
repositories (InstallerTest runs from `tests/Composer/Test`):
- `Fixtures/functional/installed-versions/` from `tests/Composer/Test/Fixtures/functional/installed-versions/`
- `DependencyResolver/Fixtures/poolbuilder/{symlinked,mirrored}-path-repo/` from `tests/Composer/Test/DependencyResolver/Fixtures/poolbuilder/`
- `install-branch-alias-composer-repo/`, `update-policy-cooldown-composer-repo/` (composer repositories used via file:// URLs)

The runner (internal/composer/installer_fixture_test.go) ports
InstallerTest::doTestIntegration and its mocks. The PHP platform comes from a
snapshot recorded with the oracle PHP 8.4.27 (internal/composer/testdata/php84-facts.json),
with ext-intl added because Composer's own test environment has it loaded.
