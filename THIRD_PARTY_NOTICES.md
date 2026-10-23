# Third-party notices

gocomposer is a Go re-implementation of Composer. Behaviour, messages,
data files and test fixtures were ported from the projects below at the
versions pinned by Composer's own `composer.lock` at the reference commit.
Each ported Go file names its original source file(s) in its header.

## Composer

Upstream: https://github.com/composer/composer
Commit: f83588392df43d4fd027d84d26de459946ef4e80
Copyright (c) Nils Adermann, Jordi Boggiano
License: MIT — `licenses/composer-MIT.txt`.

Ported: solver, pool builder/optimizer, policies, repositories, installer,
locker, downloaders, autoload generator, event dispatcher, commands,
validators, policy layer (advisories, filter lists, cooldown), JsonFile and
JsonManipulator. Verbatim copies: `res/composer-schema.json`,
`res/composer-lock-schema.json` (internal/jsonschema/res),
`src/Composer/Autoload/ClassLoader.php`, `src/Composer/InstalledVersions.php`
(internal/res), upstream test fixtures under `testdata/` and
`internal/*/testdata` (each directory has a NOTICE).

## Composer dependencies (versions from Composer's composer.lock)

| Project | Version | Used for | License |
|---|---|---|---|
| composer/semver | 3.5.1 | version parsing/constraints/intervals (internal/semver) | MIT, `licenses/composer-semver-MIT.txt` |
| composer/class-map-generator | 1.8.0 | class map generation (internal/autoload) | MIT, `licenses/composer-class-map-generator-MIT.txt` |
| composer/spdx-licenses | 1.6.0 | SPDX data files (verbatim) and validation (internal/spdx) | MIT, `licenses/composer-spdx-licenses-MIT.txt` |
| composer/pcre | 2.3.2 | Preg wrapper semantics (internal/pcre, internal/jsonmanip) | MIT, `licenses/composer-pcre-MIT.txt` |
| seld/jsonlint | 1.12.1 | JSON lint parser and messages (internal/jsonlint) | MIT, `licenses/seld-jsonlint-MIT.txt` |
| justinrainbow/json-schema | 6.13.1 | JSON schema validation (internal/jsonschema) | MIT, `licenses/justinrainbow-json-schema-MIT.txt` |
| symfony/console | 5.4.54 | console behaviour: option parsing, exception rendering, Table, ProgressBar, output formatting (internal/cli, internal/cio) | MIT, `licenses/symfony-console-MIT.txt` |

## PHP

The FILTER_VALIDATE_EMAIL pattern (php-src `ext/filter/logical_filters.c`,
pattern by Michael Rushton) is reproduced in internal/validate and
internal/jsonschema. PHP License 3.01 — `licenses/php-3.01.txt`.
This product includes PHP software, freely available from
<http://www.php.net/software/>.

Official PHP and Composer are only executed as the comparison oracle; the
Composer application itself is not bundled.

## Plugin runtime (internal/res/pluginrt, embedded in the binary)

Composer plugins are PHP code that type-hint Composer's API, so gocomposer
runs them in a PHP process with that API available. The process only hosts
plugin code; resolution, locking, downloads, extraction and autoload
generation stay in Go, and the PHP side calls back into Go through stubs.

- `internal/res/pluginrt/composer/`: Composer API classes (interfaces,
  events, IO, Package, Json, Installer base classes) copied unmodified from
  composer/composer f83588392df43d4fd027d84d26de459946ef4e80, plus
  gocomposer-written stubs (Composer, Config, Factory, EventDispatcher, …)
  whose header names the Composer file they stand in for. MIT,
  (c) Nils Adermann, Jordi Boggiano — `internal/res/pluginrt/composer/LICENSE`.
- `internal/res/pluginrt/vendor/`: the PHP libraries Composer ships with
  plugins at the versions locked in the oracle's `composer.lock`
  (symfony/console, process, finder, filesystem, string, service-contracts,
  deprecation-contracts and polyfills; composer/semver, composer/pcre,
  composer/class-map-generator; seld/jsonlint, seld/signal-handler; psr/log;
  react/promise), each with its own LICENSE file (MIT).

## Test fixtures from other projects

- `testdata/compare/packages/symfony/flex/2.11.0/`: symfony/flex v2.11.0
  (https://github.com/symfony/flex, commit
  4a6d98eea3ebc7f68d82810cb682eedca2649e99), (c) Fabien Potencier, MIT —
  its own `LICENSE` next to the sources (see `FIXTURE-ORIGIN.txt`). Used as a
  fixture plugin by the flex-* comparison scenarios.
- `testdata/integration/symfony-skeleton/composer.json`: derived from
  symfony/skeleton v7.4.99 (commit bf1abda299403468285fd4fa0514d916991c9120),
  MIT (see `SOURCE.txt`).
