// Ported for gocomposer (MIT) from justinrainbow/json-schema 6.13.1
// (MIT, Copyright (c) 2016 justinrainbow/json-schema contributors):
// src/JsonSchema/SchemaStorage.php, src/JsonSchema/Entity/JsonPointer.php,
// src/JsonSchema/Uri/UriResolver.php, src/JsonSchema/Uri/UriRetriever.php and
// src/JsonSchema/Uri/Retrievers/FileGetContents.php.
//
// Remote retrieval (http/https) and the bundled meta-schemas
// (package://dist/schema/...) are not ported: such references fail with an
// explicit error instead of being fetched.
package jsonschema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gocomposer/internal/phpjson"
)

const internalProvidedSchemaURI = "internal://provided-schema/"

const (
	draft3 = "http://json-schema.org/draft-03/schema#"
	draft4 = "http://json-schema.org/draft-04/schema#"
)

var schemaMapKeywords = []string{"properties", "patternProperties", "definitions", "$defs", "dependencies", "dependentSchemas"}

func inStrings(s string, list []string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// ---- JsonPointer ----

type jsonPointer struct {
	filename string
	paths    []string
}

var (
	pointerDecoder = strings.NewReplacer("%25", "%", "~1", "/", "~0", "~")
	pointerEncoder = strings.NewReplacer("/", "~1", "~", "~0", "%", "%25")
)

func newJSONPointer(value string) jsonPointer {
	p := jsonPointer{}
	file, frag, has := strings.Cut(value, "#")
	p.filename = file
	if has {
		for _, seg := range strings.Split(strings.Trim(frag, "/"), "/") {
			seg = pointerDecoder.Replace(seg)
			if seg != "" {
				p.paths = append(p.paths, seg)
			}
		}
	}
	return p
}

func (p jsonPointer) propertyPathAsString() string {
	enc := make([]string, len(p.paths))
	for i, s := range p.paths {
		enc[i] = pointerEncoder.Replace(s)
	}
	return strings.TrimRight("#/"+strings.Join(enc, "/"), "/")
}

func (p jsonPointer) String() string { return p.filename + p.propertyPathAsString() }

func (p jsonPointer) withAppended(seg string) jsonPointer {
	paths := make([]string, len(p.paths), len(p.paths)+1)
	copy(paths, p.paths)
	return jsonPointer{filename: p.filename, paths: append(paths, seg)}
}

// ---- UriResolver ----

var uriParseRe = mustRegex(`^(([^:/?#]+):)?(//([^/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`, "")

type uriComponents struct {
	scheme, authority, path string
	query, fragment         *string
	ok                      bool
}

func parseURI(uri string) uriComponents {
	res, _ := uriParseRe.exec(uri)
	c := uriComponents{}
	if res == nil {
		return c
	}
	n := res.count()
	if n > 5 {
		c.ok = true
		c.scheme, c.authority, c.path = res.group(2), res.group(4), res.group(5)
	}
	if n > 7 {
		q := res.group(7)
		c.query = &q
	}
	if n > 9 {
		f := res.group(9)
		c.fragment = &f
	}
	return c
}

func generateURI(c uriComponents) string {
	uri := c.scheme + "://" + c.authority + c.path
	if c.query != nil && *c.query != "" {
		uri += "?" + *c.query
	}
	if c.fragment != nil {
		uri += "#" + *c.fragment
	}
	return uri
}

var schemeSlashRe = mustRegex(`^[^/]+://`, "u")

func resolveURI(uri string, baseURI *string) string {
	if baseURI != nil && pregMatch(schemeSlashRe, *baseURI) != 1 && !looksLikeURL(*baseURI) {
		b := *baseURI
		if st, err := os.Stat(b); err == nil && !st.IsDir() {
			abs, _ := filepath.Abs(b)
			b = "file://" + abs
		} else if err == nil && st.IsDir() {
			abs, _ := filepath.Abs(b)
			b = "file://" + abs + "/"
		} else {
			wd, _ := os.Getwd()
			b = "file://" + wd + "/" + b
		}
		baseURI = &b
	}
	if uri == "" {
		if baseURI == nil {
			return ""
		}
		return *baseURI
	}
	comps := parseURI(uri)
	if comps.scheme != "" && comps.scheme != "0" {
		return uri
	}
	base := ""
	if baseURI != nil {
		base = *baseURI
	}
	bc := parseURI(base)
	basePath := bc.path
	if comps.path != "" && basePath == "" && bc.authority != "" {
		basePath = "/"
	}
	bc.path = combineRelativePathWithBasePath(comps.path, basePath)
	if comps.fragment != nil {
		bc.fragment = comps.fragment
	}
	return generateURI(bc)
}

// looksLikeURL approximates filter_var($base, FILTER_VALIDATE_URL) for bases
// such as "urn:x" or "mailto:a@b" that have a scheme but no "//".
func looksLikeURL(s string) bool {
	c := parseURI(s)
	return c.scheme != "" && (c.authority != "" || c.scheme == "mailto" || c.scheme == "news" || c.scheme == "file")
}

func normalizeURIPath(path string) string {
	var sb strings.Builder
	for i := 0; i < len(path); {
		if strings.HasPrefix(path[i:], "./") && (i == 0 || path[i-1] != '.') {
			i += 2
			continue
		}
		sb.WriteByte(path[i])
		i++
	}
	return strings.ReplaceAll(sb.String(), "//", "/")
}

func phpDirname(p string) string {
	// PHP dirname() for '/'-separated paths
	if p == "" {
		return ""
	}
	t := strings.TrimRight(p, "/")
	if t == "" {
		return "/"
	}
	i := strings.LastIndexByte(t, '/')
	if i < 0 {
		return "."
	}
	d := strings.TrimRight(t[:i], "/")
	if d == "" {
		return "/"
	}
	return d
}

func combineRelativePathWithBasePath(relativePath, basePath string) string {
	relativePath = normalizeURIPath(relativePath)
	if relativePath == "" || relativePath == "0" {
		return basePath
	}
	if relativePath[0] == '/' {
		return relativePath
	}
	if basePath == "" || basePath == "0" {
		throw("JsonSchema\\Exception\\UriResolverException", fmt.Sprintf("Unable to resolve URI '%s' from base '%s'", relativePath, basePath))
	}
	dirname := basePath
	if basePath[len(basePath)-1] != '/' {
		dirname = phpDirname(basePath)
	}
	combined := strings.TrimRight(dirname, "/") + "/" + strings.TrimLeft(relativePath, "/")
	var collapsed []string
	for _, seg := range strings.Split(combined, "/") {
		if seg == ".." {
			if len(collapsed) <= 1 {
				throw("JsonSchema\\Exception\\UriResolverException", fmt.Sprintf("Unable to resolve URI '%s' from base '%s'", relativePath, basePath))
			}
			collapsed = collapsed[:len(collapsed)-1]
		} else {
			collapsed = append(collapsed, seg)
		}
	}
	return strings.Join(collapsed, "/")
}

// ---- SchemaStorage / UriRetriever ----

type schemaStorage struct {
	schemas map[string]Value
	order   []string
	sources map[string][]byte
	cache   map[string]Value
}

func newSchemaStorage(sources map[string][]byte) *schemaStorage {
	return &schemaStorage{schemas: map[string]Value{}, sources: sources, cache: map[string]Value{}}
}

func (s *schemaStorage) set(id string, schema Value) {
	if _, ok := s.schemas[id]; !ok {
		s.order = append(s.order, id)
	}
	s.schemas[id] = schema
}

func findSchemaIDInObject(schema *phpjson.Map) *string {
	if v, ok := schema.Get("id"); ok {
		if str, ok := v.(string); ok {
			return &str
		}
	}
	if v, ok := schema.Get("$id"); ok {
		if str, ok := v.(string); ok {
			return &str
		}
	}
	return nil
}

// arrayToObject mirrors BaseConstraint::arrayToObjectRecursive for decoded
// JSON arrays: (object) json_decode(json_encode($array)).
func arrayToObject(a []Value) *phpjson.Map {
	m := phpjson.NewMap()
	for i, x := range a {
		m.Set(fmt.Sprint(i), x)
	}
	return m
}

func (s *schemaStorage) addSchema(id string, schema Value) {
	if schema == nil && id != internalProvidedSchemaURI {
		schema = s.retrieve(id, nil)
	}
	if a, ok := schema.([]Value); ok {
		schema = arrayToObject(a)
	}
	if m, ok := schema.(*phpjson.Map); ok {
		if idv, ok := m.Get("id"); ok {
			if idv == draft4 {
				if props := phpjson.AsMap(get(m, "properties")); props != nil {
					if idp := phpjson.AsMap(get(props, "id")); idp != nil {
						idp.Set("format", "uri-reference")
					}
				}
			} else if idv == draft3 {
				if props := phpjson.AsMap(get(m, "properties")); props != nil {
					if idp := phpjson.AsMap(get(props, "id")); idp != nil {
						idp.Set("format", "uri-reference")
					}
					if rp := phpjson.AsMap(get(props, "$ref")); rp != nil {
						rp.Set("format", "uri-reference")
					}
				}
			}
		}
	}
	baseID := id
	if m, ok := schema.(*phpjson.Map); ok {
		if rootID := findSchemaIDInObject(m); rootID != nil && *rootID != id {
			baseID = resolveURI(*rootID, &id)
		}
	}
	s.scanMembersForSubschemas(schema, baseID, "")
	s.expandRefs(schema, &id, nil)
	s.set(id, schema)
}

func (s *schemaStorage) expandRefs(schema Value, parentID *string, propertyStack []string) {
	m, ok := schema.(*phpjson.Map)
	if !ok {
		if a, ok := schema.([]Value); ok {
			for _, member := range a {
				s.expandRefs(member, parentID, nil)
			}
		}
		return
	}
	if ref, ok := m.Get("$ref"); ok {
		if rs, ok := ref.(string); ok {
			m.Set("$ref", newJSONPointer(resolveURI(rs, parentID)).String())
		}
	}
	parentProperty := ""
	if len(propertyStack) > 0 {
		parentProperty = propertyStack[len(propertyStack)-1]
	}
	for _, name := range m.Keys() {
		if (name == "enum" || name == "const") && !inStrings(parentProperty, schemaMapKeywords) {
			continue
		}
		member, _ := m.Get(name)
		childID := parentID
		if sid := findSchemaIDInObject(m); sid != nil && (childID == nil || *childID != *sid) {
			r := resolveURI(*sid, childID)
			childID = &r
		}
		stack := append(append([]string(nil), propertyStack...), name)
		s.expandRefs(member, childID, stack)
	}
}

func (s *schemaStorage) scanMembersForSubschemas(schema Value, parentID, parentProperty string) {
	switch schema.(type) {
	case *phpjson.Map, []Value:
	default:
		return
	}
	for _, p := range iterPairs(schema) {
		name, isStr := p.k.(string)
		if isStr && (name == "enum" || name == "const") && !inStrings(parentProperty, schemaMapKeywords) {
			continue
		}
		if a, ok := p.v.([]Value); ok {
			for _, item := range a {
				s.registerAndScanSubschema(item, parentID, toStr(p.k))
			}
			continue
		}
		s.registerAndScanSubschema(p.v, parentID, toStr(p.k))
	}
}

func (s *schemaStorage) registerAndScanSubschema(sub Value, parentID, parentProperty string) {
	m, ok := sub.(*phpjson.Map)
	if !ok {
		s.scanMembersForSubschemas(sub, parentID, parentProperty)
		return
	}
	childID := parentID
	if sid := findSchemaIDInObject(m); sid != nil {
		childID = resolveURI(*sid, &parentID)
		s.set(childID, m)
	}
	s.scanMembersForSubschemas(m, childID, parentProperty)
}

func (s *schemaStorage) getSchema(id string) Value {
	if _, ok := s.schemas[id]; !ok {
		s.addSchema(id, nil)
	}
	return s.schemas[id]
}

func urldecode(str string) string {
	var sb strings.Builder
	for i := 0; i < len(str); i++ {
		c := str[i]
		switch {
		case c == '+':
			sb.WriteByte(' ')
		case c == '%' && i+2 < len(str) && isHex(rune(str[i+1])) && isHex(rune(str[i+2])):
			sb.WriteByte(byte(hexVal(rune(str[i+1]))<<4 | hexVal(rune(str[i+2]))))
			i += 2
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

func (s *schemaStorage) resolveRef(ref string, stack []*phpjson.Map) Value {
	jp := newJSONPointer(ref)
	if jp.filename == "" {
		throw("JsonSchema\\Exception\\UnresolvableJsonPointerException", fmt.Sprintf("Could not resolve fragment '%s': no file is defined", jp.propertyPathAsString()))
	}
	refSchema := s.getSchema(jp.filename)
	for _, path := range jp.paths {
		path = urldecode(path)
		switch t := refSchema.(type) {
		case *phpjson.Map:
			if v, ok := t.Get(path); ok {
				refSchema = s.resolveRefSchema(v, stack)
				continue
			}
		case []Value:
			if idx, ok := phpjson.IntKey(path); ok && idx >= 0 && idx < int64(len(t)) {
				refSchema = s.resolveRefSchema(t[idx], stack)
				continue
			}
		}
		throw("JsonSchema\\Exception\\UnresolvableJsonPointerException", fmt.Sprintf("File: %s is found, but could not resolve fragment: %s", jp.filename, jp.propertyPathAsString()))
	}
	return refSchema
}

func (s *schemaStorage) resolveRefSchema(refSchema Value, stack []*phpjson.Map) Value {
	if m, ok := refSchema.(*phpjson.Map); ok {
		if ref, ok := m.Get("$ref"); ok {
			if rs, ok := ref.(string); ok {
				for _, x := range stack {
					if x == m {
						throw("JsonSchema\\Exception\\UnresolvableJsonPointerException", fmt.Sprintf("Dereferencing a pointer to %s results in an infinite loop", rs))
					}
				}
				next := append(append([]*phpjson.Map(nil), stack...), m)
				return s.resolveRef(rs, next)
			}
		}
		if keys := m.Keys(); len(keys) == 1 && keys[0] == "" {
			v, _ := m.Get("")
			return v
		}
	}
	return refSchema
}

// retrieve mirrors UriRetriever::retrieve.
func (s *schemaStorage) retrieve(uri string, baseURI *string) Value {
	resolved := resolveURI(uri, baseURI)
	fetch := resolved
	parts := parseURI(resolved)
	if parts.fragment != nil {
		parts.fragment = nil
		fetch = generateURI(parts)
	}
	if strings.HasPrefix(fetch, "http://json-schema.org/") || strings.HasPrefix(fetch, "https://json-schema.org/") || strings.HasPrefix(fetch, "package://") {
		throw("JsonSchema\\Exception\\ResourceNotFoundException", "gocomposer: bundled meta-schemas are not available: "+fetch)
	}
	schema := s.loadSchema(fetch)
	dialect := ""
	hasDialect := false
	if d, ok := mapGet(schema, "$schema"); ok && d != nil {
		hasDialect = true
		dialect = strings.TrimRight(toStr(d), "#")
	}
	usesDollarID := !(hasDialect && (dialect == strings.TrimRight(draft3, "#") || dialect == strings.TrimRight(draft4, "#")))
	schema = resolvePointer(schema, resolved)
	if m, ok := schema.(*phpjson.Map); ok {
		if usesDollarID {
			m.Set("$id", resolved)
		} else {
			m.Set("id", resolved)
		}
	}
	return schema
}

func resolvePointer(schema Value, uri string) Value {
	parsed := parseURI(uri)
	if parsed.fragment == nil || *parsed.fragment == "" || *parsed.fragment == "0" {
		return schema
	}
	for _, el := range strings.Split(*parsed.fragment, "/") {
		if el == "" || el == "0" {
			continue
		}
		el = strings.ReplaceAll(el, "~1", "/")
		el = strings.ReplaceAll(el, "~0", "~")
		v, ok := mapGet(schema, el)
		if !ok || !truthy(v) {
			throw("JsonSchema\\Exception\\ResourceNotFoundException", "Fragment \""+*parsed.fragment+"\" not found in "+uri)
		}
		schema = v
		if _, ok := schema.(*phpjson.Map); !ok {
			throw("JsonSchema\\Exception\\ResourceNotFoundException", "Fragment part \""+el+"\" is no object  in "+uri)
		}
	}
	return schema
}

func (s *schemaStorage) loadSchema(fetch string) Value {
	if v, ok := s.cache[fetch]; ok {
		return v
	}
	data, ok := s.sources[fetch]
	if !ok {
		if strings.HasPrefix(fetch, "file://") {
			b, err := os.ReadFile(strings.TrimPrefix(fetch, "file://"))
			if err != nil {
				throw("JsonSchema\\Exception\\ResourceNotFoundException", "file_get_contents("+fetch+"): Failed to open stream: No such file or directory")
			}
			data = b
		} else {
			throw("JsonSchema\\Exception\\ResourceNotFoundException", "gocomposer: retrieving "+fetch+" is not supported")
		}
	}
	v, err := phpjson.Decode(data)
	if err != nil {
		throw("JsonSchema\\Exception\\JsonDecodingException", "JSON syntax is malformed")
	}
	s.cache[fetch] = v
	return v
}
