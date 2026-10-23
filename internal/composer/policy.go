package composer

// Wiring of the dependency policy layer (internal/policy) into the
// Installer and PoolBuilder, porting Installer::getPolicyConfig,
// create{SecurityAudit,FilterList,Cooldown}PoolFilter,
// BaseCommand::createPolicyConfig and PoolBuilder::run*Filter
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80, MIT,
// (c) Nils Adermann, Jordi Boggiano).

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gocomposer/internal/cio"
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
	"gocomposer/internal/repo"
	"gocomposer/internal/semver"
	"gocomposer/internal/solver"
)

// boolEnv is Platform::getBoolEnv($name, false).
func boolEnv(name string) (bool, error) {
	v, ok := os.LookupEnv(name)
	if !ok || v == "" {
		return false, nil
	}
	switch v {
	case "1", "true", "on":
		return true, nil
	case "0", "false", "off":
		return false, nil
	}
	return false, &UserError{Msg: "Invalid value for " + name + ": " + v + ". Expected 0, 1, false, true, off, or on.", Code: 1}
}

// CreatePolicyConfig is BaseCommand::createPolicyConfig.
func CreatePolicyConfig(c *Composer, noBlockingOpt, noCooldownOpt bool) (*policy.PolicyConfig, error) {
	cfg, err := policy.FromConfig(c.Config)
	if err != nil {
		return nil, err
	}
	noBlocking, err := boolEnv("COMPOSER_NO_BLOCKING")
	if err != nil {
		return nil, err
	}
	if !noBlocking {
		if noBlocking, err = boolEnv("COMPOSER_NO_SECURITY_BLOCKING"); err != nil {
			return nil, err
		}
	}
	if noBlocking || noBlockingOpt {
		return cfg.WithBlockingDisabled(), nil
	}
	if noCooldownOpt {
		return cfg.WithCooldownBlockingDisabled(), nil
	}
	return cfg, nil
}

// getPolicyConfig is Installer::getPolicyConfig.
func (in *Installer) getPolicyConfig() (*policy.PolicyConfig, error) {
	if in.PolicyConfig == nil {
		cfg, err := policy.FromConfig(in.C.Config)
		if err != nil {
			return nil, err
		}
		in.PolicyConfig = cfg
	}
	return in.PolicyConfig, nil
}

func formatCount(n int) string {
	s := strconv.Itoa(n)
	var out []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}
	return string(out)
}

func policyRequest(req *solver.Request) policy.Request {
	var locked []*pkg.Package
	if req.LockedRepository != nil {
		locked = req.LockedRepository.MustPackages()
	}
	return policy.NewRequest(req.IsLockedPackage, locked)
}

// poolFilter builds the PoolBuilder hook running the security advisory,
// filter list and cooldown pool filters like PoolBuilder::buildPool. scope
// is the filter list block scope; update enables the update-only filters.
func (in *Installer) poolFilter(scope string, update bool) (func(*solver.Pool, []repo.Repository, *solver.Request) (*solver.Pool, error), error) {
	cfg, err := in.getPolicyConfig()
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}
	io := in.C.IO
	resolver := in.C.policyResolver()
	var advisory *policy.SecurityAdvisoryPoolFilter
	var filterList *policy.FilterListPoolFilter
	var cooldown *policy.CooldownPoolFilter
	if update && cfg.Advisories.ShouldBlock(policy.BlockScopeUpdate) && !in.UpdateMirrors {
		advisory = policy.NewSecurityAdvisoryPoolFilter(policy.NewAuditor(nil), cfg, io, resolver)
	}
	hasUpdateLists := cfg.GetActiveBlockFilterLists(policy.BlockScopeUpdate).Len() > 0
	hasInstallLists := cfg.GetActiveBlockFilterLists(policy.BlockScopeInstall).Len() > 0
	if (scope == policy.BlockScopeInstall && hasInstallLists) || (scope == policy.BlockScopeUpdate && (hasUpdateLists || hasInstallLists)) {
		filterList = policy.NewFilterListPoolFilter(cfg, resolver.HTTP, scope, in.C.Repositories, io, resolver)
	}
	if update && cfg.Cooldown.ShouldBlock(policy.BlockScopeUpdate) && !in.UpdateMirrors {
		cooldown = policy.NewCooldownPoolFilter(cfg.Cooldown, time.Time{})
	}
	if advisory == nil && filterList == nil && cooldown == nil {
		return nil, nil
	}
	return func(pool *solver.Pool, repos []repo.Repository, req *solver.Request) (*solver.Pool, error) {
		preq := policyRequest(req)
		apply := func(name, summary string, run func(*policy.Pool) (*policy.Pool, error)) error {
			io.WriteErrorV("Running "+name+" pool filter.", true, cio.Debug)
			before := time.Now()
			total := len(pool.Packages())
			res, err := run(&policy.Pool{Packages: pool.Packages(), PolicyRemoved: pool.AllPolicyRemovedVersions()})
			if err != nil {
				return err
			}
			pool = pool.WithPackages(res.Packages, res.PolicyRemoved)
			filtered := total - len(pool.Packages())
			if filtered == 0 {
				return nil
			}
			title := strings.ToUpper(name[:1]) + name[1:]
			io.WriteV(fmt.Sprintf("%s pool filter completed in %.3f seconds", title, time.Since(before).Seconds()), true, cio.VeryVerbose)
			if summary != "" {
				io.WriteV(fmt.Sprintf(summary, formatCount(total), formatCount(filtered), int(float64(100)/float64(total)*float64(filtered)+0.5)), true, cio.VeryVerbose)
			}
			return nil
		}
		if advisory != nil {
			if err := apply("security advisory", "<info>Found %s package versions referenced in your dependency graph. %s (%d%%) were filtered away.</info>", func(p *policy.Pool) (*policy.Pool, error) {
				return advisory.Filter(p, repos, preq)
			}); err != nil {
				return nil, err
			}
		}
		if filterList != nil {
			if err := apply("filter list", "<info>Found %s package versions referenced in your dependency graph. %s (%d%%) were filtered away by dependency policies.</info>", func(p *policy.Pool) (*policy.Pool, error) {
				return filterList.Filter(p, preq)
			}); err != nil {
				return nil, err
			}
		}
		if cooldown != nil {
			io.WriteErrorV("Running cooldown pool filter.", true, cio.Debug)
			before := time.Now()
			total := len(pool.Packages())
			res, err := cooldown.Filter(&policy.Pool{Packages: pool.Packages(), PolicyRemoved: pool.AllPolicyRemovedVersions()}, preq)
			if err != nil {
				return nil, err
			}
			pool = pool.WithPackages(res.Packages, res.PolicyRemoved)
			in.writeCooldownUnverified(cooldown)
			if filtered := total - len(pool.Packages()); filtered > 0 {
				io.WriteV(fmt.Sprintf("Cooldown pool filter completed in %.3f seconds", time.Since(before).Seconds()), true, cio.VeryVerbose)
				suffix := " (run with -vv to list them)."
				if io.IsVeryVerbose() {
					suffix = ":"
				}
				io.WriteError(fmt.Sprintf("<warning>%d package version(s) withheld by the cooldown policy%s</warning>", filtered, suffix))
				if io.IsVeryVerbose() {
					pool.AllPolicyRemovedVersions().All(func(name, _ string, reason *policy.PolicyRemovalReason) {
						info := reason.CooldownInfo()
						if info == nil || reason.PackageName() != name {
							return
						}
						io.WriteError(fmt.Sprintf("  - %s (%s) published %s, available in %s", name, reason.PrettyVersion(), info.ReleaseDate, info.AvailableIn))
					})
				}
			}
		}
		return pool, nil
	}, nil
}

