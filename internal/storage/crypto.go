package storage

import (
	"encoding/hex"
	"strings"

	"github.com/subrat-dwi/passman-cli/internal/crypto"
	"github.com/subrat-dwi/passman-cli/internal/usererror"
	"github.com/zalando/go-keyring"
)

func (s *Store) SaveSalt(salt string) error {
	return keyring.Set(s.Service, s.User+"_salt", salt)
}

func (s *Store) GetSalt() ([]byte, error) {
	salt, err := keyring.Get(s.Service, s.User+"_salt")
	if err != nil {
		return nil, usererror.ErrNoSaltFound
	}
	// Try hex first, then base64
	if b, err := hex.DecodeString(salt); err == nil {
		return b, nil
	}
	return crypto.DecodeBase64(salt)
}

func (s *Store) DeleteSalt() error {
	return keyring.Delete(s.Service, s.User+"_salt")
}

// SaveKeyVerifier stores an encrypted verification token (ciphertext:nonce)
func (s *Store) SaveKeyVerifier(ciphertext, nonce string) error {
	return keyring.Set(s.Service, s.User+"_key_verifier", ciphertext+":"+nonce)
}

// GetKeyVerifier retrieves the encrypted verification token
func (s *Store) GetKeyVerifier() (ciphertext, nonce string, err error) {
	verifier, err := keyring.Get(s.Service, s.User+"_key_verifier")
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(verifier, ":", 2)
	if len(parts) != 2 {
		return "", "", usererror.New("Stored credentials are corrupted", "Please login again to fix this")
	}
	return parts[0], parts[1], nil
}
