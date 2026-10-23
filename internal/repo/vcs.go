package repo

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/VcsRepository.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"errors"
	"regexp"
	"strings"

	"gocomposer/internal/config"
	"gocomposer/internal/httpdl"
	"gocomposer/internal/phpjson"
	"gocomposer/internal/pkg"
	"gocomposer/internal/semver"
	"gocomposer/internal/vcs"
)

// Vcs is VcsRepository.
type Vcs struct {
	*Array
	url                 string
	typ                 string
	env                 *vcs.Env
	repoConfig          *phpjson.Map
	driver              vcs.Driver
	driverType          string
	driverErr           error
	packageName         string
	onlyBranches        *regexp.Regexp
	branchErrorOccurred bool
	emptyReferences     []string
}

// vcsDriverOrder is VcsRepository::$drivers (svn must stay last).
var vcsDriverOrder = []string{"github", "gitlab", "bitbucket", "git-bitbucket", "forgejo", "git", "hg", "perforce", "fossil", "svn"}

// NewVcs is VcsRepository::__construct.
func NewVcs(repoConfig *phpjson.Map, env *vcs.Env) (*Vcs, error) {
	repoConfig = repoConfig.Clone()
	u := config.ExpandPath(repoConfig.Str("url"))
	repoConfig.Set("url", u)
	typ := "vcs"
	if repoConfig.Isset("type") {
		typ = repoConfig.Str("type")
	}
	r := &Vcs{Array: &Array{kind: KindVcs}, url: u, typ: typ, env: env, repoConfig: repoConfig}
	r.Array.self = r
	r.Array.init = r.initialize
	if ob, ok := repoConfig.Get("only-branches"); ok && ob != nil {
		var names []string
		valid := phpjson.IsArray(ob)
		phpjson.Each(ob, func(_ string, v phpjson.Value) {
			s, isStr := v.(string)
			if !isStr {
				valid = false
				return
			}
			names = append(names, s)
		})
		if !valid {
			return nil, &vcs.InvalidArgumentError{Msg: `"only-branches" key for repository ` + vcs.Sanitize(u) + " should be an array of strings"}
		}
		// branch names are case-sensitive, unlike package names
		parts := make([]string, len(names))
		for i, n := range names {
			parts[i] = strings.ReplaceAll(regexp.QuoteMeta(n), `\*`, `.*`)
		}
		r.onlyBranches = regexp.MustCompile(`^(?:` + strings.Join(parts, "|") + `)$`)
	}
	return r, nil
}

// RepoName is VcsRepository::getRepoName.
func (r *Vcs) RepoName() string {
	if _, err := r.Driver(); err != nil {
		return "vcs repo (" + r.typ + " " + vcs.Sanitize(r.url) + ")"
	}
	return "vcs repo (" + r.driverType + " " + vcs.Sanitize(r.url) + ")"
}

// RepoConfig is getRepoConfig.
func (r *Vcs) RepoConfig() *phpjson.Map { return r.repoConfig }

// HadInvalidBranches is VcsRepository::hadInvalidBranches.
func (r *Vcs) HadInvalidBranches() bool { return r.branchErrorOccurred }

// EmptyReferences is VcsRepository::getEmptyReferences.
func (r *Vcs) EmptyReferences() []string { return r.emptyReferences }

// Driver is VcsRepository::getDriver (nil, nil when no driver handles the url).
func (r *Vcs) Driver() (vcs.Driver, error) {
	if r.driver != nil || r.driverErr != nil {
		return r.driver, r.driverErr
	}
	init := func(typ string) (vcs.Driver, error) {
		d := vcs.NewDriver(typ, r.repoConfig, r.env)
		r.driver, r.driverType = d, typ
		if err := d.Initialize(); err != nil {
			r.driverErr = err
			return d, err
		}
		return d, nil
	}
	if vcs.IsDriverType(r.typ) {
		return init(r.typ)
	}
	for _, deep := range []bool{false, true} {
		for _, typ := range vcsDriverOrder {
			if typ == "git-bitbucket" {
				// same class as "bitbucket", already checked
				continue
			}
			ok, err := vcs.Supports(typ, r.env, r.url, deep)
			if err != nil {
				r.driverErr = err
				return nil, err
			}
			if ok {
				return init(typ)
			}
		}
	}
	return nil, nil
}

