package res

import "embed"

// PluginRuntime is the PHP side of gocomposer's Composer plugin support
// (pluginrt/host.php and the Composer plugin API classes it provides to
// plugins). Files under pluginrt/composer are copied or ported from Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (MIT, (c) Nils Adermann, Jordi
// Boggiano; provenance in each file header); pluginrt/vendor holds the
// libraries Composer bundles for plugins, copied unchanged from the versions
// locked by that Composer commit, each with its own LICENSE file.
//
//go:embed all:pluginrt
var PluginRuntime embed.FS
