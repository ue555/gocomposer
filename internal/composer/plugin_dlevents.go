package composer

// The plugin events of downloads and pool creation:
// PluginEvents::PRE_FILE_DOWNLOAD / POST_FILE_DOWNLOAD (FileDownloader::
// download for packages, ComposerRepository::fetchFile/
// fetchFileIfLastModified/asyncFetchFile for metadata) and
// PluginEvents::PRE_POOL_CREATE (PoolBuilder::buildPool) of Composer
// f83588392df43d4fd027d84d26de459946ef4e80 (MIT license, (c) Nils Adermann,
// Jordi Boggiano). The event objects live in the plugin runtime; what the
// listeners changed (processed URL, custom cache key, transport options,
// pool packages) is read back and applied by gocomposer.

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"os"
	"sort"
	"strings"
	"sync"

	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/repo"
	"gocomposer/internal/solver"
)

// eventDispatcherFor is the EventDispatcher the Composer instance's
// downloaders and pool builder use (the running Installer's one).
func (pm *PluginManager) eventDispatcherFor() *EventDispatcher {
	if ri := currentInstaller(); ri != nil && ri.in.C == pm.c && ri.in.Events != nil {
		return ri.in.Events
	}
	return pm.dispatcher()
}

// wantsEvent reports whether dispatching name has an observable effect.
func wantsEvent(d *EventDispatcher, name string) bool {
	return os.Getenv("COMPOSER_DEBUG_EVENTS") != "" || d.HasListeners(name)
}

// dispatchStateful dispatches an event whose object listeners may modify and
// returns the object's state afterwards (nil when no plugin listener ran).
func (pm *PluginManager) dispatchStateful(d *EventDispatcher, ev *ScriptEvent, extra *phpjson.Map) (*phpjson.Map, error) {
	pm.eventSeq++
	ev.eid = "s" + itoa(pm.eventSeq)
	if pm.eventExtras == nil {
		pm.eventExtras = map[string]*phpjson.Map{}
	}
	pm.eventExtras[ev.eid] = extra
	defer delete(pm.eventExtras, ev.eid)
	if _, err := d.dispatch(ev); err != nil {
		return nil, err
	}
	if !pm.active() {
		return nil, nil
	}
	args := phpjson.NewMap()
	args.Set("eid", ev.eid)
	v, err := pm.host.call("eventState", args)
	_, _ = pm.host.call("freeEvent", args)
	if err != nil {
		return nil, err
	}
	return phpjson.ToMap(v), nil
}

// hooks of the repositories created by the Composer instance

func (pm *PluginManager) installRepoHooks(c *Composer) {
	for _, r := range c.Repositories {
		pm.repoHook(r)
	}
}

func (pm *PluginManager) repoHook(r repo.Repository) {
	if f, ok := r.(*repo.Filter); ok {
		r = f.Repo
	}
	cr, ok := r.(*repo.Composer)
	if !ok {
		return
	}
	outer := r
	cr.PreFileDownload = func(u string, options *phpjson.Map) (string, *phpjson.Map, error) {
		return pm.metadataPreDownload(outer, u, options)
	}
	cr.PostFileDownload = func(u, sha256 string, resp *httpdl.Response) error {
		return pm.metadataPostDownload(outer, u, sha256, resp)
	}
}

func (pm *PluginManager) metadataPreDownload(r repo.Repository, u string, options *phpjson.Map) (string, *phpjson.Map, error) {
	if pm.closed {
		return u, options, nil
	}
	d := pm.eventDispatcherFor()
	if !wantsEvent(d, "pre-file-download") {
		return u, options, nil
	}
	if options == nil {
		options = phpjson.NewMap()
	}
	extra := phpjson.NewMap()
	extra.Set("url", u)
	extra.Set("type", "metadata")
	extra.Set("options", options)
	if pm.active() {
		extra.Set("repo", pm.repoDesc(pm.c, r, repoIndex(pm.c, r)))
	}
	st, err := pm.dispatchStateful(d, &ScriptEvent{Name: "pre-file-download", Kind: "prefiledl"}, extra)
	if err != nil {
		return "", nil, err
	}
	if st == nil {
		return u, options, nil
	}
	if s := st.Str("url"); s != "" {
		u = s
	}
	if o := phpjson.ToMap(mustGet(st, "options")); o != nil {
		options = o
	} else {
		options = phpjson.NewMap()
	}
	return u, options, nil
}

func (pm *PluginManager) metadataPostDownload(r repo.Repository, u, sha256 string, resp *httpdl.Response) error {
	if pm.closed {
		return nil
	}
	d := pm.eventDispatcherFor()
	if !wantsEvent(d, "post-file-download") {
		return nil
	}
	extra := phpjson.NewMap()
	extra.Set("url", u)
	extra.Set("type", "metadata")
	if sha256 != "" {
		extra.Set("checksum", sha256)
	}
	extra.Set("status", int64(resp.StatusCode))
	extra.Set("headers", responseHeaderLines(resp))
	extra.Set("body", base64.StdEncoding.EncodeToString(resp.Body))
	if pm.active() {
		extra.Set("repo", pm.repoDesc(pm.c, r, repoIndex(pm.c, r)))
	}
	_, err := pm.dispatchStateful(d, &ScriptEvent{Name: "post-file-download", Kind: "postfiledl"}, extra)
	return err
}

