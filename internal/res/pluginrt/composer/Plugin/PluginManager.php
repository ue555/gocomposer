<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Plugin\PluginManager as seen by
// plugins. Loading, allow-plugins checks and registration are implemented by
// gocomposer in Go (port of Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Plugin/PluginManager.php, MIT license, (c) Nils Adermann,
// Jordi Boggiano); the capability helpers below are ported unchanged.

namespace Composer\Plugin;

use Composer\Plugin\Capability\Capability;
use GocomposerPluginHost\Host;
use GocomposerPluginHost\Unsupported;

class PluginManager
{
    /**
     * @return array<PluginInterface>
     */
    public function getPlugins(): array
    {
        return Host::$plugins;
    }

    /**
     * @return array<string>
     */
    public function getRegisteredPlugins(): array
    {
        return array_keys(Host::$registered);
    }

    public function getGlobalComposer()
    {
        return null;
    }

    public function setRunningInGlobalDir(bool $runningInGlobalDir): void
    {
    }

    public function loadInstalledPlugins(): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function deactivateInstalledPlugins(): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function registerPackage($package, bool $failOnMissingClasses = false, bool $isGlobalPlugin = false): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function deactivatePackage($package): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function uninstallPackage($package): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function addPlugin(PluginInterface $plugin, bool $isGlobalPlugin = false, $sourcePackage = null): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function removePlugin(PluginInterface $plugin): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function uninstallPlugin(PluginInterface $plugin): void
    {
        Unsupported::fail(__METHOD__);
    }

    /**
     * @param 'local'|'global' $type
     * @return bool
     */
    public function arePluginsDisabled($type)
    {
        return false;
    }

    public function disablePlugins(): void
    {
        Unsupported::fail(__METHOD__);
    }

    public function isPluginAllowed(string $package, bool $isGlobalPlugin, bool $optional = false, bool $prompt = true): bool
    {
        return (bool) Host::cb('isPluginAllowed', ['name' => $package, 'global' => $isGlobalPlugin, 'optional' => $optional]);
    }

    /**
     * @throws \RuntimeException On empty or non-string implementation class name value
     * @return null|string       The fully qualified class of the implementation or null if Plugin is not of Capable type or does not provide it
     */
    protected function getCapabilityImplementationClassName(PluginInterface $plugin, string $capability): ?string
    {
        if (!($plugin instanceof Capable)) {
            return null;
        }

        $capabilities = (array) $plugin->getCapabilities();

        if (!empty($capabilities[$capability]) && is_string($capabilities[$capability]) && trim($capabilities[$capability])) {
            return trim($capabilities[$capability]);
        }

        if (
            array_key_exists($capability, $capabilities)
            && (empty($capabilities[$capability]) || !is_string($capabilities[$capability]) || !trim($capabilities[$capability]))
        ) {
            throw new \UnexpectedValueException('Plugin '.get_class($plugin).' provided invalid capability class name(s), got '.var_export($capabilities[$capability], true));
        }

        return null;
    }

    /**
     * @param  array<mixed>                  $ctorArgs
     */
    public function getPluginCapability(PluginInterface $plugin, $capabilityClassName, array $ctorArgs = []): ?Capability
    {
        if ($capabilityClass = $this->getCapabilityImplementationClassName($plugin, $capabilityClassName)) {
            if (!class_exists($capabilityClass)) {
                throw new \RuntimeException("Cannot instantiate Capability, as class $capabilityClass from plugin ".get_class($plugin)." does not exist.");
            }

            if (!is_a($capabilityClass, Capability::class, true) || !is_a($capabilityClass, $capabilityClassName, true)) {
                throw new \RuntimeException(
                    'Class ' . $capabilityClass . ' must implement both Composer\Plugin\Capability\Capability and '. $capabilityClassName . '.'
                );
            }

            $ctorArgs['plugin'] = $plugin;

            return new $capabilityClass($ctorArgs);
        }

        return null;
    }

    /**
     * @param  array<mixed>                  $ctorArgs
     * @return Capability[]
     */
    public function getPluginCapabilities($capabilityClassName, array $ctorArgs = []): array
    {
        $capabilities = [];
        foreach ($this->getPlugins() as $plugin) {
            $capability = $this->getPluginCapability($plugin, $capabilityClassName, $ctorArgs);
            if (null !== $capability) {
                $capabilities[] = $capability;
            }
        }

        return $capabilities;
    }
}
