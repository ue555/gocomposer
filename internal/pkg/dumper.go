package pkg

import (
	"sort"

	"gocomposer/internal/phpjson"
)

// Dump ports Composer\Package\Dumper\ArrayDumper::dump.
func Dump(p *Package) *phpjson.Map {
	d := phpjson.NewMap()
	d.Set("name", p.PrettyName)
	d.Set("version", p.PrettyVersion())
	d.Set("version_normalized", p.Version())
	if td := p.TargetDir(); td != "" {
		d.Set("target-dir", td)
	}
	if st := p.SourceType(); st != "" {
		src := phpjson.NewMap()
		src.Set("type", st)
		src.Set("url", p.SourceURL())
		if r := p.SourceReference(); r != "" {
			src.Set("reference", r)
		}
		if m := p.SourceMirrors(); len(m) > 0 {
			src.Set("mirrors", m)
		}
		d.Set("source", src)
	}
	if dt := p.DistType(); dt != "" {
		dist := phpjson.NewMap()
		dist.Set("type", dt)
		dist.Set("url", p.DistURL())
		if r := p.DistReference(); r != "" {
			dist.Set("reference", r)
		}
		if s := p.DistSha1Checksum(); s != "" || p.Base().distSha1Present {
			dist.Set("shasum", s)
		}
		if m := p.DistMirrors(); len(m) > 0 {
			dist.Set("mirrors", m)
		}
		d.Set("dist", dist)
	}
	for _, kind := range LinkKinds {
		links := p.LinksOf(kind.Type)
		if len(links) == 0 {
			continue
		}
		lm := phpjson.NewMap()
		for _, l := range links {
			lm.Set(l.Target, l.PrettyConstraint)
		}
		lm.KSort()
		d.Set(kind.Key, lm)
	}
	if s := p.Suggests(); s.Len() > 0 {
		c := s.Clone()
		c.KSort()
		d.Set("suggest", c)
	}
	if t := ReleaseTime(p.ReleaseDate()); t != "" {
		d.Set("time", t)
	}
	if t := ReleaseTime(p.Base().publishedDate); t != "" {
		d.Set("published-time", t)
	}
	if p.IsDefaultBranch() {
		d.Set("default-branch", true)
	}
	if b := p.Binaries(); len(b) > 0 {
		d.Set("bin", b)
	}
	d.Set("type", p.Type())
	if e := p.Extra(); e.Len() > 0 {
		d.Set("extra", e)
	}
	if s := p.InstallationSource(); s != "" {
		d.Set("installation-source", s)
	}
	if a := p.Autoload(); a.Len() > 0 {
		d.Set("autoload", a)
	}
	if a := p.DevAutoload(); a.Len() > 0 {
		d.Set("autoload-dev", a)
	}
	if n := p.NotificationURL(); n != "" {
		d.Set("notification-url", n)
	}
	if ip := p.IncludePaths(); len(ip) > 0 {
		d.Set("include-path", ip)
	}
	if pe := p.Base().phpExt; pe != nil && !(phpjson.IsArray(pe) && phpjson.ArrayLen(pe) == 0) {
		d.Set("php-ext", pe)
	}
	b := p.Base()
	if b.archiveName != "" || len(b.archiveExcludes) > 0 {
		ar := phpjson.NewMap()
		if b.archiveName != "" {
			ar.Set("name", b.archiveName)
		}
		if len(b.archiveExcludes) > 0 {
			ar.Set("exclude", b.archiveExcludes)
		}
		d.Set("archive", ar)
	}
	if s := p.Scripts(); s.Len() > 0 {
		d.Set("scripts", s)
	}
	if l := p.License(); len(l) > 0 {
		d.Set("license", l)
	}
	if a := p.Authors(); len(a) > 0 {
		d.Set("authors", a)
	}
	if p.Description() != "" {
		d.Set("description", p.Description())
	}
	if p.Homepage() != "" {
		d.Set("homepage", p.Homepage())
	}
	if k := p.Keywords(); len(k) > 0 {
		ks := append([]string(nil), k...)
		sort.SliceStable(ks, func(i, j int) bool { return phpjson.PHPCompareKeys(ks[i], ks[j]) < 0 })
		d.Set("keywords", ks)
	}
	if r := b.repositories; len(r) > 0 {
		d.Set("repositories", r)
	}
	if s := p.Support(); s.Len() > 0 {
		d.Set("support", s)
	}
	if f := p.Funding(); len(f) > 0 {
		d.Set("funding", f)
	}
	if p.IsAbandoned() {
		if r := p.ReplacementPackage(); r != "" {
			d.Set("abandoned", r)
		} else {
			d.Set("abandoned", true)
		}
	}
	if p.IsRoot && p.MinimumStability != "" {
		d.Set("minimum-stability", p.MinimumStability)
	}
	if to := p.TransportOptions(); to.Len() > 0 {
		d.Set("transport-options", to)
	}
	return d
}