func repoIndex(c *Composer, r repo.Repository) int {
	for i, e := range c.Repositories {
		if e == r {
			return i
		}
	}
	return -1
}

func responseHeaderLines(resp *httpdl.Response) []phpjson.Value {
	proto := resp.Proto
	if proto == "" {
		proto = "HTTP/1.1"
	}
	lines := []phpjson.Value{proto + " " + resp.Status}
	var keys []string
	for k := range resp.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range resp.Header[k] {
			lines = append(lines, strings.ToLower(k)+": "+v)
		}
	}
	return lines
}

// package downloads (FileDownloader::download)

// distCacheKey is FileDownloader's cache key generator.
func distCacheKey(p *pkg.Package, key string) string {
	sum := sha1.Sum([]byte(key))
	return p.Name + "/" + hex.EncodeToString(sum[:]) + "." + p.DistType()
}

// pluginPreFileDownload dispatches pre-file-download for a package file and
// returns the URL and cache key to use.
func (im *InstallationManager) pluginPreFileDownload(p *pkg.Package, processed, cacheKey string) (string, string, error) {
	// the deferred post-file-download events follow the download order
	postMu.Lock()
	postSeq++
	if postOrder[im] == nil {
		postOrder[im] = map[*pkg.Package]int{}
	}
	postOrder[im][p] = postSeq
	postMu.Unlock()
	pm := pluginManagerOf(im.c)
	if pm == nil || pm.closed {
		return processed, cacheKey, nil
	}
	d := pm.eventDispatcherFor()
	if !wantsEvent(d, "pre-file-download") {
		return processed, cacheKey, nil
	}
	extra := phpjson.NewMap()
	extra.Set("url", processed)
	extra.Set("type", "package")
	if pm.active() {
		extra.Set("context", pm.host.ref(p))
	}
	st, err := pm.dispatchStateful(d, &ScriptEvent{Name: "pre-file-download", Kind: "prefiledl"}, extra)
	if err != nil || st == nil {
		return processed, cacheKey, err
	}
	u := st.Str("url")
	if key, ok := mustGet(st, "cacheKey").(string); ok {
		cacheKey = distCacheKey(p, key)
	} else if u != processed {
		cacheKey = distCacheKey(p, u)
	}
	return u, cacheKey, nil
}

// retryMu serializes the pre-file-download events of retried transfers
// (they run concurrently; nothing else uses the plugin host meanwhile).
var retryMu sync.Mutex

// pluginPreFileDownloadRetry is the pre-file-download event FileDownloader
// dispatches again when it retries a download.
func (im *InstallationManager) pluginPreFileDownloadRetry(p *pkg.Package, processed, cacheKey string) (string, string, error) {
	retryMu.Lock()
	defer retryMu.Unlock()
	return im.pluginPreFileDownload(p, processed, cacheKey)
}

type pendingPost struct {
	seq                     int
	p                       *pkg.Package
	fileName, checksum, url string
}

var (
	postMu      sync.Mutex
	postSeq     int
	postPending = map[*InstallationManager][]pendingPost{}
	postOrder   = map[*InstallationManager]map[*pkg.Package]int{}
)

// pluginPostFileDownload dispatches post-file-download for a package file
// (deferred: the transfers run concurrently; Composer dispatches the events
// as the downloads complete, they are dispatched in operation order once the
// transfers are done, like with COMPOSER_MAX_PARALLEL_HTTP=1).
func (im *InstallationManager) pluginPostFileDownload(p *pkg.Package, fileName, checksum, u string, deferred bool) error {
	pm := pluginManagerOf(im.c)
	if pm == nil || pm.closed {
		return nil
	}
	if deferred {
		postMu.Lock()
		postPending[im] = append(postPending[im], pendingPost{postOrder[im][p], p, fileName, checksum, u})
		postMu.Unlock()
		return nil
	}
	d := pm.eventDispatcherFor()
	if !wantsEvent(d, "post-file-download") {
		return nil
	}
	extra := phpjson.NewMap()
	extra.Set("url", u)
	extra.Set("type", "package")
	extra.Set("file", fileName)
	if checksum != "" {
		extra.Set("checksum", checksum)
	}
	if pm.active() {
		extra.Set("context", pm.host.ref(p))
	}
	_, err := pm.dispatchStateful(d, &ScriptEvent{Name: "post-file-download", Kind: "postfiledl"}, extra)
	return err
}

