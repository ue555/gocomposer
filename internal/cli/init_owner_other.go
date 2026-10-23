//go:build !unix

package cli

// currentScriptOwner mirrors PHP's get_current_user(); file owners are not
// available here.
func currentScriptOwner() string { return "" }
