package validate

// Port of Composer\Util\ConfigValidator (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80, MIT, (c) Nils Adermann, Jordi Boggiano).

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gocomposer/internal/jsonfile"
	"gocomposer/internal/jsonlint"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/spdx"
)

// Schema modes of JsonFile.
const (
	LaxSchema    = 1
	StrictSchema = 2
	AuthSchema   = 3
	LockSchema   = 4
)

// CheckVersion is ConfigValidator::CHECK_VERSION.
const CheckVersion = 1

// JSONValidationError is Composer\Json\JsonValidationException.
type JSONValidationError struct {
	Msg    string
	Errors []string
}

func (e *JSONValidationError) Error() string { return e.Msg }

// SchemaValidator validates a JSON file against a schema mode like
// JsonFile::validateSchema, returning *JSONValidationError on schema
// errors and other errors for unreadable/unparsable files.
type SchemaValidator func(path string, schema int) error

var reUpper = regexp.MustCompile(`[A-Z]`)
var reSuggest = regexp.MustCompile(`([a-z])([A-Z])|([A-Z])([A-Z][a-z])`)

// SuggestName is the lower-cased/with-dashes suggestion of ConfigValidator.
func SuggestName(name string) string {
	return strings.ToLower(reSuggest.ReplaceAllString(name, "${1}${3}-${2}${4}"))
}

