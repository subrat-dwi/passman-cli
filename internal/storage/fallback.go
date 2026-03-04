package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// FileStore provides encrypted file-based storage as a fallback when keyring is unavailable
type FileStore struct {
	dir  string
	key  []byte
	mu   sync.RWMutex
	data map[string]string
}

var (
	fileStore     *FileStore
	fileStoreOnce sync.Once
	fileStoreErr  error
)

// getFileStore returns the singleton FileStore instance
func getFileStore() (*FileStore, error) {
	fileStoreOnce.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fileStoreErr = err
			return
		}

		storeDir := filepath.Join(homeDir, ".passman")
		if err := os.MkdirAll(storeDir, 0700); err != nil {
			fileStoreErr = err
			return
		}

		// Derive encryption key from machine-specific data
		key := deriveMachineKey()

		fileStore = &FileStore{
			dir:  storeDir,
			key:  key,
			data: make(map[string]string),
		}

		// Load existing data
		fileStore.load()
	})

	return fileStore, fileStoreErr
}

// deriveMachineKey creates a key from machine-specific information
func deriveMachineKey() []byte {
	// Combine multiple sources for uniqueness
	var machineData string

	// Home directory
	if home, err := os.UserHomeDir(); err == nil {
		machineData += home
	}

	// Hostname
	if hostname, err := os.Hostname(); err == nil {
		machineData += hostname
	}

	// OS and architecture
	machineData += runtime.GOOS + runtime.GOARCH

	// Username from environment
	if user := os.Getenv("USER"); user != "" {
		machineData += user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		machineData += user
	}

	// Application-specific salt
	machineData += "passman-cli-local-storage-v1"

	// Hash to get a 32-byte key for AES-256
	hash := sha256.Sum256([]byte(machineData))
	return hash[:]
}

func (fs *FileStore) filePath() string {
	return filepath.Join(fs.dir, "vault.enc")
}

func (fs *FileStore) load() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := os.ReadFile(fs.filePath())
	if err != nil {
		if os.IsNotExist(err) {
			fs.data = make(map[string]string)
			return nil
		}
		return err
	}

	// Decrypt data
	plaintext, err := fs.decrypt(data)
	if err != nil {
		// Corrupted file, start fresh
		fs.data = make(map[string]string)
		return nil
	}

	return json.Unmarshal(plaintext, &fs.data)
}

func (fs *FileStore) save() error {
	plaintext, err := json.Marshal(fs.data)
	if err != nil {
		return err
	}

	ciphertext, err := fs.encrypt(plaintext)
	if err != nil {
		return err
	}

	// Write atomically
	tmpPath := fs.filePath() + ".tmp"
	if err := os.WriteFile(tmpPath, ciphertext, 0600); err != nil {
		return err
	}

	return os.Rename(tmpPath, fs.filePath())
}

func (fs *FileStore) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(fs.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Prepend nonce to ciphertext
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (fs *FileStore) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(fs.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Set stores a value
func (fs *FileStore) Set(service, key, value string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	storageKey := service + ":" + key
	fs.data[storageKey] = value
	return fs.save()
}

// Get retrieves a value
func (fs *FileStore) Get(service, key string) (string, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	storageKey := service + ":" + key
	value, ok := fs.data[storageKey]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

// Delete removes a value
func (fs *FileStore) Delete(service, key string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	storageKey := service + ":" + key
	delete(fs.data, storageKey)
	return fs.save()
}
