package service

import (
	"github.com/subrat-dwi/passman-cli/internal/agent"
	"github.com/subrat-dwi/passman-cli/internal/api"
	"github.com/subrat-dwi/passman-cli/internal/crypto"
	"github.com/subrat-dwi/passman-cli/internal/storage"
	"github.com/subrat-dwi/passman-cli/internal/usererror"
)

// KeyVerifierPlaintext is the known value used to verify the master password
const KeyVerifierPlaintext = "passman-key-verify"

type AuthService struct {
	API     *api.AuthAPI
	Storage *storage.Store
}

func (s *AuthService) Register(email, password string) error {
	if email == "" || password == "" {
		return usererror.ErrEmptyFields
	}

	resp, err := s.API.Register(email, password)
	if err != nil {
		return err // API errors are already user-friendly
	}

	if err := s.Storage.SaveSalt(resp.Salt); err != nil {
		return usererror.Wrap(usererror.ErrKeyringAccess, err)
	}

	if err := s.Storage.SaveAccessToken(resp.AccessToken); err != nil {
		return usererror.Wrap(usererror.ErrKeyringAccess, err)
	}

	key := crypto.DeriveKey(password, resp.Salt)
	if err := agent.Unlock(crypto.EncodeBase64(key), 600); err != nil {
		return usererror.Wrap(usererror.ErrAgentConnection, err)
	}

	// Save key verifier for password validation later
	ciphertext, nonce, err := agent.Encrypt(KeyVerifierPlaintext)
	if err != nil {
		return usererror.Wrap(usererror.ErrEncryptFailed, err)
	}
	if err := s.Storage.SaveKeyVerifier(ciphertext, nonce); err != nil {
		return usererror.Wrap(usererror.ErrKeyringAccess, err)
	}

	return nil
}

func (s *AuthService) Login(email, password string) error {
	if email == "" || password == "" {
		return usererror.ErrEmptyFields
	}

	resp, err := s.API.Login(email, password)
	if err != nil {
		return err // API errors are already user-friendly
	}

	if err := s.Storage.SaveSalt(resp.Salt); err != nil {
		return usererror.Wrap(usererror.ErrKeyringAccess, err)
	}

	if err := s.Storage.SaveAccessToken(resp.AccessToken); err != nil {
		return usererror.Wrap(usererror.ErrKeyringAccess, err)
	}

	key := crypto.DeriveKey(password, resp.Salt)
	if err := agent.Unlock(crypto.EncodeBase64(key), 600); err != nil {
		return usererror.Wrap(usererror.ErrAgentConnection, err)
	}

	// Save key verifier for password validation later
	ciphertext, nonce, err := agent.Encrypt(KeyVerifierPlaintext)
	if err != nil {
		return usererror.Wrap(usererror.ErrEncryptFailed, err)
	}
	if err := s.Storage.SaveKeyVerifier(ciphertext, nonce); err != nil {
		return usererror.Wrap(usererror.ErrKeyringAccess, err)
	}

	return nil
}
