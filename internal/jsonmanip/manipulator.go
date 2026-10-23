// Package jsonmanip ports Composer\Json\JsonManipulator
// (src/Composer/Json/JsonManipulator.php of composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80), which edits composer.json text in
// place with PCRE patterns so formatting is preserved. The patterns are kept
// verbatim and run on internal/pcre, a PCRE2-compatible engine.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.
//
// Values use phpjson conventions: *phpjson.Map (ordered PHP array), []phpjson.Value
// (list), string, int64, float64, bool, nil. A *phpjson.Map with ForceObject
// set stands for stdClass / ArrayObject (an empty one formats as {}), a plain
// empty *phpjson.Map or empty list formats as []. The PHP `false` repository
// config is the Go value false.
//
// Methods return the same bool as PHP. Where PHP throws, they panic with
// *Exception; Try converts that into an error.
package jsonmanip

import (
	"sort"
	"strconv"
	"strings"

	"gocomposer/internal/pcre"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// DEFINES is JsonManipulator::DEFINES (PHP single-quoted string unescaped).
const DEFINES = `(?(DEFINE)
       (?<number>    -? (?= [1-9]|0(?!\d) ) \d++ (?:\.\d++)? (?:[eE] [+-]?+ \d++)? )
       (?<boolean>   true | false | null )
       (?<string>    " (?:[^"\\]*+ | \\ ["\\bfnrt\/] | \\ u [0-9A-Fa-f]{4} )* " )
       (?<array>     \[  (?:  (?&json) \s*+ (?: , (?&json) \s*+ )*+  )?+  \s*+ \] )
       (?<pair>      \s*+ (?&string) \s*+ : (?&json) \s*+ )
       (?<object>    \{  (?:  (?&pair)  (?: , (?&pair)  )*+  )?+  \s*+ \} )
       (?<json>      \s*+ (?: (?&number) | (?&boolean) | (?&string) | (?&array) | (?&object) ) )
    )`

// Manipulator is a JsonManipulator instance.
type Manipulator struct {
	contents string
	newline  string
	indent   string
}

// ErrNotObject is the InvalidArgumentException of the constructor.
var ErrNotObject = &Exception{Class: "InvalidArgumentException", Message: "The json file must be an object ({})"}

// New is new JsonManipulator($contents).
func New(contents string) (*Manipulator, error) {
	contents = strings.Trim(contents, " \t\n\r\x00\x0B")
	if contents == "" {
		contents = "{}"
	}
	if m, err := pcre.IsMatch(`#^\{(.*)\}$#s`, contents); err != nil {
		return nil, pcreException("preg_match", `#^\{(.*)\}$#s`, err)
	} else if m == nil {
		return nil, ErrNotObject
	}
	jm := &Manipulator{newline: "\n"}
	if strings.Contains(contents, "\r\n") {
		jm.newline = "\r\n"
	}
	if contents == "{}" {
		jm.contents = "{" + jm.newline + "}"
	} else {
		jm.contents = contents
	}
	ok, err := Try(func() bool { jm.detectIndenting(); return true })
	if !ok {
		return nil, err
	}
	return jm, nil
}

// Contents is getContents(): the document plus a trailing newline.
func (m *Manipulator) Contents() string { return m.contents + m.newline }

// Newline returns the detected line ending.
func (m *Manipulator) Newline() string { return m.newline }

// Indent returns the detected indentation.
func (m *Manipulator) Indent() string { return m.indent }

func (m *Manipulator) detectIndenting() { m.indent = DetectIndenting(m.contents) }

// DetectIndenting is JsonFile::detectIndenting.
func DetectIndenting(json string) string {
	if mt := isMatchStrictGroups(`#^([ \t]+)"#m`, json); mt != nil {
		return mt.Str(1)
	}
	return "    "
}

func (m *Manipulator) lazyPrefix() string {
	return `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?)`
}

func (m *Manipulator) nodeRegex(mainNode, kind string) string {
	return `{` + DEFINES + `^(?P<start> \s* \{ \s* (?: (?&string) \s* : (?&json) \s* , \s* )*?` +
		quote(encode(mainNode)) + `\s*:\s*)(?P<content>(?&` + kind + `))(?P<end>.*)}sx`
}

// AddLink is addLink($type, $package, $constraint, $sortPackages).
func (m *Manipulator) AddLink(linkType, pkgName, constraint string, sortPackages bool) bool {
	decoded := parseJSON(m.contents)

	// no link of that type yet
	if !isset(decoded, linkType) {
		return m.AddMainKey(linkType, mapOf(pkgName, constraint))
	}

	regex := `{` + DEFINES + m.lazyPrefix() +
		`(?P<property>` + quote(encode(linkType)) + `\s*:\s*)(?P<value>(?&json))(?P<end>.*)}sx`
	matches := isMatch(regex, m.contents)
	if matches == nil {
		return false
	}

	links := matches.N("value")

	// try to find existing link
	packageRegex := strings.ReplaceAll(quote(pkgName), "/", `\\?/`)
	regex = `{` + DEFINES + `"(?P<package>` + packageRegex + `)"(\s*:\s*)(?&string)}ix`
	if packageMatches := isMatch(regex, links); packageMatches != nil {
		// update existing link
		existingPackage := packageMatches.N("package")
		packageRegex = strings.ReplaceAll(quote(existingPackage), "/", `\\?/`)
		links = pregReplaceCallback(`{`+DEFINES+`"`+packageRegex+`"(?P<separator>\s*:\s*)(?&string)}ix`, func(mm *pcre.Match) string {
			return encode(strings.ReplaceAll(existingPackage, `\/`, "/")) + mm.N("separator") + `"` + constraint + `"`
		}, links)
	} else if match := isMatchStrictGroups(`#^\s*\{\s*\S+.*?(\s*\}\s*)$#s`, links); match != nil {
		// link missing but non empty links
		links = pregReplace(
			`{`+quote(match.Str(1))+`$}`,
			// addcslashes is used to double up backslashes/$ since preg_replace resolves them as back references otherwise, see #1588
			addcslashes(","+m.newline+m.indent+m.indent+encode(pkgName)+": "+encode(constraint)+match.Str(1)),
			links,
		)
	} else {
		// links empty
		links = "{" + m.newline +
			m.indent + m.indent + encode(pkgName) + ": " + encode(constraint) + m.newline +
			m.indent + "}"
	}

	if sortPackages {
		requirements, ok := decodeAssoc(links)
		if !ok || !isArray(requirements) {
			throw("TypeError", "Composer\\Json\\JsonManipulator::sortPackages(): Argument #1 ($packages) must be of type array, "+typeName(requirements)+" given")
		}
		links = m.Format(sortPackageKeys(requirements), 0, false)
	}

	m.contents = matches.N("start") + matches.N("property") + links + matches.N("end")

	return true
}

// sortPackageKeys sorts packages by importance (platform packages first,
// then PHP dependencies) and alphabetically, like sortPackages().
func sortPackageKeys(packages phpjson.Value) *phpjson.Map {
	out := toMap(packages)
	prefix := func(requirement string) string {
		if isPlatformPackage(requirement) {
			s, err := pcre.ReplaceAll(
				[]string{`/^php/`, `/^hhvm/`, `/^ext/`, `/^lib/`, `/^\D/`},
				[]string{`0-$0`, `1-$0`, `2-$0`, `3-$0`, `4-$0`},
				requirement,
			)
			if err != nil {
				panic(pcreException("preg_replace", "/^php/, /^hhvm/, /^ext/, /^lib/, /^\\D/", err))
			}
			return s
		}
		return "5-" + requirement
	}
	keys := out.Keys()
	sort.SliceStable(keys, func(i, j int) bool {
		return semver.NatCompare(prefix(keys[i]), prefix(keys[j])) < 0
	})
	sorted := phpjson.NewMap()
	for _, k := range keys {
		v, _ := out.Get(k)
		sorted.Set(k, v)
	}
	return sorted
}

// PlatformPackageRegex is PlatformRepository::PLATFORM_PACKAGE_REGEX.
const PlatformPackageRegex = `{^(?:php(?:-64bit|-ipv6|-zts|-debug)?|hhvm|(?:ext|lib)-[a-z0-9](?:[_.-]?[a-z0-9]+)*|composer(?:-(?:plugin|runtime)-api)?)$}iD`

func isPlatformPackage(name string) bool { return isMatch(PlatformPackageRegex, name) != nil }

// AddRepository is addRepository($name, $config, $append). config is an
// array value (*phpjson.Map or list) or false.
func (m *Manipulator) AddRepository(name string, config phpjson.Value, appendItem bool) bool {
	if name != "" && !m.doRemoveRepository(name) {
		return false
	}

	if !m.doConvertRepositoriesFromAssocToList() {
		return false
	}

	if isArray(config) && !isNumeric(name) && name != "" {
		config = prependName(name, config)
	} else if b, ok := config.(bool); ok && !b {
		config = mapOf(name, config)
	}

	return m.AddListItem("repositories", config, appendItem)
}

func (m *Manipulator) doConvertRepositoriesFromAssocToList() bool {
	decoded, _ := decodeObject(m.contents)
	repos := objectProperty(decoded, "repositories")

	if repositories, ok := repos.(*phpjson.Map); ok && isObject(repositories) {
		// delete from bottom to top, to ensure keys stay the same
		keys := repositories.Keys()
		for i := len(keys) - 1; i >= 0; i-- {
			if !m.RemoveSubNode("repositories", keys[i]) {
				return false
			}
		}

		m.ChangeEmptyMainKeyFromAssocToList("repositories")

		// re-add in order
		for _, repositoryName := range keys {
			repository, _ := repositories.Get(repositoryName)
			if !isObject(repository) {
				if !m.AddListItem("repositories", mapOf(repositoryName, repository), true) {
					return false
				}
			} else if isIntKey(repositoryName) || isNumeric(repositoryName) {
				if !m.AddListItem("repositories", repository, true) {
					return false
				}
			} else {
				// prepend name property
				if !m.AddListItem("repositories", prependName(repositoryName, repository), true) {
					return false
				}
			}
		}
	}

	return true
}

// objectProperty is "$decoded->$name ?? null" on a json_decode() object.
func objectProperty(decoded phpjson.Value, name string) phpjson.Value {
	if !isObject(decoded) {
		return nil
	}
	v, _ := decoded.(*phpjson.Map).Get(name)
	return v
}

// SetRepositoryUrl is setRepositoryUrl($name, $url).
func (m *Manipulator) SetRepositoryUrl(name, url string) bool {
	decoded := parseJSON(m.contents)
	repositoryIndex := ""
	found := false

	repos := get(decoded, "repositories")
	phpjson.Each(arrayOrEmpty(repos), func(index string, repository phpjson.Value) {
		if found {
			return
		}
		if !isIntKey(index) && name == index {
			repositoryIndex, found = index, true
			return
		}
		if s, ok := get(repository, "name").(string); ok && isArray(repository) && s == name {
			repositoryIndex, found = index, true
		}
	})

	if !found {
		return false
	}

	listRegex := ""
	if i, ok := phpjson.IntKey(repositoryIndex); ok {
		if i < 0 {
			i = 0
		}
		listRegex = `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?"repositories"\s*:\s*\[\s*((?&json)\s*+,\s*+){` + strconv.FormatInt(i, 10) + `})(?P<repository>(?&object))(?P<end>.*)}sx`
	}

	objectRegex := `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?"repositories"\s*:\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?` + quote(encodeKey(repositoryIndex)) + `\s*:\s*)(?P<repository>(?&object))(?P<end>.*)}sx`

	var matches *pcre.Match
	if listRegex != "" {
		matches = isMatch(listRegex, m.contents)
	}
	if matches == nil {
		matches = isMatch(objectRegex, m.contents)
	}
	if matches != nil {
		// invalid match due to un-regexable content, abort
		if v, ok := decodeObject(matches.N("repository")); ok && v == false {
			return false
		}

		repositoryRegex := `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?"url"\s*:\s*)(?P<url>(?&string))(?P<end>.*)}sx`

		m.contents = matches.N("start") + pregReplaceCallback(repositoryRegex, func(rm *pcre.Match) string {
			return rm.N("start") + encode(url) + rm.N("end")
		}, matches.N("repository")) + matches.N("end")

		return true
	}

	return false
}

func arrayOrEmpty(v phpjson.Value) phpjson.Value {
	if isArray(v) {
		return v
	}
	return []phpjson.Value{}
}

// InsertRepository is insertRepository($name, $config, $referenceName, $offset).
func (m *Manipulator) InsertRepository(name string, config phpjson.Value, referenceName string, offset int) bool {
	if name != "" && !m.doRemoveRepository(name) {
		return false
	}

	if !m.doConvertRepositoriesFromAssocToList() {
		return false
	}

	indexToInsert := ""
	found := false
	decoded := parseJSON(m.contents)

	phpjson.Each(arrayOrEmpty(get(decoded, "repositories")), func(repositoryIndex string, repository phpjson.Value) {
		if found {
			return
		}
		if s, ok := get(repository, "name").(string); ok && isArray(repository) && s == referenceName {
			indexToInsert, found = repositoryIndex, true
			return
		}
		if !isIntKey(repositoryIndex) && repositoryIndex == referenceName {
			indexToInsert, found = repositoryIndex, true
			return
		}
		if strictEqualFalseMap(referenceName, repository) {
			indexToInsert, found = repositoryIndex, true
		}
	})

	if !found {
		return false
	}

	if isArray(config) && !isNumeric(name) && name != "" {
		config = prependName(name, config)
	} else if b, ok := config.(bool); ok && !b {
		config = mapOf("name", config)
	}

	idx, ok := phpjson.IntKey(indexToInsert)
	if !ok {
		throw("TypeError", "Unsupported operand types: string + int")
	}
	return m.InsertListItem("repositories", config, int(idx)+offset)
}

// RemoveRepository is removeRepository($name).
func (m *Manipulator) RemoveRepository(name string) bool {
	return m.doRemoveRepository(name) && m.RemoveMainKeyIfEmpty("repositories")
}

func (m *Manipulator) doRemoveRepository(name string) bool {
	decoded, _ := decodeObject(m.contents)
	repos := objectProperty(decoded, "repositories")
	isAssoc := isObject(repos)

	var reposArr phpjson.Value = []phpjson.Value{}
	if repos != nil {
		reposArr = castArray(repos)
	}

	type entry struct {
		k string
		v phpjson.Value
	}
	var entries []entry
	phpjson.Each(reposArr, func(k string, v phpjson.Value) { entries = append(entries, entry{k, v}) })

	for _, e := range entries {
		repositoryIndex, repository := e.k, e.v
		if !isIntKey(repositoryIndex) && repositoryIndex == name && isAssoc {
			if !m.RemoveSubNode("repositories", repositoryIndex) {
				return false
			}

			break
		}

		if s, ok := objectProperty(repository, "name").(string); ok && s == name {
			if isAssoc {
				if !m.RemoveSubNode("repositories", repositoryIndex) {
					return false
				}
			} else {
				if !m.RemoveListItem("repositories", intCast(repositoryIndex)) {
					return false
				}
			}

			break
		}

		if isAssoc {
			if b, ok := repository.(bool); ok && !b && !isIntKey(repositoryIndex) && name == repositoryIndex {
				if !m.RemoveSubNode("repositories", repositoryIndex) {
					return false
				}

				return true
			}
		} else {
			repositoryAsArray := castArray(repository)

			if b, ok := get(repositoryAsArray, name).(bool); ok && !b && count(repositoryAsArray) == 1 {
				if !m.RemoveListItem("repositories", intCast(repositoryIndex)) {
					return false
				}

				return true
			}
		}
	}

	return true
}

// castArray is PHP's (array) cast.
func castArray(v phpjson.Value) phpjson.Value {
	switch t := v.(type) {
	case nil:
		return []phpjson.Value{}
	case *phpjson.Map:
		c := phpjson.NewMap()
		for _, k := range t.Keys() {
			x, _ := t.Get(k)
			c.Set(k, x)
		}
		return c
	case []phpjson.Value:
		return t
	}
	return []phpjson.Value{v}
}

// intCast is (int) on an array key.
func intCast(k string) int {
	if i, ok := phpjson.IntKey(k); ok {
		return int(i)
	}
	n, _ := strconv.Atoi(k)
	return n
}

// AddConfigSetting is addConfigSetting($name, $value).
func (m *Manipulator) AddConfigSetting(name string, value phpjson.Value) bool {
	// policy config has one more nesting level than the rest of the config e.g. policy.list-name.key
	// Rewrite as a policy.<list> update with the merged list block so addSubNode can do a
	// surgical edit at the list level rather than falling back to a whole-file rewrite.
	if strings.HasPrefix(name, "policy.") && strings.Count(name, ".") >= 2 {
		listName, fieldName, _ := strings.Cut(name[7:], ".")
		if strings.Contains(fieldName, ".") {
			return false
		}
		listBlock := m.getCurrentPolicyListBlock(listName)
		listBlock.Set(fieldName, value)

		return m.AddSubNode("config", "policy."+listName, listBlock, true)
	}

	return m.AddSubNode("config", name, value, true)
}

// RemoveConfigSetting is removeConfigSetting($name).
func (m *Manipulator) RemoveConfigSetting(name string) bool {
	// policy.<list>.<field>: drop the field from the list block, then either rewrite the
	// block (surgical) or defer to the whole-file fallback when the block ends up empty so
	// the cascade-cleanup of empty ancestors can run.
	if strings.HasPrefix(name, "policy.") && strings.Count(name, ".") >= 2 {
		listName, fieldName, _ := strings.Cut(name[7:], ".")
		if strings.Contains(fieldName, ".") {
			return false
		}

		listBlock := m.getCurrentPolicyListBlock(listName)
		if !listBlock.Has(fieldName) {
			return true
		}
		listBlock.Delete(fieldName)
		if listBlock.Len() == 0 {
			return m.RemoveConfigSetting("policy." + listName)
		}

		return m.AddSubNode("config", "policy."+listName, listBlock, true)
	}

	return m.RemoveSubNode("config", name)
}

func (m *Manipulator) getCurrentPolicyListBlock(listName string) *phpjson.Map {
	decoded := parseJSON(m.contents)
	if !isArray(decoded) {
		return phpjson.NewMap()
	}
	policySection := get(decoded, "config", "policy")
	if !isArray(policySection) {
		return phpjson.NewMap()
	}
	listBlock := get(policySection, listName)
	if !isArray(listBlock) {
		return phpjson.NewMap()
	}
	return toMap(listBlock)
}

// AddProperty is addProperty($name, $value).
func (m *Manipulator) AddProperty(name string, value phpjson.Value) bool {
	if strings.HasPrefix(name, "suggest.") {
		return m.AddSubNode("suggest", name[8:], value, true)
	}

	if strings.HasPrefix(name, "extra.") {
		return m.AddSubNode("extra", name[6:], value, true)
	}

	if strings.HasPrefix(name, "scripts.") {
		return m.AddSubNode("scripts", name[8:], value, true)
	}

	return m.AddMainKey(name, value)
}

// RemoveProperty is removeProperty($name).
func (m *Manipulator) RemoveProperty(name string) bool {
	if strings.HasPrefix(name, "suggest.") {
		return m.RemoveSubNode("suggest", name[8:])
	}

	if strings.HasPrefix(name, "extra.") {
		return m.RemoveSubNode("extra", name[6:])
	}

	if strings.HasPrefix(name, "scripts.") {
		return m.RemoveSubNode("scripts", name[8:])
	}

	if strings.HasPrefix(name, "autoload.") {
		return m.RemoveSubNode("autoload", name[9:])
	}

	if strings.HasPrefix(name, "autoload-dev.") {
		return m.RemoveSubNode("autoload-dev", name[13:])
	}

	return m.RemoveMainKey(name)
}

func splitSubName(mainNode, name string) (string, *string) {
	if (mainNode == "config" || mainNode == "extra" || mainNode == "scripts") && strings.Contains(name, ".") {
		n, sub, _ := strings.Cut(name, ".")
		return n, &sub
	}
	return name, nil
}

// AddSubNode is addSubNode($mainNode, $name, $value, $append).
func (m *Manipulator) AddSubNode(mainNode, name string, value phpjson.Value, appendItem bool) bool {
	decoded := parseJSON(m.contents)

	name, subName := splitSubName(mainNode, name)

	// no main node yet
	if !isset(decoded, mainNode) {
		if subName != nil {
			m.AddMainKey(mainNode, mapOf(name, mapOf(*subName, value)))
		} else {
			m.AddMainKey(mainNode, mapOf(name, value))
		}

		return true
	}

	// main node content not match-able
	nodeRegex := m.nodeRegex(mainNode, "object")

	match, ok := isMatchCatch(nodeRegex, m.contents, pregBacktrackLimitError)
	if !ok || match == nil {
		return false
	}

	children := match.N("content")
	// invalid match due to un-regexable content, abort
	if v, ok := decodeObject(children); !ok || !truthy(v) {
		return false
	}

	// child exists
	if mainVal := get(decoded, mainNode); isArray(mainVal) && arrayKeyExists(mainVal, name) {
		matches := m.matchChild(name, children)
		if matches == nil {
			return false
		}

		if subName != nil {
			curVal, ok := decodeAssoc(matches.N("content"))
			cur := phpjson.NewMap()
			if ok && isArray(curVal) {
				cur = toMap(curVal)
			}
			cur.Set(*subName, value)
			value = cur
		}

		children = matches.N("start") + m.Format(value, 1, false) + matches.N("end") + children[len(matches.Str(0)):]
	} else if match := isMatch(`#^\{(?P<leadingspace>\s*?)(?P<content>\S+.*?)?(?P<trailingspace>\s*)\}$#s`, children); match != nil {
		whitespace := match.N("trailingspace")
		if _, hasContent := match.Named("content"); hasContent {
			if subName != nil {
				value = mapOf(*subName, value)
			}

			// child missing but non empty children
			if appendItem {
				children = pregReplace(
					`#`+whitespace+`}$#`,
					addcslashes(","+m.newline+m.indent+m.indent+encode(name)+": "+m.Format(value, 1, false)+whitespace+"}"),
					children,
				)
			} else {
				whitespace = match.N("leadingspace")
				children = pregReplace(
					`#^{`+whitespace+`#`,
					addcslashes("{"+whitespace+encode(name)+": "+m.Format(value, 1, false)+","+m.newline+m.indent+m.indent),
					children,
				)
			}
		} else {
			if subName != nil {
				value = mapOf(*subName, value)
			}

			// children present but empty
			children = "{" + m.newline + m.indent + m.indent + encode(name) + ": " + m.Format(value, 1, false) + whitespace + "}"
		}
	} else {
		throw("LogicException", "Nothing matched above for: "+children)
	}

	m.contents = pregReplaceCallback(nodeRegex, func(mm *pcre.Match) string {
		return mm.N("start") + children + mm.N("end")
	}, m.contents)

	return true
}

// RemoveSubNode is removeSubNode($mainNode, $name).
func (m *Manipulator) RemoveSubNode(mainNode, name string) bool {
	decoded := parseJSON(m.contents)

	// no node or empty node
	if !truthy(get(decoded, mainNode)) {
		return true
	}

	// no node content match-able
	nodeRegex := m.nodeRegex(mainNode, "object")
	match, ok := isMatchCatch(nodeRegex, m.contents, pregBacktrackLimitError)
	if !ok || match == nil {
		return false
	}

	children := match.N("content")

	// invalid match due to un-regexable content, abort
	if v, ok := decodeAssoc(children); !ok || !truthy(v) {
		return false
	}

	name, subName := splitSubName(mainNode, name)

	// no node to remove
	if !isset(decoded, mainNode, name) || (subName != nil && phpjson.Truthy(*subName) && !isset(decoded, mainNode, name, *subName)) {
		return true
	}

	matches := m.matchChild(name, children)
	if matches == nil {
		return false
	}

	// drop the pair along with the comma before it, or after it if it is the first pair
	rest := children[len(matches.Str(0))-len(matches.N("end")):]
	var childrenClean string
	if isMatch(`{,\s*$}`, matches.N("prefix")) != nil {
		childrenClean = pregReplace(`{,\s*$}`, "", matches.N("prefix")) + rest
	} else {
		childrenClean = matches.N("prefix") + pregReplace(`{^\s*,?\s*}`, "", rest)
	}

	// no child data left, $name was the only key in
	if match := isMatch(`#^\{\s*?(?P<content>\S+.*?)?(?P<trailingspace>\s*)\}$#s`, childrenClean); match != nil {
		if _, hasContent := match.Named("content"); !hasContent {
			newline := m.newline
			indent := m.indent

			m.contents = pregReplaceCallback(nodeRegex, func(mm *pcre.Match) string {
				return mm.N("start") + "{" + newline + indent + "}" + mm.N("end")
			}, m.contents)

			// we have a subname, so we restore the rest of $name
			if subName != nil {
				curVal, _ := decodeAssoc(children)
				m.AddSubNode(mainNode, name, unsetSub(curVal, name, *subName), true)
			}

			return true
		}
	}

	m.contents = pregReplaceCallback(nodeRegex, func(mm *pcre.Match) string {
		clean := childrenClean
		if subName != nil {
			curVal, _ := decodeAssoc(mm.N("content"))
			cur := toMap(curVal)
			cur.Set(name, unsetSub(curVal, name, *subName))
			clean = m.Format(cur, 0, true)
		}

		return mm.N("start") + clean + mm.N("end")
	}, m.contents)

	return true
}

// unsetSub performs unset($curVal[$name][$subName]) and returns the new
// $curVal[$name], replacing an empty array by new \ArrayObject().
func unsetSub(curVal phpjson.Value, name, subName string) phpjson.Value {
	v, _ := index(curVal, name)
	switch t := v.(type) {
	case *phpjson.Map, []phpjson.Value:
		c := toMap(t)
		c.Delete(subName)
		if c.Len() == 0 {
			return &phpjson.Map{ForceObject: true}
		}
		return c
	case string:
		throw("Error", "Cannot unset string offsets")
	case nil:
		// unset on null is a no-op; null === [] is false
		return nil
	default:
		throw("Error", "Cannot unset offset in a non-array variable")
	}
	return nil
}

func (m *Manipulator) arrayNodeMatch(mainNode string) (*pcre.Match, string, bool) {
	nodeRegex := m.nodeRegex(mainNode, "array")
	match, ok := isMatchCatch(nodeRegex, m.contents, pregBacktrackLimitError)
	if !ok || match == nil {
		return nil, nodeRegex, false
	}
	return match, nodeRegex, true
}

// AddListItem is addListItem($mainNode, $value, $append).
func (m *Manipulator) AddListItem(mainNode string, value phpjson.Value, appendItem bool) bool {
	decoded := parseJSON(m.contents)

	// no main node yet
	if !isset(decoded, mainNode) {
		if !m.AddMainKey(mainNode, []phpjson.Value{}) {
			return false
		}
	}

	// main node content not match-able
	match, nodeRegex, ok := m.arrayNodeMatch(mainNode)
	if !ok {
		return false
	}

	children := match.N("content")
	// invalid match due to un-regexable content, abort
	if v, ok := decodeObject(children); ok && v == false {
		return false
	}

	if match := isMatch(`#^\[(?P<leadingspace>\s*?)(?P<content>\S+.*?)?(?P<trailingspace>\s*)\]$#s`, children); match != nil {
		leadingWhitespace := match.N("leadingspace")
		whitespace := match.N("trailingspace")
		leadingItemWhitespace := m.newline + m.indent + m.indent
		trailingItemWhitespace := whitespace
		itemDepth := 1

		// keep oneline lists as one line
		if !strings.Contains(whitespace, m.newline) {
			leadingItemWhitespace = leadingWhitespace
			trailingItemWhitespace = leadingWhitespace
			itemDepth = 0
		}

		if _, hasContent := match.Named("content"); hasContent {
			// child missing but non empty children
			if appendItem {
				children = pregReplace(
					`#`+whitespace+`]$#`,
					addcslashes(","+leadingItemWhitespace+m.Format(value, itemDepth, false)+trailingItemWhitespace+"]"),
					children,
				)
			} else {
				whitespace = match.N("leadingspace")
				children = pregReplace(
					`#^\[`+whitespace+`#`,
					addcslashes("["+whitespace+m.Format(value, itemDepth, false)+","+leadingItemWhitespace),
					children,
				)
			}
		} else {
			// children present but empty
			children = "[" + leadingItemWhitespace + m.Format(value, itemDepth, false) + trailingItemWhitespace + "]"
		}
	} else {
		throw("LogicException", "Nothing matched above for: "+children)
	}

	m.contents = pregReplaceCallback(nodeRegex, func(mm *pcre.Match) string {
		return mm.N("start") + children + mm.N("end")
	}, m.contents)

	return true
}

// InsertListItem is insertListItem($mainNode, $value, $index).
func (m *Manipulator) InsertListItem(mainNode string, value phpjson.Value, idx int) bool {
	if idx < 0 {
		throw("InvalidArgumentException", "Index can only be positive integer")
	}

	if idx == 0 {
		return m.AddListItem(mainNode, value, false)
	}

	decoded := parseJSON(m.contents)

	// no main node yet
	if !isset(decoded, mainNode) {
		if !m.AddMainKey(mainNode, []phpjson.Value{}) {
			return false
		}
	}

	node := get(decoded, mainNode)
	if !isArray(node) {
		throw("TypeError", "count(): Argument #1 ($value) must be of type Countable|array, "+typeName(node)+" given")
	}
	if count(node) == idx {
		return m.AddListItem(mainNode, value, true)
	}

	// main node content not match-able
	match, nodeRegex, ok := m.arrayNodeMatch(mainNode)
	if !ok {
		return false
	}

	children := match.N("content")
	// invalid match due to un-regexable content, abort
	if v, ok := decodeObject(children); ok && v == false {
		return false
	}

	listSkipToItemRegex := `{` + DEFINES + `^(?P<start>\[\s*((?&json)\s*+,\s*?){` + strconv.Itoa(max(0, idx)) + `})(?P<space_before_item>(\s*))(?P<end>.*)}sx`

	children = pregReplaceCallback(listSkipToItemRegex, func(mm *pcre.Match) string {
		return mm.N("start") + mm.N("space_before_item") + m.Format(value, 1, false) + "," + mm.N("space_before_item") + mm.N("end")
	}, children)

	m.contents = pregReplaceCallback(nodeRegex, func(mm *pcre.Match) string {
		return mm.N("start") + children + mm.N("end")
	}, m.contents)

	return true
}

// RemoveListItem is removeListItem($mainNode, $nodeIndex).
func (m *Manipulator) RemoveListItem(mainNode string, nodeIndex int) bool {
	// invalid index, that cannot be removed anyway
	if nodeIndex < 0 {
		return true
	}

	decoded := parseJSON(m.contents)

	// no node or empty node
	if node, ok := index(decoded, mainNode); ok && isArray(node) && count(node) == 0 {
		return true
	}

	// no node content match-able
	match, _, ok := m.arrayNodeMatch(mainNode)
	if !ok {
		return false
	}

	children := match.N("content")

	// invalid match due to un-regexable content, abort
	if v, ok := decodeAssoc(children); ok && v == false {
		return false
	}

	// no node to remove
	if !isset(decoded, mainNode, strconv.Itoa(nodeIndex)) {
		return true
	}

	contentRegex := `(?&json)`
	var startRegex, endRegex string

	if nodeIndex > 1 {
		startRegex = `(?&json)\s*+(?:,(?&json)\s*+){` + strconv.Itoa(nodeIndex-1) + `}`
		// remove leading array separator in case we might remove the last
		contentRegex = `\s*+,?\s*+` + contentRegex
		endRegex = `(?:(\s*+,\s*+(?&json))*(?:\s*+(?&json))?)\s*+`
	} else if nodeIndex > 0 {
		startRegex = `(?&json)\s*+`
		// remove leading array separator in case we might remove the last
		contentRegex = `\s*+,?\s*+` + contentRegex
		endRegex = `(?:(\s*+,\s*+(?&json))*(?:\s*+(?&json))?)\s*+`
	} else {
		startRegex = `\s*+`
		// remove trailing array separator when we delete first
		contentRegex = contentRegex + `\s*+,?\s*+`
		endRegex = `(?:((?&json)\s*+,\s*+)*(?:\s*+(?&json))?)\s*+`
	}

	if childMatch := isMatch(`{`+DEFINES+`(?P<start>\[`+startRegex+`)(?P<content>`+contentRegex+`)(?P<end>`+endRegex+`\])}sx`, children); childMatch != nil {
		m.contents = match.N("start") + childMatch.N("start") + childMatch.N("end") + match.N("end")

		return true
	}

	return false
}

// AddMainKey is addMainKey($key, $content).
func (m *Manipulator) AddMainKey(key string, content phpjson.Value) bool {
	decoded := parseJSON(m.contents)
	formatted := m.Format(content, 0, false)

	// key exists already
	regex := `{` + DEFINES + m.lazyPrefix() +
		`(?P<key>` + quote(encode(key)) + `\s*:\s*(?&json))(?P<end>.*)}sx`
	if isset(decoded, key) {
		if matches := isMatch(regex, m.contents); matches != nil {
			// invalid match due to un-regexable content, abort
			if v, ok := decodeObject("{" + matches.N("key") + "}"); !ok || !truthy(v) {
				return false
			}

			m.contents = matches.N("start") + encode(key) + ": " + formatted + matches.N("end")

			return true
		}
	}

	// append at the end of the file and keep whitespace
	if match := isMatch(`#[^{\s](\s*)\}$#`, m.contents); match != nil {
		m.contents = pregReplace(
			`#`+match.Str(1)+`\}$#`,
			addcslashes(","+m.newline+m.indent+encode(key)+": "+formatted+m.newline+"}"),
			m.contents,
		)

		return true
	}

	// append at the end of the file
	m.contents = pregReplace(
		`#\}$#`,
		addcslashes(m.indent+encode(key)+": "+formatted+m.newline+"}"),
		m.contents,
	)

	return true
}

func requireArray(fn string, v phpjson.Value) {
	if !isArray(v) {
		throw("TypeError", fn+"(): Argument #2 ($array) must be of type array, "+typeName(v)+" given")
	}
}

// RemoveMainKey is removeMainKey($key).
func (m *Manipulator) RemoveMainKey(key string) bool {
	decoded := parseJSON(m.contents)

	requireArray("array_key_exists", decoded)
	if !arrayKeyExists(decoded, key) {
		return true
	}

	// key exists already
	regex := `{` + DEFINES + m.lazyPrefix() +
		`(?P<removal>` + quote(encode(key)) + `\s*:\s*(?&json))\s*,?\s*(?P<end>.*)}sx`
	if matches := isMatch(regex, m.contents); matches != nil {
		// invalid match due to un-regexable content, abort
		if v, ok := decodeObject("{" + matches.N("removal") + "}"); !ok || !truthy(v) {
			return false
		}

		start := matches.N("start")
		// check that we are not leaving a dangling comma on the previous line if the last line was removed
		if isMatchStrictGroups(`#,\s*$#`, start) != nil && isMatch(`#^\}$#`, matches.N("end")) != nil {
			start = rtrimChars(pregReplace(`#,(\s*)$#`, "$1", start), m.indent)
		}

		m.contents = start + matches.N("end")
		if isMatch(`#^\{\s*\}\s*$#`, m.contents) != nil {
			m.contents = "{\n}"
		}

		return true
	}

	return false
}

// ChangeEmptyMainKeyFromAssocToList is changeEmptyMainKeyFromAssocToList($key).
func (m *Manipulator) ChangeEmptyMainKeyFromAssocToList(key string) bool {
	decoded := parseJSON(m.contents)

	requireArray("array_key_exists", decoded)
	if !arrayKeyExists(decoded, key) {
		return true
	}

	regex := `{` + DEFINES + `^(?P<start>\s*\{\s*(?:(?&string)\s*:\s*(?&json)\s*,\s*)*?` + quote(encode(key)) + `\s*:\s*)(?P<removal>\{(?P<removal_space>\s*+)\})(?P<end>\s*,?\s*.*)}sx`
	if matches := isMatch(regex, m.contents); matches != nil {
		// invalid match due to un-regexable content, abort
		if v, ok := decodeObject(matches.N("removal")); ok && v == false {
			return false
		}

		m.contents = matches.N("start") + "[" + matches.N("removal_space") + "]" + matches.N("end")

		return true
	}

	return false
}

// RemoveMainKeyIfEmpty is removeMainKeyIfEmpty($key).
func (m *Manipulator) RemoveMainKeyIfEmpty(key string) bool {
	decoded := parseJSON(m.contents)

	requireArray("array_key_exists", decoded)
	if !arrayKeyExists(decoded, key) {
		return true
	}

	if v := get(decoded, key); isArray(v) && count(v) == 0 {
		return m.RemoveMainKey(key)
	}

	return true
}

// matchChild matches the top-level pair for name in an object, skipping
// nested keys with the same name; nil if it could not be matched.
func (m *Manipulator) matchChild(name, children string) *pcre.Match {
	// spaces are escaped for the x modifier, slashes may have been written escaped
	key := strings.NewReplacer(" ", `\ `, "/", `\\?/`).Replace(quote(encode(name)))

	// possessive loop over the other pairs, a lazy one backtracks exponentially when the key is missing
	childRegex := `{` + DEFINES + `^(?P<start>(?P<prefix>\{ \s* (?: (?!` + key + `\s*:) (?&string) \s* : (?&json) \s* , \s* )*+) ` + key + `\s*:\s*)(?P<content>(?&json))(?P<end>,?)}x`

	matches, ok := isMatchCatch(childRegex, children, pregBacktrackLimitError, pregJitStackLimitError)
	if !ok {
		return nil
	}
	return matches
}

// Format is format($data, $depth, $wasObject).
func (m *Manipulator) Format(data phpjson.Value, depth int, wasObject bool) string {
	if t, ok := data.(*phpjson.Map); ok && t != nil && t.ForceObject {
		wasObject = true
	}

	if isArray(data) || isObject(data) {
		if count(data) == 0 {
			if wasObject {
				return "{" + m.newline + strings.Repeat(m.indent, depth+1) + "}"
			}
			return "[]"
		}

		if l, isList := data.([]phpjson.Value); isList || data.(*phpjson.Map).IsList() {
			var items []string
			if isList {
				for _, v := range l {
					items = append(items, m.Format(v, depth+1, false))
				}
			} else {
				phpjson.Each(data, func(_ string, v phpjson.Value) {
					items = append(items, m.Format(v, depth+1, false))
				})
			}
			return "[" + strings.Join(items, ", ") + "]"
		}

		out := "{" + m.newline
		var elems []string
		phpjson.Each(data, func(key string, val phpjson.Value) {
			elems = append(elems, strings.Repeat(m.indent, depth+2)+encode(key)+": "+m.Format(val, depth+1, false))
		})

		return out + strings.Join(elems, ","+m.newline) + m.newline + strings.Repeat(m.indent, depth+1) + "}"
	}

	return encode(data)
}
