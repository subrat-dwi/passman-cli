package service

// version holds the current version, set at build time via main package
var version = "dev"

// SetVersion sets the version (called from main)
func SetVersion(v string) {
	version = v
}

// Version returns the current version
func Version() string {
	return version
}
