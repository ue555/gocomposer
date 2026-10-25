// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80:
// src/Composer/Util/Http/ResponseWarnings.php.
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package httpdl

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"gocomposer/internal/cio"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/semver"
)

// composerVersion is Composer::getVersion() of the reference source checkout
// (the same value as platform.ComposerVersion, which httpdl cannot import).
const composerVersion = "2.11.999-dev+source"

var (
	shownWarningsMu sync.Mutex
	// shownWarnings is ResponseWarnings::$shownWarnings (keyed by host and
	// the warning/info keys of the body).
	shownWarnings = map[string]bool{}
	reANSIColor   = regexp.MustCompile("\x1b\\[[;\\d]*m")
)

// ResetResponseWarnings is ResponseWarnings::reset.
func ResetResponseWarnings() {
	shownWarningsMu.Lock()
	shownWarnings = map[string]bool{}
	shownWarningsMu.Unlock()
}

// OutputResponseWarnings is ResponseWarnings::output: it shows the
// warning(s)/info(s) a repository sent in a decoded response body, unless the
// same host already showed the same ones. It reports whether any was shown,
// now or before.
func OutputResponseWarnings(writeError func(string), decorated bool, rawURL string, data phpjson.Value) bool {
	m := phpjson.AsMap(data)
	if m == nil || !phpjson.IsArray(data) {
		return false
	}
	host := rawURL
	if u, err := url.Parse(rawURL); err == nil && u.Host != "" {
		host = u.Hostname()
		if p := u.Port(); p != "" {
			if _, err := strconv.Atoi(p); err == nil {
				host += ":" + p
			}
		}
	} else if h, p, ok := strings.Cut(rawURL, ":"); ok && !strings.Contains(rawURL, "/") {
		// parse_url("host:port") yields the host and the port
		if _, err := strconv.Atoi(p); err == nil && h != "" {
			host = h + ":" + p
		}
	}
	if host == "repo.packagist.org" {
		host = "packagist.org"
	}
	shown := phpjson.NewMap()
	for _, k := range m.Keys() {
		switch k {
		case "warning", "warning-versions", "info", "info-versions", "warnings", "infos":
			v, _ := m.Get(k)
			shown.Set(k, v)
		}
	}
	enc, _ := phpjson.Encode(shown, 0)
	key := host + ":" + enc
	shownWarningsMu.Lock()
	already := shownWarnings[key]
	shownWarningsMu.Unlock()
	if already {
		return true
	}
	clean := func(v phpjson.Value) string {
		s, _ := phpjson.ScalarString(v)
		if !decorated {
			s = reANSIColor.ReplaceAllString(s, "")
		}
		return s
	}
	composerMatches := func(versions phpjson.Value) bool {
		vs, _ := phpjson.ScalarString(versions)
		c, err := semver.ParseConstraints(vs)
		if err != nil {
			return false
		}
		n, err := semver.Normalize(composerVersion)
		if err != nil {
			return false
		}
		return c.Matches(semver.NewConstraint("==", n))
	}
	wrote := false
	for _, typ := range []string{"warning", "info"} {
		v, _ := m.Get(typ)
		if !phpjson.Truthy(v) {
			continue
		}
		if vv, _ := m.Get(typ + "-versions"); phpjson.Truthy(vv) && !composerMatches(vv) {
			continue
		}
		writeError("<" + typ + ">" + strings.ToUpper(typ[:1]) + typ[1:] + " from " + Sanitize(rawURL) + ": " + clean(v) + "</" + typ + ">")
		wrote = true
	}
	for _, k := range []string{"warnings", "infos"} {
		list, _ := m.Get(k)
		if !phpjson.Truthy(list) {
			continue
		}
		typ := k[:len(k)-1]
		phpjson.Each(list, func(_ string, spec phpjson.Value) {
			sm := phpjson.AsMap(spec)
			if sm == nil {
				return
			}
			versions, _ := sm.Get("versions")
			if !composerMatches(versions) {
				return
			}
			msg, _ := sm.Get("message")
			writeError("<" + typ + ">" + strings.ToUpper(typ[:1]) + typ[1:] + " from " + Sanitize(rawURL) + ": " + clean(msg) + "</" + typ + ">")
			wrote = true
		})
	}
	if wrote {
		shownWarningsMu.Lock()
		shownWarnings[key] = true
		shownWarningsMu.Unlock()
	}
	return wrote
}

// OutputResponseWarnings shows the warnings/infos of a decoded response body
// through the downloader's IO (ResponseWarnings::output($io, $url, $data)).
func (d *Downloader) OutputResponseWarnings(rawURL string, data phpjson.Value) bool {
	write := func(s string) {}
	decorated := false
	switch {
	case d.IO != nil:
		// ConsoleIO::writeError sanitizes the message (escape sequences
		// and control characters are removed even on a decorated output)
		write = func(s string) { d.IO.WriteError(cio.Sanitize(s)) }
		decorated = d.IO.Decorated
	case d.Warn != nil:
		write = d.Warn
	}
	return OutputResponseWarnings(write, decorated, rawURL, data)
}

// isJSONContentType is CurlDownloader::isJsonResponse.
func isJSONContentType(contentType string) bool {
	ct := strings.ToLower(contentType)
	return ct == "application/json" || ct == "application/json; charset=utf-8"
}
