// Ported for gocomposer (MIT) from justinrainbow/json-schema 6.13.1
// (MIT, Copyright (c) 2016 justinrainbow/json-schema contributors):
// src/JsonSchema/Constraints/{BaseConstraint,Constraint,SchemaConstraint,
// UndefinedConstraint,TypeConstraint,ObjectConstraint,CollectionConstraint,
// StringConstraint,NumberConstraint,EnumConstraint,ConstConstraint}.php and
// src/JsonSchema/ConstraintError.php.
//
// This is the code path Validator::validate() takes in the default check mode
// (CHECK_MODE_NORMAL), which is what Composer uses: SchemaConstraint ->
// UndefinedConstraint and friends. The Drafts/Draft06.. constraint sets are
// only used with CHECK_MODE_STRICT and are not ported.
package jsonschema

import (
	"errors"
	"math"
	"strings"

	"gocomposer/internal/phpjson"
)

// Constraint error names (ConstraintError constants) and their messages.
var constraintMessages = map[string]string{
	"additionalItems":   "The item %s[%s] is not defined and the definition does not allow additional items",
	"additionalProp":    "The property %s is not defined and the definition does not allow additional properties",
	"allOf":             "Failed to match all schemas",
	"anyOf":             "Failed to match at least one schema",
	"dependencies":      "%s depends on %s, which is missing",
	"disallow":          "Disallowed value was matched",
	"divisibleBy":       "Is not divisible by %d",
	"enum":              "Does not have a value in the enumeration %s",
	"const":             "Does not have a value equal to %s",
	"exclusiveMinimum":  "Must have a minimum value greater than %d",
	"exclusiveMaximum":  "Must have a maximum value less than %d",
	"false":             "Boolean schema false",
	"colorFormat":       "Invalid color",
	"dateFormat":        "Invalid date %s, expected format YYYY-MM-DD",
	"dateTimeFormat":    "Invalid date-time %s, expected format YYYY-MM-DDThh:mm:ssZ or YYYY-MM-DDThh:mm:ss+hh:mm",
	"dateUtcFormat":     "Invalid time %s, expected integer of milliseconds since Epoch",
	"emailFormat":       "Invalid email",
	"styleHostName":     "Invalid hostname",
	"ipFormat":          "Invalid IP address",
	"phoneFormat":       "Invalid phone number",
	"regexFormat":       "Invalid regex format %s",
	"styleFormat":       "Invalid style",
	"timeFormat":        "Invalid time %s, expected format hh:mm:ss",
	"urlFormat":         "Invalid URL format",
	"invalidSchema":     "Schema is not valid",
	"maxLength":         "Must be at most %d characters long",
	"minLength":         "Must be at least %d characters long",
	"maxItems":          "There must be a maximum of %d items in the array, %d found",
	"maximum":           "Must have a maximum value less than or equal to %d",
	"minItems":          "There must be a minimum of %d items in the array, %d found",
	"minimum":           "Must have a minimum value greater than or equal to %d",
	"missingMaximum":    "Use of exclusiveMaximum requires presence of maximum",
	"missingMinimum":    "Use of exclusiveMinimum requires presence of minimum",
	"multipleOf":        "Must be a multiple of %s",
	"not":               "Matched a schema which it should not",
	"oneOf":             "Failed to match exactly one schema",
	"required":          "The property %s is required",
	"requires":          "The presence of the property %s requires that %s also be present",
	"pattern":           "Does not match the regex pattern %s",
	"pregrex":           "The pattern %s is invalid",
	"minProperties":     "Must contain a minimum of %d properties",
	"maxProperties":     "Must contain no more than %d properties",
	"type":              "%s value found, but %s is required",
	"uniqueItems":       "There are no duplicates allowed in the array",
	"propertyNames":     "Property name %s is invalid",
	"contains":          "Does not have a value valid to contains schema",
	"durationFormat":    "Invalid duration %s, expected ISO 8601 duration such as P1DT12H",
	"jsonPointerFormat": "Invalid JSON pointer",
	"uriTemplateFormat": "Invalid URI template format",
	"urlRefFormat":      "Invalid URL reference format",
	"uuidFormat":        "Invalid UUID format",
	"maxContains":       "There must be a maximum of %d valid items in the array, %d found",
	"minContains":       "There must be a minimum of %d valid items in the array, %d found",
	"contentMediaType":  "Value is not valid with content media type",
	"contentEncoding":   "Value is not valid with content encoding",
}

