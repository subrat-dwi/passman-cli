//go:build !linux
// +build !linux

package storage

import (
	"sync"

	"github.com/zalando/go-keyring"
)

var (
	keyringAvailable bool
	keyringChecked   bool
	keyringCheckMu   sync.Mutex
)

// checkKeyringAvailable tests if the system keyring is working
func checkKeyringAvailable() bool {
	keyringCheckMu.Lock()
	defer keyringCheckMu.Unlock()

	if keyringChecked {
		return keyringAvailable
	}

	// Use defer/recover to catch any panics from keyring library
	defer func() {
		if r := recover(); r != nil {
			keyringChecked = true
			keyringAvailable = false
		}
	}()

	// Test keyring by setting and getting a test value
	testService := "passman-keyring-test"
	testKey := "availability-check"
	testValue := "test"

	err := keyring.Set(testService, testKey, testValue)
	if err != nil {
		keyringChecked = true
		keyringAvailable = false
		return false
	}

	_, err = keyring.Get(testService, testKey)
	if err != nil {
		keyringChecked = true
		keyringAvailable = false
		return false
	}

	// Cleanup
	_ = keyring.Delete(testService, testKey)

	keyringChecked = true
	keyringAvailable = true
	return true
}

// secureSet stores a value, trying keyring first then falling back to file
func secureSet(service, key, value string) error {
	if checkKeyringAvailable() {
		err := keyring.Set(service, key, value)
		if err == nil {
			return nil
		}
	}

	// Use file-based fallback
	fs, err := getFileStore()
	if err != nil {
		return err
	}
	return fs.Set(service, key, value)
}

// secureGet retrieves a value, trying keyring first then file fallback
func secureGet(service, key string) (string, error) {
	if checkKeyringAvailable() {
		value, err := keyring.Get(service, key)
		if err == nil {
			return value, nil
		}
	}

	// Try file-based fallback
	fs, err := getFileStore()
	if err != nil {
		return "", err
	}
	return fs.Get(service, key)
}

// secureDelete removes a value from both storage backends
func secureDelete(service, key string) error {
	var keyringErr, fileErr error

	// Try to delete from keyring
	if checkKeyringAvailable() {
		keyringErr = keyring.Delete(service, key)
	}

	// Also try to delete from file fallback
	fs, err := getFileStore()
	if err == nil {
		fileErr = fs.Delete(service, key)
	}

	// Return nil if either succeeded
	if keyringErr == nil || fileErr == nil {
		return nil
	}

	if keyringAvailable && keyringErr != nil {
		return keyringErr
	}
	return fileErr
}
