//go:build linux
// +build linux

package storage

// On Linux (including WSL), keyring access often fails or crashes.
// Always use file-based storage for reliability.

func checkKeyringAvailable() bool {
	return false
}

// secureSet stores a value using file-based storage
func secureSet(service, key, value string) error {
	fs, err := getFileStore()
	if err != nil {
		return err
	}
	return fs.Set(service, key, value)
}

// secureGet retrieves a value from file-based storage
func secureGet(service, key string) (string, error) {
	fs, err := getFileStore()
	if err != nil {
		return "", err
	}
	return fs.Get(service, key)
}

// secureDelete removes a value from file-based storage
func secureDelete(service, key string) error {
	fs, err := getFileStore()
	if err != nil {
		return err
	}
	return fs.Delete(service, key)
}