var errNoEndDelimiter = errors.New("no ending delimiter")

type param struct {
	Name  string
	Value Value
}

// Error is one entry of Validator::getErrors().
type Error struct {
	// Property is the dotted property path ("require.foo", "authors[0].name").
	Property string
	// Pointer is the JSON pointer of the instance ("/require/foo"), "" for the root.
	Pointer string
	// Message is the human readable message (ucfirst'ed).
	Message string
	// Constraint is the ConstraintError name ("type", "required", "additionalProp", ...).
	Constraint string
	// Params are the constraint parameters ($error['constraint']['params']) in order.
	Params []param
}

// ParamsJSON returns json_encode($error['constraint']['params']).
func (e Error) ParamsJSON() string { return encodeParams(e.Params) }

// Param returns the named constraint parameter.
func (e Error) Param(name string) (Value, bool) {
	for _, p := range e.Params {
		if p.Name == name {
			return p.Value, true
		}
	}
	return nil, false
}

// inst models the error list of one constraint instance.
type inst struct{ errs []Error }

func cloneErrs(e []Error) []Error { return append([]Error(nil), e...) }

type validator struct {
	st            *schemaStorage
	disableFormat bool
}

func convertPointerToPropertyPath(p jsonPointer) string {
	var sb strings.Builder
	for _, seg := range p.paths {
		if isNumeric(seg) {
			sb.WriteString("[" + vsprintf("%d", []Value{seg}) + "]")
		} else {
			sb.WriteString("." + seg)
		}
	}
	return strings.Trim(sb.String(), ".")
}

func (v *validator) addError(in *inst, name string, path jsonPointer, params ...param) {
	args := make([]Value, len(params))
	for i, p := range params {
		val := normalize(p.Value)
		if isScalar(val) {
			if b, ok := val.(bool); ok {
				if b {
					args[i] = "true"
				} else {
					args[i] = "false"
				}
			} else {
				args[i] = val
			}
		} else {
			args[i] = jsonEncode(val)
		}
	}
	in.errs = append(in.errs, Error{
		Property:   convertPointerToPropertyPath(path),
		Pointer:    strings.TrimLeft(path.String(), "#"),
		Message:    ucfirst(vsprintf(constraintMessages[name], args)),
		Constraint: name,
		Params:     params,
	})
}

// incrementPath mirrors Constraint::incrementPath.
func incrementPath(path jsonPointer, i Value) jsonPointer {
	if i == nil {
		return path
	}
	if s, ok := i.(string); ok && s == "" {
		return path
	}
	return path.withAppended(toStr(i))
}

// checkUndefined mirrors Constraint::checkUndefined.
func (v *validator) checkUndefined(in *inst, value, schema Value, path jsonPointer, i Value) {
	sub := &inst{}
	v.undefinedCheck(sub, value, v.st.resolveRefSchema(schema, nil), path, i)
	in.errs = append(in.errs, sub.errs...)
}

// ---- UndefinedConstraint ----

func (v *validator) undefinedCheck(in *inst, value, schema Value, path jsonPointer, i Value) {
	sm, ok := schema.(*phpjson.Map)
	if !ok {
		return
	}
	path = incrementPath(path, i)
	v.validateCommonProperties(in, value, sm, path, i)
	v.validateOfProperties(in, value, sm, path)
	v.validateTypes(in, value, sm, path, i)
}

