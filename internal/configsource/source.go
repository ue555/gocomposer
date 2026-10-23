// Package configsource ports Composer\Config\JsonConfigSource
// (src/Composer/Config/JsonConfigSource.php of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80): edits composer.json /
// config.json / auth.json through jsonmanip and, when the in-place edit is
// not possible, rewrites the whole file from decoded data.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// Errors that PHP throws are returned as *jsonmanip.Exception (Class and
// Message as in PHP); I/O failures that PHP only warns about are returned as
// plain errors instead of being ignored.
package configsource

import (
	"fmt"
	"os"
	"strings"

	"gocomposer/internal/jsonmanip"
	"gocomposer/internal/pcre"
	"gocomposer/internal/phpjson"
)

// Source is a JsonConfigSource.
type Source struct {
	file       *jsonFile
	authConfig bool
	// ValidateSchema is JsonFile::validateSchema(LAX_SCHEMA) on the written
	// file. Returning *ValidationError makes the update revert.
	ValidateSchema func(path string) error
}

// New is new JsonConfigSource(new JsonFile($path), $authConfig).
func New(path string, authConfig bool) *Source {
	return &Source{file: &jsonFile{path: path, indent: indentDefault}, authConfig: authConfig, ValidateSchema: validateSchemaLax}
}

// Name is getName(): the file path.
func (s *Source) Name() string { return s.file.path }

type fallback func(config *phpjson.Map, args []phpjson.Value)

// AddRepository is addRepository($name, $config, $append); config is an
// array (*phpjson.Map / list) or false.
func (s *Source) AddRepository(name string, config phpjson.Value, appendItem bool) error {
	return s.manipulateJSON("addRepository", func(cfg *phpjson.Map, a []phpjson.Value) {
		repo := a[0].(string)
		repoConfig := normalize(a[1])
		convertRepositoriesToList(cfg)

		if b, ok := repoConfig.(bool); ok && !b {
			if isset(cfg, "repositories") {
				repos := requireArrayArg("foreach", 1, "array", get(cfg, "repositories"))
				for _, k := range repos.Keys() {
					repository, _ := repos.Get(k)
					if n, ok := get(repository, "name").(string); ok && n == repo {
						repos.Set(k, mapOf(repo, repoConfig))
						return
					}
					if isFalseMap(repository, repo) {
						return
					}
				}
			} else {
				cfg.Set("repositories", phpjson.NewMap())
			}
			push(slot(cfg, "repositories"), mapOf(repo, repoConfig))
			return
		}

		if m := arr(repoConfig); m != nil && repo != "" && !isset(m, "name") {
			repoConfig = prependName(repo, m)
		}

		// ensure uniqueness by removing any existing entries which use the same name
		repos := filterRepos(get(cfg, "repositories"), repo)
		if a[2].(bool) {
			push(repos, repoConfig)
		} else {
			repos = listMap(append([]phpjson.Value{repoConfig}, values(repos)...))
		}
		cfg.Set("repositories", repos)
	}, name, config, appendItem)
}

// convertRepositoriesToList converts an associative repositories map to a
// list, keeping the name (the shared preamble of the repository fallbacks).
func convertRepositoriesToList(cfg *phpjson.Map) {
	reposV := get(cfg, "repositories")
	if reposV == nil {
		return // array_is_list([]) is true
	}
	repos := requireArrayArg("array_is_list", 1, "array", reposV)
	if repos.IsList() {
		return
	}
	var list []phpjson.Value
	for _, idx := range repos.Keys() {
		repository, _ := repos.Get(idx)
		switch {
		case isStringKey(idx) && isArr(repository):
			// convert to list entry with name
			if !isset(repository, "name") {
				repository = prependName(idx, arr(repository))
			}
			list = append(list, repository)
		case isStringKey(idx):
			// keep boolean entries (e.g. 'packagist.org' => false)
			list = append(list, mapOf(idx, repository))
		default:
			list = append(list, repository)
		}
	}
	cfg.Set("repositories", listMap(list))
}