// overwriteError is $this->io->overwriteError($msg, false) (the erased
// length is the IO's last stderr message, like ConsoleIO).
func (r *Vcs) overwriteError(msg string) { r.env.IO.OverwriteError(msg, false) }

func shouldRethrowTransport(err error) bool {
	var te *httpdl.TransportError
	if !errors.As(err, &te) {
		return false
	}
	c := te.StatusCode
	return c == 401 || c == 403 || c == 429 || c >= 500
}

func isTransport(err error) (*httpdl.TransportError, bool) {
	var te *httpdl.TransportError
	ok := errors.As(err, &te)
	return te, ok
}

var (
	reTagDevSuffix    = regexp.MustCompile(`(?i)[.-]?dev$`)
	reTagDevAffix     = regexp.MustCompile(`(?i)(^dev-|[.-]?dev$)`)
	reBranchNines     = regexp.MustCompile(`(\.9{7})+`)
	reInvalidTagIsDev = regexp.MustCompile(`(?i)(^dev-|[.-]?dev$)`)
)

// initialize is VcsRepository::initialize.
func (r *Vcs) initialize() error {
	io := r.env.IO
	isVerbose, isVeryVerbose := io.IsVerbose(), io.IsVeryVerbose()
	driver, err := r.Driver()
	if err != nil {
		return err
	}
	if driver == nil {
		return &vcs.InvalidArgumentError{Msg: "No driver found to handle VCS repository " + vcs.Sanitize(r.url)}
	}
	loader := &pkg.Loader{}
	hasRootComposer := false
	func() {
		rootID, err := driver.RootIdentifier()
		var has bool
		if err == nil {
			has, err = driver.HasComposerFile(rootID)
		}
		if err == nil {
			hasRootComposer = has
			if has {
				var data *phpjson.Map
				data, err = driver.ComposerInformation(rootID)
				if err == nil && data != nil {
					if n, ok := data.Get("name"); ok && phpjson.Truthy(n) {
						r.packageName, _ = phpjson.ScalarString(n)
					}
				}
			}
		}
		if err != nil {
			if shouldRethrowTransport(err) {
				r.driverErr = err
				return
			}
			if isVeryVerbose {
				// the message calls getRootIdentifier() again, outside the try
				again, aerr := driver.RootIdentifier()
				if aerr != nil {
					r.driverErr = aerr
					return
				}
				io.WriteError("<error>Skipped parsing " + again + ", " + err.Error() + "</error>")
			}
		}
	}()
	if r.driverErr != nil {
		return r.driverErr
	}
	displayName := func() string {
		if r.packageName != "" {
			return r.packageName
		}
		return vcs.Sanitize(r.url)
	}
	tags, err := driver.Tags()
	if err != nil {
		return err
	}
	for _, t := range tags {
		tag, identifier := t.Name, t.ID
		msg := "Reading composer.json of <info>" + displayName() + "</info> (<comment>" + tag + "</comment>)"
		// strip the release- prefix from tags if present
		tag = strings.ReplaceAll(tag, "release-", "")
		parsedTag, perr := semver.Normalize(tag)
		if perr != nil || parsedTag == "" {
			if isVeryVerbose {
				io.WriteError("<warning>Skipped tag " + tag + ", invalid tag name</warning>")
			}
			continue
		}
		if isVeryVerbose {
			io.WriteError(msg)
		} else if isVerbose {
			r.overwriteError(msg)
		}
		err := func() error {
			data, err := driver.ComposerInformation(identifier)
			if err != nil {
				return err
			}
			if data == nil {
				if isVeryVerbose {
					io.WriteError("<warning>Skipped tag " + tag + ", no composer file</warning>")
				}
				r.emptyReferences = append(r.emptyReferences, identifier)
				return nil
			}
			data = data.Clone()
			if v, ok := data.Get("version"); ok && v != nil {
				vs, _ := phpjson.ScalarString(v)
				n, err := semver.Normalize(vs)
				if err != nil {
					return err
				}
				data.Set("version_normalized", n)
			} else {
				data.Set("version", tag)
				data.Set("version_normalized", parsedTag)
			}
			vs, _ := phpjson.ScalarString(mustGetVal(data, "version"))
			data.Set("version", reTagDevSuffix.ReplaceAllString(vs, ""))
			vn := reTagDevAffix.ReplaceAllString(data.Str("version_normalized"), "")
			data.Set("version_normalized", vn)
			data.Delete("default-branch")
			if vn != parsedTag {
				if isVeryVerbose {
					if reInvalidTagIsDev.MatchString(parsedTag) {
						io.WriteError("<warning>Skipped tag " + tag + ", invalid tag name, tags can not use dev prefixes or suffixes</warning>")
					} else {
						io.WriteError("<warning>Skipped tag " + tag + ", tag (" + parsedTag + ") does not match version (" + vn + ") in composer.json</warning>")
					}
				}
				return nil
			}
			tagPackageName := r.packageName
			if tagPackageName == "" {
				tagPackageName = data.Str("name")
			}
			if existing := FindInList(r.Array.packages, tagPackageName, semver.NewConstraint("==", vn)); len(existing) > 0 {
				if isVeryVerbose {
					io.WriteError("<warning>Skipped tag " + tag + ", it conflicts with an another tag (" + existing[0].PrettyVersion() + ") as both resolve to " + vn + " internally</warning>")
				}
				return nil
			}
			if isVeryVerbose {
				io.WriteError("Importing tag " + tag + " (" + vn + ")")
			}
			pd, err := r.preProcess(driver, data, identifier)
			if err != nil {
				return err
			}
			p, err := loader.Load(pd)
			if err != nil {
				return err
			}
			r.Array.AddPackage(p)
			return nil
		}()
		if err != nil {
			var se *SecurityError
			if errors.As(err, &se) {
				return err
			}
			if te, ok := isTransport(err); ok {
				if te.StatusCode == 404 {
					r.emptyReferences = append(r.emptyReferences, identifier)
				}
				if shouldRethrowTransport(err) {
					return err
				}
			}
			if isVeryVerbose {
				reason := err.Error()
				if te, ok := isTransport(err); ok {
					reason = "no composer file was found (" + itoa(te.StatusCode) + " HTTP status code)"
				}
				io.WriteError("<warning>Skipped tag " + tag + ", " + reason + "</warning>")
			}
			continue
		}
	}
	if !isVeryVerbose {
		r.overwriteError("")
	}
	branches, err := driver.Branches()
	if err != nil {
		return err
	}
	rootID, err := driver.RootIdentifier()
	if err != nil {
		return err
	}
	// make sure the root identifier branch gets loaded first
	if hasRootComposer {
		for i, b := range branches {
			if b.Name == rootID {
				reordered := append([]vcs.Ref{b}, branches[:i]...)
				branches = append(reordered, branches[i+1:]...)
				break
			}
		}
	}
	branches = r.filterBranches(branches, rootID)
	for _, b := range branches {
		branch, identifier := b.Name, b.ID
		msg := "Reading composer.json of <info>" + displayName() + "</info> (<comment>" + branch + "</comment>)"
		if isVeryVerbose {
			io.WriteError(msg)
		} else if isVerbose {
			r.overwriteError(msg)
		}
		parsedBranch, ok := validateBranch(branch)
		if !ok {
			if isVeryVerbose {
				io.WriteError("<warning>Skipped branch " + branch + ", invalid name</warning>")
			}
			continue
		}
		var version string
		if strings.HasPrefix(parsedBranch, "dev-") || parsedBranch == semver.DefaultBranchAlias {
			version = "dev-" + strings.ReplaceAll(branch, "#", "+")
			parsedBranch = strings.ReplaceAll(parsedBranch, "#", "+")
		} else {
			prefix := ""
			if strings.HasPrefix(branch, "v") {
				prefix = "v"
			}
			version = prefix + reBranchNines.ReplaceAllString(parsedBranch, ".x")
		}
		err := func() error {
			data, err := driver.ComposerInformation(identifier)
			if err != nil {
				return err
			}
			if data == nil {
				if isVeryVerbose {
					io.WriteError("<warning>Skipped branch " + branch + ", no composer file</warning>")
				}
				r.emptyReferences = append(r.emptyReferences, identifier)
				return nil
			}
			data = data.Clone()
			data.Set("version", version)
			data.Set("version_normalized", parsedBranch)
			data.Delete("default-branch")
			if rootID == branch {
				data.Set("default-branch", true)
			}
			if isVeryVerbose {
				io.WriteError("Importing branch " + branch + " (" + version + ")")
			}
			pd, err := r.preProcess(driver, data, identifier)
			if err != nil {
				return err
			}
			p, err := loader.Load(pd)
			if err != nil {
				return err
			}
			r.Array.AddPackage(p)
			return nil
		}()
		if err != nil {
			if te, ok := isTransport(err); ok {
				if te.StatusCode == 404 {
					r.emptyReferences = append(r.emptyReferences, identifier)
				}
				if shouldRethrowTransport(err) {
					return err
				}
				if isVeryVerbose {
					io.WriteError("<warning>Skipped branch " + branch + ", no composer file was found (" + itoa(te.StatusCode) + " HTTP status code)</warning>")
				}
				continue
			}
			var se *SecurityError
			if errors.As(err, &se) {
				return err
			}
			if !isVeryVerbose {
				io.WriteError("")
			}
			r.branchErrorOccurred = true
			io.WriteError("<error>Skipped branch " + branch + ", " + err.Error() + "</error>")
			io.WriteError("")
			continue
		}
	}
	driver.Cleanup()
	if !isVeryVerbose {
		r.overwriteError("")
	}
	if len(r.Array.packages) == 0 {
		return &InvalidRepositoryError{"No valid composer.json was found in any branch or tag of " + vcs.Sanitize(r.url) + ", could not load a package from it."}
	}
	return nil
}

