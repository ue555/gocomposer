package semver

// BranchNames returns Intervals::get($constraint)['branches']['names']
// (composer/semver, MIT, as locked by composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80).
func BranchNames(c Constraint) []string {
	return append([]string(nil), getIntervals(c).branches.names...)
}