// InsertRepository is insertRepository($name, $config, $referenceName, $offset).
func (s *Source) InsertRepository(name string, config phpjson.Value, referenceName string, offset int) error {
	return s.manipulateJSON("insertRepository", func(cfg *phpjson.Map, a []phpjson.Value) {
		name := a[0].(string)
		repoConfig := normalize(a[1])
		referenceName := a[2].(string)
		offset := int(a[3].(int64))
		convertRepositoriesToList(cfg)

		// ensure uniqueness by removing any existing entries which use the same name
		repos := filterRepos(get(cfg, "repositories"), name)
		cfg.Set("repositories", repos)

		indexToInsert := -1
		for i, repository := range values(repos) {
			if n, ok := get(repository, "name").(string); ok && n == referenceName {
				indexToInsert = i
				break
			}
			if isFalseMap(repository, referenceName) {
				indexToInsert = i
				break
			}
		}

		if indexToInsert < 0 {
			phpError("RuntimeException", fmt.Sprintf("The referenced repository \"%s\" does not exist.", referenceName))
		}

		if m := arr(repoConfig); m != nil && name != "" && !isset(m, "name") {
			repoConfig = prependName(name, m)
		}

		// array_splice($config['repositories'], $indexToInsert + $offset, 0, [$repoConfig])
		list := values(repos)
		pos := indexToInsert + offset
		if pos < 0 {
			pos += len(list)
			if pos < 0 {
				pos = 0
			}
		}
		if pos > len(list) {
			pos = len(list)
		}
		list = append(list[:pos], append([]phpjson.Value{repoConfig}, list[pos:]...)...)
		cfg.Set("repositories", listMap(list))
	}, name, config, referenceName, int64(offset))
}

// SetRepositoryUrl is setRepositoryUrl($name, $url).
func (s *Source) SetRepositoryUrl(name, url string) error {
	return s.manipulateJSON("setRepositoryUrl", func(cfg *phpjson.Map, a []phpjson.Value) {
		name, url := a[0].(string), a[1].(string)
		repos := arr(get(cfg, "repositories"))
		if repos == nil {
			return
		}
		for _, index := range repos.Keys() {
			repository, _ := repos.Get(index)
			if isStringKey(index) && name == index {
				setPath(cfg, url, "repositories", index, "url")
				return
			}
			if n, ok := get(repository, "name").(string); ok && name == n {
				setPath(cfg, url, "repositories", index, "url")
				return
			}
		}
	}, name, url)
}

// RemoveRepository is removeRepository($name).
func (s *Source) RemoveRepository(name string) error {
	return s.manipulateJSON("removeRepository", func(cfg *phpjson.Map, a []phpjson.Value) {
		repo := a[0].(string)
		if isset(cfg, "repositories", repo) {
			unsetPath(cfg, "repositories", repo)
		} else {
			cfg.Set("repositories", filterRepos(get(cfg, "repositories"), repo))
		}

		if v, _ := cfg.Get("repositories"); isEmptyArray(v) {
			cfg.Delete("repositories")
		}
	}, name)
}

const authKeyPattern = `{^(bitbucket-oauth|github-oauth|gitlab-oauth|gitlab-token|bearer|http-basic|custom-headers|forgejo-token|platform)\.}`

func isAuthKey(key string) bool {
	m, err := pcre.IsMatch(authKeyPattern, key)
	if err != nil {
		panic(err)
	}
	return m != nil
}

// AddConfigSetting is addConfigSetting($name, $value).
func (s *Source) AddConfigSetting(name string, value phpjson.Value) error {
	authConfig := s.authConfig
	return s.manipulateJSON("addConfigSetting", func(cfg *phpjson.Map, a []phpjson.Value) {
		key, val := a[0].(string), normalize(a[1])
		if isAuthKey(key) {
			k, host, _ := strings.Cut(key, ".")
			if authConfig {
				setPath(cfg, val, k, host)
			} else {
				setPath(cfg, val, "config", k, host)
			}
		} else if strings.HasPrefix(key, "policy.") {
			if !isset(cfg, "config") || !isArr(get(cfg, "config")) {
				cfg.Set("config", phpjson.NewMap())
			}

			bits := strings.Split(key, ".")
			last := bits[len(bits)-1]
			bits = bits[:len(bits)-1]
			cur := arr(get(cfg, "config"))
			for _, bit := range bits {
				if !isset(cur, bit) || !isArr(get(cur, bit)) {
					cur.Set(bit, phpjson.NewMap())
				}
				cur = arr(get(cur, bit))
			}
			cur.Set(last, val)
		} else {
			setPath(cfg, val, "config", key)
		}
	}, name, value)
}