func (in *Installer) writeCooldownUnverified(f *policy.CooldownPoolFilter) {
	io := in.C.IO
	f.UnverifiedVersions().Each(func(repoName string, versions []string) {
		suffix := " (run with -vv to list them)."
		if io.IsVeryVerbose() {
			suffix = ":"
		}
		io.WriteError(fmt.Sprintf("<warning>The publication date of %d package version(s) from %s could not be verified, the cooldown policy relied on the time field set by the package authors%s</warning>", len(versions), repoName, suffix))
		if io.IsVeryVerbose() {
			for _, v := range versions {
				io.WriteError("  - " + v)
			}
		}
	})
}

// runPolicyAudit is the Auditor::audit call of Installer::run.
func (in *Installer) runPolicyAudit(packages []*pkg.Package) (int, error) {
	cfg, err := in.getPolicyConfig()
	if err != nil {
		return 0, err
	}
	resolver := in.C.policyResolver()
	src := &policy.RepositoryAdvisorySource{Repositories: in.C.Repositories, Resolver: resolver}
	var filterSet policy.FilterListSource
	if cfg.Enabled {
		fs, err := policy.CreateFilterListProviderSet(cfg, in.C.Repositories, resolver.HTTP, resolver)
		if err != nil {
			return 0, err
		}
		filterSet = fs
	}
	format := in.AuditFormat
	if format == "" {
		format = "summary"
	}
	return policy.NewAuditor(nil).Audit(in.C.IO, src, cfg, packages, format, true, filterSet)
}

// completeAdvisories is Problem::loadCompleteAdvisories: the pool may only
// hold partial advisories, so the complete ones are loaded for messages.
func (in *Installer) completeAdvisories(rs *solver.RepositorySet) func([]policy.VersionReason) []policy.VersionReason {
	return func(removals []policy.VersionReason) []policy.VersionReason {
		var names []string
		byName := map[string][]semver.Constraint{}
		for _, r := range removals {
			if r.Reason.Type() != policy.RemovalAdvisories {
				continue
			}
			n := r.Reason.PackageName()
			if _, ok := byName[n]; !ok {
				names = append(names, n)
			}
			byName[n] = append(byName[n], semver.NewConstraint("==", r.Version))
		}
		complete := map[string]*policy.Advisory{}
		src := &policy.RepositoryAdvisorySource{Repositories: rs.Repositories(), Resolver: in.C.policyResolver()}
		for _, n := range names {
			packages := rs.FindPackages(n, semver.CreateMulti(byName[n], false), 0)
			res, err := src.GetMatchingSecurityAdvisories(packages, false, true)
			if err != nil || res == nil || res.Advisories == nil {
				continue
			}
			res.Advisories.Each(func(_ string, list []*policy.Advisory) {
				for _, a := range list {
					complete[a.AdvisoryID] = a
				}
			})
		}
		if len(complete) == 0 {
			return removals
		}
		out := make([]policy.VersionReason, len(removals))
		for i, r := range removals {
			out[i] = r
			if r.Reason.Type() != policy.RemovalAdvisories {
				continue
			}
			var advs []*policy.Advisory
			for _, a := range r.Reason.Advisories() {
				if c, ok := complete[a.AdvisoryID]; ok {
					advs = append(advs, c)
				} else {
					advs = append(advs, a)
				}
			}
			out[i].Reason = policy.RemovalReasonAdvisories(r.Reason.PackageName(), r.Reason.PrettyVersion(), advs)
		}
		return out
	}
}