func (v *validator) validateTypes(in *inst, value Value, schema *phpjson.Map, path jsonPointer, i Value) {
	if isArray(value) {
		sub := &inst{}
		v.collectionCheck(sub, value.([]Value), schema, path, i)
		in.errs = append(in.errs, sub.errs...)
	}
	if looseIsObject(value) {
		var props, addl, pattern Value
		if isset(schema, "properties") {
			props = get(schema, "properties")
		}
		if isset(schema, "additionalProperties") {
			addl = get(schema, "additionalProperties")
		}
		if isset(schema, "patternProperties") {
			pattern = get(schema, "patternProperties")
		}
		sub := &inst{}
		v.objectCheck(sub, value, schema, path, props, addl, pattern)
		in.errs = append(in.errs, sub.errs...)
	}
	if s, ok := value.(string); ok {
		sub := &inst{}
		v.stringCheck(sub, s, schema, path, i)
		in.errs = append(in.errs, sub.errs...)
	}
	if isNumeric(value) {
		sub := &inst{}
		v.numberCheck(sub, value, schema, path, i)
		in.errs = append(in.errs, sub.errs...)
	}
	if isset(schema, "enum") {
		sub := &inst{}
		v.enumCheck(sub, value, schema, path)
		in.errs = append(in.errs, sub.errs...)
	}
	if isset(schema, "const") {
		sub := &inst{}
		v.constCheck(sub, value, schema, path)
		in.errs = append(in.errs, sub.errs...)
	}
}

func propertyExists(obj Value, prop Value) bool {
	s, ok := prop.(string)
	if !ok {
		throw("TypeError", "property_exists(): Argument #2 ($property) must be of type string, "+debugType(prop)+" given")
	}
	m, ok := obj.(*phpjson.Map)
	if !ok {
		return false
	}
	return m.Has(s)
}

func (v *validator) validateCommonProperties(in *inst, value Value, schema *phpjson.Map, path jsonPointer, i Value) {
	if isset(schema, "extends") {
		ext := get(schema, "extends")
		switch t := ext.(type) {
		case string:
			throw("JsonSchema\\Exception\\RuntimeException", "gocomposer: draft-03 \"extends\" by URI is not supported")
		case []Value:
			for _, e := range t {
				v.checkUndefined(in, value, e, path, i)
			}
		default:
			v.checkUndefined(in, value, ext, path, i)
		}
	}

	if isObject(value) {
		req, hasReq := mapGet(schema, "required")
		hasReq = hasReq && req != nil
		_, reqIsArray := req.([]Value)
		if !isUndef(value) && hasReq && reqIsArray {
			for _, r := range req.([]Value) {
				if !propertyExists(value, r) {
					v.addError(in, "required", incrementPath(path, r), param{"property", r})
				}
			}
		} else if hasReq && !reqIsArray {
			if truthy(req) && isUndef(value) {
				var last Value = false
				if len(path.paths) > 0 {
					last = path.paths[len(path.paths)-1]
				}
				v.addError(in, "required", path, param{"property", last})
			}
		} else if isUndef(value) {
			return
		}
	}

	if !isUndef(value) {
		v.checkType(in, value, schema, path, i)
	}

	if isset(schema, "disallow") {
		initErrors := cloneErrs(in.errs)
		typeSchema := phpjson.NewMap()
		typeSchema.Set("type", get(schema, "disallow"))
		v.checkType(in, value, typeSchema, path, nil)
		if len(in.errs) == len(initErrors) {
			v.addError(in, "disallow", path)
		} else {
			in.errs = initErrors
		}
	}

	if isset(schema, "not") {
		initErrors := cloneErrs(in.errs)
		v.checkUndefined(in, value, get(schema, "not"), path, i)
		if len(in.errs) == len(initErrors) {
			v.addError(in, "not", path)
		} else {
			in.errs = initErrors
		}
	}

	if isset(schema, "dependencies") && isObject(value) {
		v.validateDependencies(in, value, get(schema, "dependencies"), path)
	}
}

