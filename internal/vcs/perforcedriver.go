package vcs

// Ported from Composer f83588392df43d4fd027d84d26de459946ef4e80
// (src/Composer/Repository/Vcs/PerforceDriver.php).
// Copyright (c) Nils Adermann, Jordi Boggiano. MIT License.

import (
	"time"

	"gocomposer/internal/phpjson"
)

// PerforceDriver is Repository\Vcs\PerforceDriver.
type PerforceDriver struct {
	base
	depot    string
	branch   string
	perforce *Perforce
}

// NewPerforceDriver is new PerforceDriver($repoConfig, ...).
func NewPerforceDriver(repoConfig *phpjson.Map, env *Env) *PerforceDriver {
	d := &PerforceDriver{base: newBase(repoConfig, env)}
	d.reader = d
	return d
}

// Initialize is PerforceDriver::initialize.
func (d *PerforceDriver) Initialize() error {
	v, ok := d.repoConfig.Get("depot")
	if !ok {
		return &RuntimeError{Msg: `Undefined array key "depot"`}
	}
	d.depot, _ = phpjson.ScalarString(v)
	d.branch = ""
	if b, ok := d.repoConfig.Get("branch"); ok && phpjson.Truthy(b) {
		d.branch, _ = phpjson.ScalarString(b)
	}
	if err := d.initPerforce(); err != nil {
		return err
	}
	if err := d.perforce.P4Login(); err != nil {
		return err
	}
	if _, err := d.perforce.CheckStream(); err != nil {
		return err
	}
	if err := d.perforce.WriteP4ClientSpec(); err != nil {
		return err
	}
	d.perforce.ConnectClient()
	return nil
}

func (d *PerforceDriver) initPerforce() error {
	if d.perforce != nil {
		return nil
	}
	vcsDir := d.Config.GetString("cache-vcs-dir")
	if !CacheUsable(vcsDir) {
		return &RuntimeError{Msg: "PerforceDriver requires a usable cache directory, and it looks like you set it to be disabled"}
	}
	p, err := NewPerforce(d.repoConfig, d.URL(), vcsDir+"/"+d.depot, d.Env)
	if err != nil {
		return err
	}
	d.perforce = p
	return nil
}

// FileContent is PerforceDriver::getFileContent.
func (d *PerforceDriver) FileContent(file, identifier string) (string, bool, error) {
	return d.perforce.FileContent(file, identifier)
}

// ChangeDate is PerforceDriver::getChangeDate (always null).
func (d *PerforceDriver) ChangeDate(string) (time.Time, bool, error) { return time.Time{}, false, nil }

// RootIdentifier is PerforceDriver::getRootIdentifier.
func (d *PerforceDriver) RootIdentifier() (string, error) { return d.branch, nil }

// Branches is PerforceDriver::getBranches.
func (d *PerforceDriver) Branches() ([]Ref, error) { return d.perforce.Branches() }

// Tags is PerforceDriver::getTags.
func (d *PerforceDriver) Tags() ([]Ref, error) { return d.perforce.Tags() }

// Dist is PerforceDriver::getDist.
func (d *PerforceDriver) Dist(string) (*phpjson.Map, error) { return nil, nil }

// Source is PerforceDriver::getSource.
func (d *PerforceDriver) Source(identifier string) *phpjson.Map {
	m := phpjson.NewMap()
	m.Set("type", "perforce")
	u, _ := d.repoConfig.Get("url")
	m.Set("url", u)
	m.Set("reference", identifier)
	if user, ok := d.perforce.User(); ok {
		m.Set("p4user", user)
	} else {
		m.Set("p4user", nil)
	}
	return m
}

// URL is PerforceDriver::getUrl.
func (d *PerforceDriver) URL() string { return d.url }

// HasComposerFile is PerforceDriver::hasComposerFile.
func (d *PerforceDriver) HasComposerFile(identifier string) (bool, error) {
	v, err := d.perforce.ComposerInformation("//" + d.depot + "/" + identifier)
	if err != nil {
		return false, err
	}
	return phpjson.Truthy(v), nil
}

// Cleanup is PerforceDriver::cleanup.
func (d *PerforceDriver) Cleanup() {
	if d.perforce != nil {
		_ = d.perforce.CleanupClientSpec()
		d.perforce = nil
	}
}
