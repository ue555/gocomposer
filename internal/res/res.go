// Package res embeds PHP runtime files that Composer copies into
// vendor/composer. They are copied unchanged from Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (src/Composer/InstalledVersions.php
// and src/Composer/Autoload/ClassLoader.php), Copyright (c) Nils Adermann,
// Jordi Boggiano, MIT license (see LICENSE in this directory).
package res

import _ "embed"

//go:embed InstalledVersions.php
var InstalledVersions []byte

//go:embed ClassLoader.php
var ClassLoader []byte

//go:embed LICENSE
var License []byte