func (v *validator) validateOfProperties(in *inst, value Value, schema *phpjson.Map, path jsonPointer) {
	if isUndef(value) {
		return
	}
	if isset(schema, "allOf") {
		valid := true
		for _, s := range iterValues(get(schema, "allOf")) {
			init := len(in.errs)
			v.checkUndefined(in, value, s, path, "")
			valid = valid && len(in.errs) == init
		}
		if !valid {
			v.addError(in, "allOf", path)
		}
	}
	if isset(schema, "anyOf") {
		valid := false
		startErrors := cloneErrs(in.errs)
		for _, s := range iterValues(get(schema, "anyOf")) {
			init := len(in.errs)
			v.checkUndefined(in, value, s, path, "")
			valid = len(in.errs) == init
			if valid {
				break
			}
		}
		if !valid {
			v.addError(in, "anyOf", path)
		} else {
			in.errs = startErrors
		}
	}
	if isset(schema, "oneOf") {
		var allErrors []Error
		matched := 0
		startErrors := cloneErrs(in.errs)
		for _, s := range iterValues(get(schema, "oneOf")) {
			in.errs = nil
			v.checkUndefined(in, value, s, path, "")
			if len(in.errs) == 0 {
				matched++
			}
			allErrors = append(allErrors, in.errs...)
		}
		if matched != 1 {
			merged := append(cloneErrs(allErrors), startErrors...)
			in.errs = append(cloneErrs(in.errs), merged...)
			v.addError(in, "oneOf", path)
		} else {
			in.errs = startErrors
		}
	}
}

func (v *validator) validateDependencies(in *inst, value Value, deps Value, path jsonPointer) {
	for _, d := range iterPairs(deps) {
		if !propertyExists(value, d.k) {
			continue
		}
		switch dep := d.v.(type) {
		case string:
			if !propertyExists(value, dep) {
				v.addError(in, "dependencies", path, param{"key", d.k}, param{"dependency", dep})
			}
		case []Value:
			for _, x := range dep {
				if !propertyExists(value, x) {
					v.addError(in, "dependencies", path, param{"key", d.k}, param{"dependency", dep})
				}
			}
		case *phpjson.Map:
			v.checkUndefined(in, value, dep, path, "")
		}
	}
}

// ---- TypeConstraint ----

var typeWording = map[string]Value{
	"integer": "an integer",
	"number":  "a number",
	"boolean": "a boolean",
	"object":  "an object",
	"array":   "an array",
	"string":  "a string",
	"null":    "a null",
	"any":     nil,
	"0":       nil,
}

func (v *validator) checkType(in *inst, value Value, schema *phpjson.Map, path jsonPointer, i Value) {
	sub := &inst{}
	v.typeCheck(sub, value, schema, path)
	in.errs = append(in.errs, sub.errs...)
}

func (v *validator) typeCheck(in *inst, value Value, schema *phpjson.Map, path jsonPointer) {
	var typ Value
	if isset(schema, "type") {
		typ = get(schema, "type")
	}
	valid := false
	var wording []Value
	switch t := typ.(type) {
	case []Value:
		v.validateTypesArray(value, t, &wording, &valid, path)
	case *phpjson.Map:
		v.checkUndefined(in, value, t, path, nil)
		return
	default:
		valid = v.validateType(value, typ)
	}
	if !valid {
		if _, isArr := typ.([]Value); !isArr {
			validateTypeNameWording(typ)
			wording = append(wording, typeWording[wordingKey(typ)])
		}
		v.addError(in, "type", path, param{"found", gettype(value)}, param{"expected", implodeWith(wording, ", ", "or")})
	}
}

func wordingKey(t Value) string {
	if k, ok := arrayKey(t); ok {
		if k == 0 {
			return "0"
		}
		return toStr(k)
	}
	return toStr(t)
}

func validateTypeNameWording(t Value) {
	if _, ok := typeWording[wordingKey(t)]; ok {
		return
	}
	throw("UnexpectedValueException", "No wording for "+varExport(t)+" available, expected wordings are: [an integer, a number, a boolean, an object, an array, a string, a null]")
}