// ValidateConfig is ConfigValidator::validate. fatal is an exception the
// PHP code does not catch (it aborts the command).
func ValidateConfig(file string, arrayLoaderFlags, flags int, schema SchemaValidator) (errs, publishErrors, warnings []string, fatal error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return []string{"Could not read " + file + "\n\n" + err.Error()}, nil, nil, nil
	}
	value, err := jsonfile.Parse(raw, file)
	if err != nil {
		return []string{err.Error()}, nil, nil, nil
	}
	laxValid := false
	if schema != nil {
		err := schema(file, LaxSchema)
		if err == nil {
			laxValid = true
			err = schema(file, StrictSchema)
		}
		var ve *JSONValidationError
		if errors.As(err, &ve) {
			for _, m := range ve.Errors {
				if laxValid {
					publishErrors = append(publishErrors, m)
				} else {
					errs = append(errs, m)
				}
			}
		} else if err != nil {
			return append(errs, err.Error()), publishErrors, warnings, nil
		}
	}
	manifest := phpjson.ToMap(value)
	if phpjson.IsArray(value) {
		if derr := jsonlint.Lint(string(raw), jsonlint.DetectKeyConflicts); derr != nil {
			var de *jsonlint.DuplicateKeyError
			if errors.As(derr, &de) {
				warnings = append(warnings, fmt.Sprintf("Key %s is a duplicate in %s at line %d", de.Key, file, de.Line))
			}
		}
	}
	if manifest == nil {
		manifest = phpjson.NewMap()
	}
	get := func(k string) phpjson.Value { v, _ := manifest.Get(k); return v }

	if !phpjson.Truthy(get("license")) {
		warnings = append(warnings, `No license specified, it is recommended to do so. For closed-source software you may use "proprietary" as license.`)
	} else {
		var licenses []string
		eachCast(get("license"), func(_ string, v phpjson.Value) {
			if s, ok := v.(string); ok && s != "proprietary" {
				licenses = append(licenses, s)
			}
		})
		for _, license := range licenses {
			info := spdx.LicenseByIdentifier(license)
			if info == nil || !info.Deprecated {
				continue
			}
			switch {
			case reGPLPlus.MatchString(license):
				warnings = append(warnings, fmt.Sprintf(`License "%s" is a deprecated SPDX license identifier, use "`+strings.ReplaceAll(license, "+", "")+`-or-later" instead`, license))
			case reGPLBare.MatchString(license):
				warnings = append(warnings, fmt.Sprintf(`License "%s" is a deprecated SPDX license identifier, use "`+license+`-only" or "`+license+`-or-later" instead`, license))
			default:
				warnings = append(warnings, fmt.Sprintf(`License "%s" is a deprecated SPDX license identifier, see https://spdx.org/licenses/`, license))
			}
		}
	}
	if flags&CheckVersion != 0 && manifest.Isset("version") {
		warnings = append(warnings, "The version field is present, it is recommended to leave it out if the package is published on Packagist.")
	}
	if name, ok := get("name").(string); ok && phpjson.Truthy(name) && reUpper.MatchString(name) {
		publishErrors = append(publishErrors, fmt.Sprintf(`Name "%s" does not match the best practice (e.g. lower-cased/with-dashes). We suggest using "%s" instead. As such you will not be able to submit it to Packagist.`, name, SuggestName(name)))
	}
	if t, ok := get("type").(string); ok && t == "composer-installer" {
		warnings = append(warnings, "The package type 'composer-installer' is deprecated. Please distribute your custom installers as plugins from now on. See https://getcomposer.org/doc/articles/plugins.md for plugin documentation.")
	}
	if manifest.Isset("require") && manifest.Isset("require-dev") {
		dev := phpjson.ToMap(get("require-dev"))
		var overrides []string
		phpjson.Each(get("require"), func(k string, _ phpjson.Value) {
			if dev != nil && dev.Has(k) {
				overrides = append(overrides, k)
			}
		})
		if len(overrides) > 0 {
			plural := "is"
			if len(overrides) > 1 {
				plural = "are"
			}
			warnings = append(warnings, strings.Join(overrides, ", ")+" "+plural+" required both in require and require-dev, this can lead to unexpected behavior")
		}
	}
	for _, linkType := range []string{"provide", "replace"} {
		if !manifest.Isset(linkType) {
			continue
		}
		for _, requireType := range []string{"require", "require-dev"} {
			if !manifest.Isset(requireType) {
				continue
			}
			req := phpjson.ToMap(get(requireType))
			phpjson.Each(get(linkType), func(provide string, _ phpjson.Value) {
				if req != nil && req.Isset(provide) {
					warnings = append(warnings, "The package "+provide+" in "+requireType+" is also listed in "+linkType+" which satisfies the requirement. Remove it from "+linkType+" if you wish to install it.")
				}
			})
		}
	}
	// array_merge($require, $requireDev)
	merged := phpjson.NewMap()
	var order []string
	for _, k := range []string{"require", "require-dev"} {
		phpjson.Each(get(k), func(name string, v phpjson.Value) {
			if !merged.Has(name) {
				order = append(order, name)
			}
			merged.Set(name, v)
		})
	}
	for _, name := range order {
		if v, ok := merged.Get(name); ok {
			if s, ok := v.(string); ok && strings.Contains(s, "#") {
				warnings = append(warnings, fmt.Sprintf(`The package "%s" is pointing to a commit-ref, this is bad practice and can cause unforeseen issues.`, name))
			}
		}
	}
	scripts := phpjson.ToMap(get("scripts"))
	phpjson.Each(get("scripts-descriptions"), func(name string, _ phpjson.Value) {
		if scripts == nil || !scripts.Has(name) {
			warnings = append(warnings, fmt.Sprintf(`Description for non-existent script "%s" found in "scripts-descriptions"`, name))
		}
	})
	phpjson.Each(get("scripts-aliases"), func(name string, _ phpjson.Value) {
		if scripts == nil || !scripts.Has(name) {
			warnings = append(warnings, fmt.Sprintf(`Aliases for non-existent script "%s" found in "scripts-aliases"`, name))
		}
	})
	if al := manifest.Map("autoload"); al != nil {
		if m := al.Map("psr-0"); m != nil && m.Isset("") {
			warnings = append(warnings, "Defining autoload.psr-0 with an empty namespace prefix is a bad idea for performance")
		}
		if m := al.Map("psr-4"); m != nil && m.Isset("") {
			warnings = append(warnings, "Defining autoload.psr-4 with an empty namespace prefix is a bad idea for performance")
		}
	}

	loader := &ArrayLoader{Flags: arrayLoaderFlags}
	m := manifest.Clone()
	if !m.Isset("version") {
		m.Set("version", "1.0.0")
	}
	if !m.Isset("name") {
		m.Set("name", "dummy/dummy")
	}
	if _, err := loader.Load(m); err != nil {
		var ie *InvalidPackageError
		if !errors.As(err, &ie) {
			return errs, publishErrors, warnings, err
		}
		errs = append(errs, ie.Errors...)
	}
	warnings = append(warnings, loader.Warnings()...)
	return errs, publishErrors, warnings, nil
}