// RemoveConfigSetting is removeConfigSetting($name).
func (s *Source) RemoveConfigSetting(name string) error {
	authConfig := s.authConfig
	return s.manipulateJSON("removeConfigSetting", func(cfg *phpjson.Map, a []phpjson.Value) {
		key := a[0].(string)
		if isAuthKey(key) {
			k, host, _ := strings.Cut(key, ".")
			if authConfig {
				unsetPath(cfg, k, host)
			} else {
				unsetPath(cfg, "config", k, host)
			}
		} else if strings.HasPrefix(key, "policy.") {
			if !isset(cfg, "config") || !isArr(get(cfg, "config")) {
				return
			}
			bits := strings.Split(key, ".")
			last := bits[len(bits)-1]
			bits = bits[:len(bits)-1]
			cur := arr(get(cfg, "config"))
			for _, bit := range bits {
				if !isset(cur, bit) || !isArr(get(cur, bit)) {
					return
				}
				cur = arr(get(cur, bit))
			}
			cur.Delete(last)

			// cascade: drop now-empty ancestors within the policy subtree (stops before config itself)
			for len(bits) > 0 {
				leafKey := bits[len(bits)-1]
				bits = bits[:len(bits)-1]
				var parent phpjson.Value = get(cfg, "config")
				for _, bit := range bits {
					if !isset(parent, bit) {
						return
					}
					parent = get(parent, bit)
				}
				if isset(parent, leafKey) && isEmptyArray(get(parent, leafKey)) {
					unsetPath(parent.(*phpjson.Map), leafKey)
				} else {
					break
				}
			}
		} else {
			unsetPath(cfg, "config", key)
		}
	}, name)
}

// propertyBase resolves "$arr = &$config[reset($bits)]": taking the
// reference creates the key with null when it is missing.
func propertyBase(cfg *phpjson.Map, first string) {
	if _, ok := cfg.Get(first); !ok {
		cfg.Set(first, nil)
	}
}

// AddProperty is addProperty($name, $value).
func (s *Source) AddProperty(name string, value phpjson.Value) error {
	return s.manipulateJSON("addProperty", func(cfg *phpjson.Map, a []phpjson.Value) {
		key, val := a[0].(string), normalize(a[1])
		if strings.HasPrefix(key, "extra.") || strings.HasPrefix(key, "scripts.") {
			bits := strings.Split(key, ".")
			last := bits[len(bits)-1]
			bits = bits[:len(bits)-1]
			propertyBase(cfg, bits[0])
			// $arr starts at $config[reset($bits)] and then walks every bit,
			// including the first one again
			cur := slot(cfg, bits[0])
			for _, bit := range bits {
				if !isset(cur, bit) {
					cur.Set(bit, phpjson.NewMap())
				}
				cur = slot(cur, bit)
			}
			cur.Set(last, val)
		} else {
			cfg.Set(key, val)
		}
	}, name, value)
}

// RemoveProperty is removeProperty($name).
func (s *Source) RemoveProperty(name string) error {
	return s.manipulateJSON("removeProperty", func(cfg *phpjson.Map, a []phpjson.Value) {
		key := a[0].(string)
		lower := strings.ToLower(key)
		if strings.HasPrefix(key, "extra.") || strings.HasPrefix(key, "scripts.") || strings.HasPrefix(lower, "autoload.") || strings.HasPrefix(lower, "autoload-dev.") {
			bits := strings.Split(key, ".")
			last := bits[len(bits)-1]
			bits = bits[:len(bits)-1]
			propertyBase(cfg, bits[0])
			var cur phpjson.Value
			cur, _ = cfg.Get(bits[0])
			for _, bit := range bits {
				if !isset(cur, bit) {
					return
				}
				cur = get(cur, bit)
			}
			switch t := cur.(type) {
			case *phpjson.Map:
				if t.ForceObject {
					phpError("Error", "Cannot use object of type stdClass as array")
				}
				t.Delete(last)
			case string:
				phpError("Error", "Cannot unset string offsets")
			default:
				phpError("Error", "Cannot unset offset in a non-array variable")
			}
		} else {
			cfg.Delete(key)
		}
	}, name)
}