// SecurityError is Composer\Exception\SecurityException.
type SecurityError struct{ Msg string }

func (e *SecurityError) Error() string { return e.Msg }

// preProcess is VcsRepository::preProcess.
func (r *Vcs) preProcess(driver vcs.Driver, data *phpjson.Map, identifier string) (*phpjson.Map, error) {
	var dataName phpjson.Value
	if v, ok := data.Get("name"); ok {
		dataName = v
	}
	if r.packageName != "" {
		data.Set("name", r.packageName)
	} else {
		data.Set("name", dataName)
	}
	data.Delete("published-time")
	if v, ok := data.Get("dist"); !ok || v == nil {
		dist, err := driver.Dist(identifier)
		if err != nil {
			return nil, err
		}
		if dist != nil {
			data.Set("dist", dist)
		} else {
			data.Set("dist", nil)
		}
	}
	if v, ok := data.Get("source"); !ok || v == nil {
		data.Set("source", driver.Source(identifier))
	}
	if dist := data.Map("dist"); dist != nil && !dist.Isset("reference") {
		if src := data.Map("source"); src != nil && src.Isset("reference") {
			ref, _ := src.Get("reference")
			dist = dist.Clone()
			dist.Set("reference", ref)
			data.Set("dist", dist)
		}
	}
	return data, nil
}

// filterBranches is VcsRepository::filterBranches.
func (r *Vcs) filterBranches(branches []vcs.Ref, rootID string) []vcs.Ref {
	if r.onlyBranches == nil {
		return branches
	}
	var out []vcs.Ref
	for _, b := range branches {
		if b.Name == rootID || b.ID == rootID || r.onlyBranches.MatchString(b.Name) {
			out = append(out, b)
			continue
		}
		if r.env.IO.IsVeryVerbose() {
			r.env.IO.WriteError("<warning>Skipped branch " + b.Name + ", not matching the only-branches filter</warning>")
		}
	}
	return out
}

// validateBranch is VcsRepository::validateBranch.
func validateBranch(branch string) (string, bool) {
	n := semver.NormalizeBranch(branch)
	if _, err := semver.ParseConstraints(n); err != nil {
		return "", false
	}
	return n, true
}
