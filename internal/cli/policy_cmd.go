package cli

// Port of Composer\Command\PolicyCommand and Composer\Policy\SourceValidator
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano). Shell completion is not ported.

import (
	"strconv"
	"strings"

	"gocomposer/internal/jsonfile"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/policy"
)

func init() { registerCommand(policyCommand) }

func policyCommand() *Command {
	return &Command{Name: "policy", Run: runPolicy}
}

func runPolicy(ctx *Context) (int, error) {
	b, err := initBaseConfig(ctx)
	if err != nil {
		return 1, err
	}
	in := ctx.In
	action := strings.ToLower(in.Arg("action"))
	_, hasName := in.args["name"]
	_, hasArg1 := in.args["arg1"]
	_, hasArg2 := in.args["arg2"]
	listName, arg1, arg2 := in.Arg("name"), in.Arg("arg1"), in.Arg("arg2")

	if action != "add-source" {
		return 1, configErr(`Unknown action "` + action + `". Use add-source.`)
	}
	if !hasName {
		return 1, configErr("You must pass a dependency policy name. Example: composer policy add-source my-policy url https://example.org")
	}
	if err := assertCustomListName(listName); err != nil {
		return 1, err
	}
	if !hasArg1 {
		return 1, configErr("You must pass the source type and a url, or a JSON string.")
	}
	var source *phpjson.Map
	if reJSONObjectStart.MatchString(arg1) {
		v, err := jsonfile.Parse([]byte(arg1), "")
		if err != nil {
			return 1, err
		}
		if !phpjson.IsArray(v) {
			return 1, configErr("Source JSON must be an object.")
		}
		source = phpjson.ToMap(v)
	} else {
		if !hasArg2 {
			return 1, configErr("You must pass the source type and a url. Example: composer policy add-source my-policy url https://example.org")
		}
		source = phpjson.NewMap()
		source.Set("type", arg1)
		source.Set("url", arg2)
	}
	if err := validatePolicySource(listName, source); err != nil {
		return 1, err
	}

	data, err := readJSONMap(b.configFile)
	if err != nil {
		return 1, err
	}
	var current phpjson.Value = []phpjson.Value{}
	if cfgMap := phpjson.ToMap(mustMapGet(data, "config")); cfgMap != nil {
		if pol := phpjson.ToMap(mustMapGet(cfgMap, "policy")); pol != nil {
			if list := phpjson.ToMap(mustMapGet(pol, listName)); list != nil {
				if s, ok := list.Get("sources"); ok && phpjson.IsArray(s) {
					current = s
				}
			}
		}
	}
	srcType, _ := source.Get("type")
	srcURL, _ := source.Get("url")
	found := false
	phpjson.Each(current, func(_ string, e phpjson.Value) {
		m := phpjson.ToMap(e)
		if found || m == nil || !phpjson.IsArray(e) {
			return
		}
		t, _ := m.Get("type")
		u, _ := m.Get("url")
		if phpStrictEqual(t, srcType) && phpStrictEqual(u, srcURL) {
			found = true
		}
	})
	if found {
		u, _ := phpjson.ScalarString(srcURL)
		ctx.IO.Write("<info>Source " + u + " already present in policy " + listName + "</info>")
		return 0, nil
	}
	var sources phpjson.Value
	switch t := current.(type) {
	case []phpjson.Value:
		sources = append(append([]phpjson.Value(nil), t...), source)
	case *phpjson.Map:
		m := t.Clone()
		next := int64(0)
		for _, k := range m.Keys() {
			if n, ok := phpjson.IntKey(k); ok && n >= next {
				next = n + 1
			}
		}
		m.Set(strconv.FormatInt(next, 10), source)
		sources = m
	}
	return srcResult(b.source.AddConfigSetting("policy."+listName+".sources", sources))
}

// assertCustomListName is PolicyCommand::assertCustomListName.
func assertCustomListName(name string) error {
	for _, n := range policy.BuiltinListNames {
		if n == name {
			return configErr(`Built-in dependency policy "` + name + `" does not support sources. Use ` + "`composer config policy." + name + ".<field>`" + ` to configure it.`)
		}
	}
	if e := policy.GetFutureReservedListNameError(name); e != "" {
		return configErr(e)
	}
	if name == "" || strings.Contains(name, ".") {
		return configErr(`Invalid dependency policy name "` + name + `".`)
	}
	return nil
}

// validatePolicySource is SourceValidator::validate.
func validatePolicySource(listName string, source *phpjson.Map) error {
	t, ok := source.Get("type")
	if !ok || t == nil {
		return configErr(`Source configuration is missing the "type" field.`)
	}
	if t != "url" {
		ts, _ := phpjson.ScalarString(t)
		if phpjson.IsArray(t) {
			ts = "Array"
		}
		return configErr(`Unsupported source type "` + ts + `". Only "url" is currently supported.`)
	}
	u, ok := source.Get("url")
	us, isStr := u.(string)
	if !ok || !isStr {
		return configErr(`Source configuration is missing a string "url" field.`)
	}
	if !strings.HasPrefix(us, "https://") {
		return configErr(`Source URL for policy list "` + listName + `" must start with "https://"; got "` + us + `".`)
	}
	return nil
}

// phpStrictEqual is === for JSON scalars (null when missing).
func phpStrictEqual(a, b phpjson.Value) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case string, bool, int64, float64:
		return a == b
	default:
		_ = x
		return false
	}
}