// AddLink is addLink($type, $name, $value).
func (s *Source) AddLink(linkType, name, value string) error {
	return s.manipulateJSON("addLink", func(cfg *phpjson.Map, a []phpjson.Value) {
		setPath(cfg, a[2], a[0].(string), a[1].(string))
	}, linkType, name, value)
}

// RemoveLink is removeLink($type, $name).
func (s *Source) RemoveLink(linkType, name string) error {
	if err := s.manipulateJSON("removeSubNode", func(cfg *phpjson.Map, a []phpjson.Value) {
		unsetPath(cfg, a[0].(string), a[1].(string))
	}, linkType, name); err != nil {
		return err
	}
	return s.manipulateJSON("removeMainKeyIfEmpty", func(cfg *phpjson.Map, a []phpjson.Value) {
		t := a[0].(string)
		v, _ := cfg.Get(t)
		if !isArr(v) {
			phpError("TypeError", "count(): Argument #1 ($value) must be of type Countable|array, "+typeName(v)+" given")
		}
		if arr(v).Len() == 0 {
			cfg.Delete(t)
		}
	}, linkType)
}

// callManipulator dispatches call_user_func_array([$manipulator, $method], $args).
func callManipulator(m *jsonmanip.Manipulator, method string, a []phpjson.Value) bool {
	str := func(i int) string {
		if i >= len(a) {
			return ""
		}
		s, _ := a[i].(string)
		return s
	}
	switch method {
	case "addRepository":
		return m.AddRepository(str(0), a[1], a[2].(bool))
	case "insertRepository":
		return m.InsertRepository(str(0), a[1], str(2), int(a[3].(int64)))
	case "setRepositoryUrl":
		return m.SetRepositoryUrl(str(0), str(1))
	case "removeRepository":
		return m.RemoveRepository(str(0))
	case "addConfigSetting":
		return m.AddConfigSetting(str(0), a[1])
	case "removeConfigSetting":
		return m.RemoveConfigSetting(str(0))
	case "addProperty":
		return m.AddProperty(str(0), a[1])
	case "removeProperty":
		return m.RemoveProperty(str(0))
	case "addLink":
		return m.AddLink(str(0), str(1), str(2), false)
	case "addSubNode":
		return m.AddSubNode(str(0), str(1), a[2], true)
	case "removeSubNode":
		return m.RemoveSubNode(str(0), str(1))
	case "removeMainKeyIfEmpty":
		return m.RemoveMainKeyIfEmpty(str(0))
	}
	panic("configsource: unknown manipulator method " + method)
}