func varExport(v Value) string {
	switch t := normalize(v).(type) {
	case nil:
		return "NULL"
	case bool:
		if t {
			return "true"
		}
		return "false"
	case string:
		return "'" + strings.ReplaceAll(strings.ReplaceAll(t, `\`, `\\`), `'`, `\'`) + "'"
	case int64, float64:
		s := toStr(t)
		if f, ok := t.(float64); ok && !strings.ContainsAny(s, ".EN") {
			_ = f
			s += ".0"
		}
		return s
	}
	return jsonEncode(v)
}

func (v *validator) validateTypesArray(value Value, types []Value, wording *[]Value, valid *bool, path jsonPointer) {
	for _, tp := range types {
		if *valid {
			return
		}
		if m, ok := tp.(*phpjson.Map); ok {
			sub := &inst{}
			subSchema := phpjson.NewMap()
			subSchema.Set("type", m)
			v.typeCheck(sub, value, subSchema, path)
			*valid = len(sub.errs) == 0
			*wording = append(*wording, typeWording["object"])
		} else {
			validateTypeNameWording(tp)
			*wording = append(*wording, typeWording[wordingKey(tp)])
			if !*valid {
				*valid = v.validateType(value, tp)
			}
		}
	}
}

func implodeWith(elements []Value, delim, listEnd string) string {
	strs := make([]string, len(elements))
	for i, e := range elements {
		strs[i] = toStr(e)
	}
	if len(elements) < 2 || elements[1] == nil {
		return strings.Join(strs, delim)
	}
	first := strings.Join(strs[:len(strs)-1], delim)
	return first + " " + listEnd + " " + strs[len(strs)-1]
}

func (v *validator) validateType(value Value, typ Value) bool {
	if !truthy(typ) {
		return true
	}
	switch {
	case typ == "any":
		return true
	case typ == "object":
		return isObject(value)
	case typ == "array":
		return isArray(value)
	case typ == "integer":
		_, ok := normalize(value).(int64)
		return ok
	case typ == "number":
		switch normalize(value).(type) {
		case int64, float64:
			return true
		}
		return false
	case typ == "boolean":
		_, ok := value.(bool)
		return ok
	case typ == "string":
		_, ok := value.(string)
		return ok
	case typ == "null":
		return value == nil
	}
	var vs string
	if isObject(value) {
		vs = "object"
	} else {
		vs = toStr(value)
	}
	throw("JsonSchema\\Exception\\InvalidArgumentException", vs+" is an invalid type for "+toStr(typ))
	return false
}

// ---- ObjectConstraint ----

func getProperty(element Value, prop Value) (Value, bool) {
	switch t := element.(type) {
	case []Value:
		if k, ok := arrayKey(prop); ok && k >= 0 && k < int64(len(t)) {
			return t[k], true
		}
	case *phpjson.Map:
		return t.Get(toStr(prop))
	}
	return nil, false
}

func inArrayLoose(needle Value, list []Value) bool {
	for _, x := range list {
		if looseEqual(needle, x) {
			return true
		}
	}
	return false
}

func (v *validator) objectCheck(in *inst, element Value, schema *phpjson.Map, path jsonPointer, properties, additionalProp, patternProperties Value) {
	if isUndef(element) {
		return
	}
	var matches []Value
	if truthy(patternProperties) {
		matches = v.validatePatternProperties(in, element, path, patternProperties)
	}
	if truthy(properties) {
		v.validateProperties(in, element, properties, path)
	}
	v.validateElement(in, element, matches, schema, path, properties, additionalProp)
}

func jsonPatternToPhpRegex(pattern Value) (*reProg, error) {
	s, ok := pattern.(string)
	if !ok {
		throw("TypeError", "JsonSchema\\Constraints\\BaseConstraint::jsonPatternToPhpRegex(): Argument #1 ($pattern) must be of type string, "+debugType(pattern)+" given")
	}
	return compileDelimited("~" + strings.ReplaceAll(s, "~", `\~`) + "~u")
}

// compileDelimited mirrors how preg_* split a delimited PHP regex into the
// pattern body and its modifiers.
func compileDelimited(full string) (*reProg, error) {
	delim := full[0]
	i := 1
	for i < len(full) {
		if full[i] == '\\' && i+1 < len(full) {
			i += 2
			continue
		}
		if full[i] == delim {
			break
		}
		i++
	}
	if i >= len(full) {
		return nil, errNoEndDelimiter
	}
	return compileRegex(full[1:i], full[i+1:])
}

func (v *validator) validatePatternProperties(in *inst, element Value, path jsonPointer, patternProperties Value) []Value {
	var matches []Value
	for _, pp := range iterPairs(patternProperties) {
		re, err := jsonPatternToPhpRegex(pp.k)
		if err == nil && pregMatch(re, "") == -1 {
			err = errReLimit
		}
		if err != nil {
			v.addError(in, "pregrex", path, param{"pregex", pp.k})
			continue
		}
		for _, el := range iterPairs(element) {
			if pregMatch(re, toStr(el.k)) == 1 {
				matches = append(matches, el.k)
				sch := pp.v
				if !truthy(sch) {
					sch = phpjson.NewMap()
				}
				v.checkUndefined(in, el.v, sch, path, el.k)
			}
		}
	}
	return matches
}

func (v *validator) validateElement(in *inst, element Value, matches []Value, schema *phpjson.Map, path jsonPointer, properties, additionalProp Value) {
	v.validateMinMaxConstraint(in, element, schema, path)
	for _, el := range iterPairs(element) {
		definition, _ := getProperty(properties, el.k)
		inMatches := inArrayLoose(el.k, matches)
		if !inMatches && additionalProp == false && el.k != "$schema" && !truthy(definition) {
			v.addError(in, "additionalProp", path, param{"property", el.k})
		}
		if !inMatches && truthy(additionalProp) && !truthy(definition) {
			if additionalProp == true {
				v.checkUndefined(in, el.v, nil, path, el.k)
			} else {
				v.checkUndefined(in, el.v, additionalProp, path, el.k)
			}
		}
		require, _ := getProperty(definition, "requires")
		if truthy(require) {
			rv, _ := getProperty(element, require)
			if !truthy(rv) {
				v.addError(in, "requires", path, param{"property", el.k}, param{"requiredProperty", require})
			}
		}
		property, ok := getProperty(element, el.k)
		if !ok {
			property = undefined
		}
		if isObject(property) {
			target := property
			if isUndef(property) {
				target = element
			}
			v.validateMinMaxConstraint(in, target, definition, path)
		}
	}
}

func (v *validator) validateProperties(in *inst, element Value, properties Value, path jsonPointer) {
	for _, p := range iterPairs(properties) {
		property, ok := getProperty(element, p.k)
		if !ok {
			property = undefined
		}
		if def, ok := p.v.(*phpjson.Map); ok {
			v.checkUndefined(in, property, def, path, p.k)
		}
	}
}

func (v *validator) validateMinMaxConstraint(in *inst, element Value, def Value, path jsonPointer) {
	if !isObject(element) {
		return
	}
	count := int64(0)
	if m, ok := element.(*phpjson.Map); ok {
		count = int64(m.Len())
	}
	if isset(def, "minProperties") {
		if mp, ok := normalize(get(def, "minProperties")).(int64); ok {
			if count < max(0, mp) {
				v.addError(in, "minProperties", path, param{"minProperties", mp})
			}
		}
	}
	if isset(def, "maxProperties") {
		if mp, ok := normalize(get(def, "maxProperties")).(int64); ok {
			if count > max(0, mp) {
				v.addError(in, "maxProperties", path, param{"maxProperties", mp})
			}
		}
	}
}

// ---- CollectionConstraint ----

func (v *validator) collectionCheck(in *inst, value []Value, schema *phpjson.Map, path jsonPointer, i Value) {
	n := int64(len(value))
	if isset(schema, "minItems") && compare(n, get(schema, "minItems")) < 0 {
		v.addError(in, "minItems", path, param{"minItems", get(schema, "minItems")}, param{"found", n})
	}
	if isset(schema, "maxItems") && compare(n, get(schema, "maxItems")) > 0 {
		v.addError(in, "maxItems", path, param{"maxItems", get(schema, "maxItems")}, param{"found", n})
	}
	if isset(schema, "uniqueItems") && truthy(get(schema, "uniqueItems")) {
	outer:
		for x := 0; x < len(value)-1; x++ {
			for y := x + 1; y < len(value); y++ {
				if deepEqual(value[x], value[y]) {
					v.addError(in, "uniqueItems", path)
					break outer
				}
			}
		}
	}
	v.validateItems(in, value, schema, path, i)
}

func (v *validator) validateItems(in *inst, value []Value, schema *phpjson.Map, path jsonPointer, i Value) {
	if !isset(schema, "items") {
		return
	}
	items := get(schema, "items")
	if items == true {
		return
	}
	if _, ok := items.(*phpjson.Map); ok {
		var secondErrors []Error
		secondSet := false
		for k, x := range value {
			initErrors := cloneErrs(in.errs)
			v.checkUndefined(in, x, items, path, int64(k))
			if len(initErrors) < len(in.errs) && isset(schema, "additionalItems") && get(schema, "additionalItems") != false {
				secondErrors = cloneErrs(in.errs)
				secondSet = true
				v.checkUndefined(in, x, get(schema, "additionalItems"), path, int64(k))
			}
			if secondSet && len(secondErrors) < len(in.errs) {
				in.errs = cloneErrs(secondErrors)
			} else if secondSet && len(secondErrors) == len(in.errs) {
				in.errs = initErrors
			}
		}
		return
	}
	list, ok := items.([]Value)
	if !ok {
		if len(value) > 0 {
			throw("TypeError", "array_key_exists(): Argument #2 ($array) must be of type array, "+debugType(items)+" given")
		}
		return
	}
	for k, x := range value {
		if k < len(list) {
			v.checkUndefined(in, x, list[k], path, int64(k))
			continue
		}
		if ai, ok := mapGet(schema, "additionalItems"); ok {
			if ai != false {
				v.checkUndefined(in, x, ai, path, int64(k))
			} else {
				v.addError(in, "additionalItems", path, param{"item", i}, param{"property", int64(k)}, param{"additionalItems", ai})
			}
		} else {
			v.checkUndefined(in, x, phpjson.NewMap(), path, int64(k))
		}
	}
	if len(value) > 0 {
		for k := len(value); k < len(list); k++ {
			v.checkUndefined(in, undefined, list[k], path, int64(k))
		}
	}
}

// ---- StringConstraint ----

func (v *validator) stringCheck(in *inst, element string, schema *phpjson.Map, path jsonPointer, i Value) {
	length := int64(len([]rune(element)))
	if isset(schema, "maxLength") && compare(length, get(schema, "maxLength")) > 0 {
		v.addError(in, "maxLength", path, param{"maxLength", get(schema, "maxLength")})
	}
	if isset(schema, "minLength") && compare(length, get(schema, "minLength")) < 0 {
		v.addError(in, "minLength", path, param{"minLength", get(schema, "minLength")})
	}
	if isset(schema, "pattern") {
		re, err := jsonPatternToPhpRegex(get(schema, "pattern"))
		if err != nil || pregMatch(re, element) != 1 {
			v.addError(in, "pattern", path, param{"pattern", get(schema, "pattern")})
		}
	}
	v.formatCheck(in, element, schema, path)
}

// ---- NumberConstraint ----

func (v *validator) numberCheck(in *inst, element Value, schema *phpjson.Map, path jsonPointer, i Value) {
	if isset(schema, "exclusiveMinimum") {
		if isset(schema, "minimum") {
			min := get(schema, "minimum")
			if truthy(get(schema, "exclusiveMinimum")) && compare(element, min) <= 0 {
				v.addError(in, "exclusiveMinimum", path, param{"minimum", min})
			} else if compare(element, min) < 0 {
				v.addError(in, "minimum", path, param{"minimum", min})
			}
		} else {
			v.addError(in, "missingMinimum", path)
		}
	} else if isset(schema, "minimum") && compare(element, get(schema, "minimum")) < 0 {
		v.addError(in, "minimum", path, param{"minimum", get(schema, "minimum")})
	}

	if isset(schema, "exclusiveMaximum") {
		if isset(schema, "maximum") {
			max := get(schema, "maximum")
			if truthy(get(schema, "exclusiveMaximum")) && compare(element, max) >= 0 {
				v.addError(in, "exclusiveMaximum", path, param{"maximum", max})
			} else if compare(element, max) > 0 {
				v.addError(in, "maximum", path, param{"maximum", max})
			}
		} else {
			v.addError(in, "missingMaximum", path)
		}
	} else if isset(schema, "maximum") && compare(element, get(schema, "maximum")) > 0 {
		v.addError(in, "maximum", path, param{"maximum", get(schema, "maximum")})
	}

	if isset(schema, "divisibleBy") && phpFmod(element, get(schema, "divisibleBy")) != 0 {
		v.addError(in, "divisibleBy", path, param{"divisibleBy", get(schema, "divisibleBy")})
	}
	if isset(schema, "multipleOf") && phpFmod(element, get(schema, "multipleOf")) != 0 {
		v.addError(in, "multipleOf", path, param{"multipleOf", get(schema, "multipleOf")})
	}
	v.formatCheck(in, element, schema, path)
}

func arithOperand(x Value, op string, other Value) float64 {
	switch t := normalize(x).(type) {
	case int64, float64:
		return toFloat(t)
	case bool:
		return toFloat(t)
	case nil:
		return 0
	case string:
		n, whole, ok := numericPrefix(t)
		if !ok {
			throw("TypeError", "Unsupported operand types: string "+op+" "+debugType(other))
		}
		_ = whole // leading-numeric strings only warn
		return toFloat(n)
	}
	throw("TypeError", "Unsupported operand types: "+debugType(x)+" "+op+" "+debugType(other))
	return 0
}

func phpRound(f float64) float64 { return math.Round(f) }

func phpFmod(n1, n2 Value) float64 {
	a := arithOperand(n1, "/", n2)
	b := arithOperand(n2, "/", n1)
	if b == 0 {
		throw("DivisionByZeroError", "Division by zero")
	}
	mod := a - phpRound(a/b)*b
	const precision = 0.0000000001
	if -precision < mod && mod < precision {
		return 0
	}
	return mod
}

// ---- EnumConstraint / ConstConstraint ----

func (v *validator) enumCheck(in *inst, element Value, schema *phpjson.Map, path jsonPointer) {
	if isUndef(element) && (!isset(schema, "required") || !truthy(get(schema, "required"))) {
		return
	}
	typ := gettype(element)
	for _, e := range iterValues(get(schema, "enum")) {
		if typ == gettype(e) && deepEqual(element, e) {
			return
		}
		if isNumeric(element) && isNumeric(e) && deepEqual(toFloat(element), toFloat(e)) {
			return
		}
	}
	v.addError(in, "enum", path, param{"enum", get(schema, "enum")})
}

func (v *validator) constCheck(in *inst, element Value, schema *phpjson.Map, path jsonPointer) {
	if isUndef(element) && (!isset(schema, "required") || !truthy(get(schema, "required"))) {
		return
	}
	c := get(schema, "const")
	if deepEqual(element, c) {
		return
	}
	v.addError(in, "const", path, param{"const", c})
}
