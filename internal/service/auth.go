package service

import (
	"fmt"

	"github.com/subrat-dwi/passman-cli/internal/agent"
	"github.com/subrat-dwi/passman-cli/internal/api"
	"github.com/subrat-dwi/passman-cli/internal/crypto"
	"github.com/subrat-dwi/passman-cli/internal/storage"
)

// KeyVerifierPlaintext is the known value used to verify the master password
const KeyVerifierPlaintext = "passman-key-verify"

type AuthService struct {
	API     *api.AuthAPI
	Storage *storage.Store
}

func (s *AuthService) Register(email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("email and password cannot be empty")
	}

	resp, err := s.API.Register(email, password)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	if err := s.Storage.SaveSalt(resp.Salt); err != nil {
		return fmt.Errorf("failed to save salt: %w", err)
	}

	if err := s.Storage.SaveAccessToken(resp.AccessToken); err != nil {
		return fmt.Errorf("failed to save access token: %w", err)
	}

	key := crypto.DeriveKey(password, resp.Salt)
	if err := agent.Unlock(crypto.EncodeBase64(key), 600); err != nil {
		return fmt.Errorf("failed to unlock agent: %w", err)
	}

	// Save key verifier for password validation later
	ciphertext, nonce, err := agent.Encrypt(KeyVerifierPlaintext)
	if err != nil {
		return fmt.Errorf("failed to create key verifier: %w", err)
	}
	if err := s.Storage.SaveKeyVerifier(ciphertext, nonce); err != nil {
		return fmt.Errorf("failed to save key verifier: %w", err)
	}

	return nil
}

func (s *AuthService) Login(email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("email and password cannot be empty")
	}

	resp, err := s.API.Login(email, password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	if err := s.Storage.SaveSalt(resp.Salt); err != nil {
		return fmt.Errorf("failed to save salt: %w", err)
	}

	if err := s.Storage.SaveAccessToken(resp.AccessToken); err != nil {
		return fmt.Errorf("failed to save access token: %w", err)
	}

	key := crypto.DeriveKey(password, resp.Salt)
	if err := agent.Unlock(crypto.EncodeBase64(key), 600); err != nil {
		return fmt.Errorf("failed to unlock agent: %w", err)
	}

	// Save key verifier for password validation later
	ciphertext, nonce, err := agent.Encrypt(KeyVerifierPlaintext)
	if err != nil {
		return fmt.Errorf("failed to create key verifier: %w", err)
	}
	if err := s.Storage.SaveKeyVerifier(ciphertext, nonce); err != nil {
		return fmt.Errorf("failed to save key verifier: %w", err)
	}

	return nil
}