// manipulateJSON is JsonConfigSource::manipulateJson().
func (s *Source) manipulateJSON(method string, fb fallback, args ...phpjson.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*jsonmanip.Exception); ok {
				err = e
				return
			}
			panic(r)
		}
	}()
	path := s.file.path

	var contents string
	if s.file.exists() {
		if !isWritable(path) {
			phpError("RuntimeException", fmt.Sprintf("The file \"%s\" is not writable.", path))
		}

		if !isReadable(path) {
			phpError("RuntimeException", fmt.Sprintf("The file \"%s\" is not readable.", path))
		}

		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		contents = string(data)
	} else if s.authConfig {
		contents = "{\n}\n"
	} else {
		contents = "{\n    \"config\": {\n    }\n}\n"
	}

	manipulator, merr := jsonmanip.New(contents)
	if merr != nil {
		return merr
	}

	newFile := !s.file.exists()

	// override manipulator method for auth config files
	if s.authConfig && (method == "addConfigSetting" || method == "removeConfigSetting") {
		mainNode, name, found := strings.Cut(args[0].(string), ".")
		target := "addSubNode"
		if method == "removeConfigSetting" {
			target = "removeSubNode"
		}
		if !found {
			// [$mainNode, $name] = explode('.', $name, 2) leaves $name null
			phpError("TypeError", "Composer\\Json\\JsonManipulator::"+target+"(): Argument #2 ($name) must be of type string, null given")
		}
		if method == "addConfigSetting" {
			args = []phpjson.Value{mainNode, name, args[1]}
		} else {
			args = []phpjson.Value{mainNode, name}
		}
		method = target
	}

	// try to update cleanly
	// PHP only warns when writing fails; the error is kept and reported
	// after the schema validation (which usually fails first, like in PHP).
	var writeErr error
	if callManipulator(manipulator, method, args) {
		writeErr = safeFilePutContents(path, manipulator.Contents())
	} else {
		// on failed clean update, call the fallback and rewrite the whole file
		decoded := s.file.read()
		config := arr(normalize(decoded))
		if config == nil {
			if decoded != nil {
				phpError("TypeError", "Composer\\Json\\JsonFile::write(): Argument #1 ($hash) must be of type array, "+typeName(decoded)+" given")
			}
			// a null document auto-vivifies into an array when written to
			config = phpjson.NewMap()
		}
		fb(config, args)
		if decoded == nil && config.Len() == 0 {
			phpError("TypeError", "Composer\\Json\\JsonFile::write(): Argument #1 ($hash) must be of type array, null given")
		}
		coerceObjects(config)
		writeErr = s.file.write(config)
	}

	if verr := s.ValidateSchema(path); verr != nil {
		if ve, ok := verr.(*ValidationError); ok {
			// restore contents to the original state
			if werr := safeFilePutContents(path, contents); werr != nil {
				return werr
			}
			return &jsonmanip.Exception{
				Class:   "RuntimeException",
				Message: "Failed to update composer.json with a valid format, reverting to the original content. Please report an issue to us with details (command you run and a copy of your composer.json). " + "\n" + strings.Join(ve.Errors, "\n"),
				Err:     ve,
			}
		}
		return verr
	}

	if newFile {
		_ = os.Chmod(path, 0o600)
	}
	return writeErr
}

// coerceObjects avoids ending up with arrays for keys that should be objects
// (applied bottom-up so parent coercion to stdClass doesn't hide deeper
// sub-keys).
func coerceObjects(config *phpjson.Map) {
	obj := func() *phpjson.Map { return &phpjson.Map{ForceObject: true} }
	if isset(config, "config", "policy") && isArr(get(config, "config", "policy")) {
		policy := arr(get(config, "config", "policy"))
		for _, listName := range policy.Keys() {
			v, _ := policy.Get(listName)
			if isEmptyArray(v) {
				policy.Set(listName, obj())
			}
		}
		if policy.Len() == 0 {
			arr(get(config, "config")).Set("policy", obj())
		}
	}
	for _, prop := range []string{"platform", "http-basic", "bearer", "gitlab-token", "gitlab-oauth", "github-oauth", "custom-headers", "forgejo-token", "preferred-install"} {
		if isset(config, "config", prop) && isEmptyArray(get(config, "config", prop)) {
			arr(get(config, "config")).Set(prop, obj())
		}
	}
	for _, prop := range []string{"psr-0", "psr-4"} {
		if isset(config, "autoload", prop) && isEmptyArray(get(config, "autoload", prop)) {
			arr(get(config, "autoload")).Set(prop, obj())
		}
		if isset(config, "autoload-dev", prop) && isEmptyArray(get(config, "autoload-dev", prop)) {
			arr(get(config, "autoload-dev")).Set(prop, obj())
		}
	}
	for _, prop := range []string{"require", "require-dev", "conflict", "provide", "replace", "suggest", "config", "autoload", "autoload-dev", "scripts", "scripts-descriptions", "scripts-aliases", "support"} {
		if isset(config, prop) && isEmptyArray(get(config, prop)) {
			config.Set(prop, obj())
		}
	}
}