// flushPluginPostDownloads dispatches the deferred post-file-download events.
func (im *InstallationManager) flushPluginPostDownloads() error {
	postMu.Lock()
	list := postPending[im]
	delete(postPending, im)
	delete(postOrder, im)
	postMu.Unlock()
	sort.Slice(list, func(i, j int) bool { return list[i].seq < list[j].seq })
	for _, e := range list {
		if err := im.pluginPostFileDownload(e.p, e.fileName, e.checksum, e.url, false); err != nil {
			return err
		}
	}
	return nil
}

// pool creation (PoolBuilder::buildPool)

// prePoolCreateHook returns the PRE_POOL_CREATE dispatcher of c's builders.
func prePoolCreateHook(c *Composer) func(ev *solver.PrePoolCreateEvent) error {
	return func(ev *solver.PrePoolCreateEvent) error {
		pm := pluginManagerOf(c)
		if pm == nil || pm.closed {
			return nil
		}
		d := pm.eventDispatcherFor()
		if !wantsEvent(d, "pre-pool-create") {
			return nil
		}
		extra := phpjson.NewMap()
		if pm.active() {
			h := pm.host
			refs := func(ps []*pkg.Package) []phpjson.Value {
				out := []phpjson.Value{}
				for _, p := range ps {
					out = append(out, h.ref(p))
				}
				return out
			}
			repos := []phpjson.Value{}
			for _, r := range ev.Repositories {
				repos = append(repos, pm.repoDesc(c, r, repoIndex(c, r)))
			}
			extra.Set("repos", repos)
			req := phpjson.NewMap()
			req.Set("fixed", refs(ev.Request.FixedPackages()))
			req.Set("locked", refs(ev.Request.LockedPackages()))
			requires := []phpjson.Value{}
			for _, r := range ev.Request.Requires() {
				m := phpjson.NewMap()
				m.Set("name", r.Name)
				if r.Constraint != nil {
					m.Set("pretty", r.Constraint.PrettyString())
					m.Set("c", r.Constraint.String())
				}
				requires = append(requires, m)
			}
			req.Set("requires", requires)
			if list := ev.Request.UpdateAllowList(); len(list) > 0 {
				allow := []phpjson.Value{}
				for _, n := range list {
					allow = append(allow, n)
				}
				req.Set("allow", allow)
				transitive := int64(solver.UpdateOnlyListed)
				if ev.Request.UpdateAllowTransitiveRootDependencies() {
					transitive = solver.UpdateListedWithTransitiveDeps
				} else if ev.Request.UpdateAllowTransitiveDependencies() {
					transitive = solver.UpdateListedWithTransitiveDepsNoRoot
				}
				req.Set("transitive", transitive)
			}
			extra.Set("request", req)
			intMap := func(in map[string]int) *phpjson.Map {
				m := phpjson.NewMap()
				var keys []string
				for k := range in {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					m.Set(k, int64(in[k]))
				}
				return m
			}
			extra.Set("acceptable", intMap(ev.Acceptable))
			extra.Set("flags", intMap(ev.StabilityFlags))
			aliases := phpjson.NewMap()
			var names []string
			for n := range ev.RootAliases {
				names = append(names, n)
			}
			sort.Strings(names)
			for _, n := range names {
				vm := phpjson.NewMap()
				var vers []string
				for v := range ev.RootAliases[n] {
					vers = append(vers, v)
				}
				sort.Strings(vers)
				for _, v := range vers {
					e := ev.RootAliases[n][v]
					am := phpjson.NewMap()
					am.Set("alias", e.Alias)
					am.Set("alias_normalized", e.AliasNormalized)
					vm.Set(v, am)
				}
				aliases.Set(n, vm)
			}
			extra.Set("rootAliases", aliases)
			rr := phpjson.NewMap()
			names = names[:0]
			for n := range ev.RootReferences {
				names = append(names, n)
			}
			sort.Strings(names)
			for _, n := range names {
				rr.Set(n, ev.RootReferences[n])
			}
			extra.Set("rootRefs", rr)
			extra.Set("packages", refs(ev.Packages))
			extra.Set("unacceptable", refs(ev.Unacceptable))
		}
		st, err := pm.dispatchStateful(d, &ScriptEvent{Name: "pre-pool-create", Kind: "prepool"}, extra)
		if err != nil || st == nil {
			return err
		}
		list := func(v phpjson.Value) ([]*pkg.Package, error) {
			var out []*pkg.Package
			var perr error
			phpjson.Each(v, func(_ string, r phpjson.Value) {
				if perr != nil {
					return
				}
				p, err := pm.host.pkgFromRef(r)
				if err != nil {
					perr = err
					return
				}
				if p != nil {
					out = append(out, p)
				}
			})
			return out, perr
		}
		packages, err := list(mustGet(st, "packages"))
		if err != nil {
			return err
		}
		unacceptable, err := list(mustGet(st, "unacceptable"))
		if err != nil {
			return err
		}
		ev.Packages, ev.Unacceptable = packages, unacceptable
		return nil
	}
}
