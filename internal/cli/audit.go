package cli

// Port of Composer\Command\AuditCommand (plus BaseCommand::getAuditFormat)
// (composer/composer f83588392df43d4fd027d84d26de459946ef4e80,
// src/Composer/Command/AuditCommand.php, MIT,
// (c) Nils Adermann, Jordi Boggiano).
// The audit itself is internal/policy's Auditor, shared with the --audit
// run of install/update.

import (
	"strings"

	"gocomposer/internal/composer"
	"gocomposer/internal/pkg"
	"gocomposer/internal/policy"
)

func init() { registerCommand(auditCommand) }

func auditCommand() *Command {
	return &Command{Name: "audit", Run: runAudit}
}

func runAudit(ctx *Context) (int, error) {
	io, in := ctx.IO, ctx.In
	c, err := ctx.composerFor(true)
	if err != nil {
		return 1, err
	}
	packages, err := auditPackages(c, in)
	if err != nil {
		return 1, err
	}
	root := c.Package
	if len(packages) == 0 {
		if len(root.Requires()) > 0 || (!in.Bool("no-dev") && len(root.DevRequires()) > 0) {
			io.WriteError(`No installed packages found. Please run "composer install" before running "audit" or pass "--locked" to audit the lock file.`)
			return policy.StatusFailed, nil
		}
		io.WriteError("No packages - skipping audit.")
		return policy.StatusOK, nil
	}

	// filtered after the empty check so an empty result is not mistaken for a missing install
	if in.Bool("direct") {
		packages = filterDirectPackages(root, in.Bool("no-dev"), packages)
		if len(packages) == 0 {
			io.WriteError("No direct dependencies - skipping audit.")
			return policy.StatusOK, nil
		}
	}

	auditModes := []string{policy.AuditIgnore, policy.AuditReport, policy.AuditFail}
	abandoned, hasAbandoned := in.String("abandoned")
	if hasAbandoned && !inList(abandoned, auditModes) {
		return 1, &composer.UserError{Msg: "--abandoned must be one of " + strings.Join(auditModes, ", ") + ".", Code: 1}
	}

	cfg, err := composer.CreatePolicyConfig(c, false, false)
	if err != nil {
		return 1, err
	}
	if hasAbandoned {
		cfg = cfg.WithAudit(abandoned)
	}
	if sev := in.Strings("ignore-severity"); len(sev) > 0 {
		cfg = cfg.WithIgnoreSeverity(append([]string(nil), sev...))
	}
	if in.Bool("ignore-unreachable") {
		if cfg, err = cfg.WithIgnoreUnreachable("audit"); err != nil {
			return 1, err
		}
	}

	resolver := c.PolicyResolver
	if resolver == nil {
		resolver = &policy.ProviderResolver{IO: c.IO}
		if c.HTTP != nil {
			resolver.HTTP = c.HTTP
		}
		c.PolicyResolver = resolver
	}
	var filterSet policy.FilterListSource
	if cfg.Enabled {
		fs, err := policy.CreateFilterListProviderSet(cfg, c.Repositories, resolver.HTTP, resolver)
		if err != nil {
			return 1, err
		}
		filterSet = fs
	}

	// BaseCommand::getAuditFormat
	format, _ := in.String("format")
	if !inList(format, policy.Formats) {
		return 1, &composer.UserError{Msg: "--format must be one of " + strings.Join(policy.Formats, ", ") + ".", Code: 1}
	}

	src := &policy.RepositoryAdvisorySource{Repositories: c.Repositories, Resolver: resolver}
	status, err := policy.NewAuditor(nil).Audit(io, src, cfg, packages, format, false, filterSet)
	if err != nil {
		return 1, err
	}
	if status > 255 {
		status = 255
	}
	return status, nil
}

// auditPackages is AuditCommand::getPackages.
func auditPackages(c *composer.Composer, in *Input) ([]*pkg.Package, error) {
	if in.Bool("locked") {
		if locked, _ := c.Locker.IsLocked(); !locked {
			return nil, &composer.UserError{Msg: "Valid composer.json and composer.lock files are required to run this command with --locked", Code: 1}
		}
		lr, err := c.Locker.LockedRepository(!in.Bool("no-dev"))
		if err != nil {
			return nil, err
		}
		return lr.MustPackages(), nil
	}
	installed := c.LocalRepo.MustPackages()
	if in.Bool("no-dev") {
		return filterRequiredPackages(installed, c.Package, false, nil), nil
	}
	return installed, nil
}

// filterDirectPackages is AuditCommand::filterDirectPackages.
func filterDirectPackages(root *pkg.Package, noDev bool, packages []*pkg.Package) []*pkg.Package {
	requires := root.Requires()
	if !noDev {
		requires = append(append(pkg.Links(nil), requires...), root.DevRequires()...)
	}
	var out []*pkg.Package
	for _, p := range packages {
		for _, n := range p.Names(true) {
			if requires.Has(n) {
				out = append(out, p)
				break
			}
		}
	}
	return out
}
