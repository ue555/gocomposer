package composer

// ArchiveManager for plugins (Composer::getArchiveManager): the runtime has
// Composer's ArchiveManager (internal/res/pluginrt/composer/Package/
// Archiver), whose built-in archivers and package downloads are
// gocomposer's. When plugins change the archive manager of the running
// instance (setArchiveManager, addArchiver), the archive command archives
// through it like Composer's ArchiveCommand does with
// $composer->getArchiveManager() (Composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT license, (c) Nils Adermann,
// Jordi Boggiano).

import (
	"errors"

	"gocomposer/internal/archiver"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
)

// archiveCallback serves the archive callbacks of the runtime (ok is false
// for other ops).
func (pm *PluginManager) archiveCallback(op string, a *phpjson.Map) (phpjson.Value, bool, error) {
	switch op {
	case "archiveSources":
		p, err := pm.host.pkgFromRef(mustGet(a, "p"))
		if err != nil {
			return nil, true, err
		}
		return nil, true, NewInstallationManager(pm.c).DownloadTo(p, a.Str("path"))
	case "archiveRun":
		var excludes []string
		phpjson.Each(mustGet(a, "excludes"), func(_ string, v phpjson.Value) {
			if s, ok := v.(string); ok {
				excludes = append(excludes, s)
			}
		})
		path, err := archiver.Archive(a.Str("sources"), a.Str("target"), a.Str("format"), excludes, mustGet(a, "ignoreFilters") == true)
		if err != nil {
			var ae *archiver.ArchiveError
			var se *archiver.SourceError
			var pe *archiver.PatternError
			if errors.As(err, &ae) || errors.As(err, &se) || errors.As(err, &pe) {
				return nil, true, &UserError{Msg: err.Error(), Code: 1}
			}
			return nil, true, err
		}
		return path, true, nil
	}
	return nil, false, nil
}

// PluginArchive is ArchiveCommand's call of $composer->getArchiveManager()
// ->archive() when plugins changed the archive manager of c (handled is
// false otherwise: gocomposer archives itself).
func PluginArchive(c *Composer, p *pkg.Package, format, dest, fileName string, ignoreFilters bool) (path string, handled bool, err error) {
	pm := pluginManagerOf(c)
	if !pm.active() {
		return "", false, nil
	}
	v, err := pm.host.call("archiveChanged", nil)
	if err != nil {
		return "", true, err
	}
	if v != true {
		return "", false, nil
	}
	args := phpjson.NewMap()
	args.Set("p", pm.host.ref(p))
	args.Set("format", format)
	args.Set("dest", dest)
	if fileName != "" {
		args.Set("fileName", fileName)
	}
	args.Set("ignoreFilters", ignoreFilters)
	args.Set("frames", pm.frames())
	v, err = pm.host.call("archiveVia", args)
	if err == nil {
		err = pm.takePendingErr()
	}
	if err != nil {
		return "", true, err
	}
	s, _ := v.(string)
	return s, true, nil
}
