package composer

import "gocomposer/internal/solver"

// LockTransaction is Installer::getLockTransaction (nil when no update ran).
func (in *Installer) LockTransaction() *solver.LockTransaction { return in.lockTransaction }
