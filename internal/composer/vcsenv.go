package composer

import (
	"sync"
	"time"

	"gocomposer/internal/vcs"
)

var vcsEnvs sync.Map // *Composer -> *vcs.Env

// VcsEnv returns the shared ProcessExecutor/IO/auth environment used by the
// VCS repositories and downloaders of this Composer instance.
func (c *Composer) VcsEnv() *vcs.Env {
	if e, ok := vcsEnvs.Load(c); ok {
		return e.(*vcs.Env)
	}
	timeout := 300 * time.Second
	if t, ok := c.Config.Get("process-timeout").(int64); ok {
		timeout = time.Duration(t) * time.Second
	}
	env := &vcs.Env{IO: c.IO, Config: c.Config, HTTP: c.HTTP, Process: vcs.NewProcess(c.IO, timeout), Auth: vcs.NewAuthStore(c.Config).Link(c.HTTP)}
	actual, _ := vcsEnvs.LoadOrStore(c, env)
	return actual.(*vcs.Env)
}

// FactoryInvalidArgumentError is an \InvalidArgumentException thrown while
// Factory::createComposer runs: Application::getComposer(true) prints its
// message as a plain line and exits with code 1 (no exception box).
type FactoryInvalidArgumentError struct{ Msg string }

func (e *FactoryInvalidArgumentError) Error() string { return e.Msg }
