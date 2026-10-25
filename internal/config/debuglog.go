package config

// Debug output and deprecation notices of Composer\Factory and
// Composer\Config (composer/composer f83588392df43d4fd027d84d26de459946ef4e80,
// MIT, (c) Nils Adermann, Jordi Boggiano).

// Deprecation receives the E_USER_DEPRECATED notices Config::get triggers
// (Composer's ErrorHandler prints them); nil ignores them.
var Deprecation func(msg string)

// SetDebug sets the receiver of the IOInterface::DEBUG messages written by
// Factory::loadComposerAuthEnv for this config (nil: none).
func (c *Config) SetDebug(debug func(string)) { c.debug = debug }

func (c *Config) debugf(msg string) {
	if c.debug != nil {
		c.debug(msg)
	}
}
