package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/subrat-dwi/passman-cli/internal/agent"
	"github.com/subrat-dwi/passman-cli/internal/api"
	"github.com/subrat-dwi/passman-cli/internal/crypto"
	"github.com/subrat-dwi/passman-cli/internal/passwordmanager"
	"github.com/subrat-dwi/passman-cli/internal/storage"
	"github.com/subrat-dwi/passman-cli/internal/usererror"
	"golang.org/x/term"
)

type PasswordService struct {
	API     *api.PasswordAPI
	Storage *storage.Store
}

// promptAndUnlock prompts for master password, verifies it, and unlocks the agent
func (s *PasswordService) PromptAndUnlock() error {
	salt, err := s.Storage.GetSalt()
	if err != nil {
		return usererror.ErrNoSaltFound
	}

	fmt.Print("\nVault locked. Enter master password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return usererror.New("Could not read password", "Make sure you're running this in a terminal")
	}

	if len(passwordBytes) == 0 {
		return usererror.New("Password cannot be empty", "Please enter your master password")
	}

	key := crypto.DeriveKey(string(passwordBytes), crypto.EncodeBase64(salt))
	if err := agent.Unlock(crypto.EncodeBase64(key), 600); err != nil {
		return usererror.Wrap(usererror.ErrAgentConnection, err)
	}

	// Verify the password by decrypting the key verifier
	ciphertext, nonce, err := s.Storage.GetKeyVerifier()
	if err != nil {
		// No verifier stored - older login, allow but warn
		fmt.Println("  Note: Password verification not available. Re-login to enable it.")
		return nil
	}

	plaintext, err := agent.Decrypt(ciphertext, nonce)
	if err != nil || plaintext != KeyVerifierPlaintext {
		// Wrong password - lock agent and return error
		agent.Lock()
		return usererror.ErrInvalidPassword
	}

	fmt.Println("  Vault unlocked successfully.")
	return nil
}

// isAgentLocked checks if the error is due to agent being locked
func isAgentLocked(err error) bool {
	return err != nil && strings.Contains(err.Error(), "agent locked")
}

func (s *PasswordService) Create(name, username, password string) error {
	token, err := s.Storage.GetAccessToken()
	if err != nil {
		return usererror.ErrNotLoggedIn
	}

	ciphertext, nonce, err := agent.Encrypt(password)
	if isAgentLocked(err) {
		if unlockErr := s.PromptAndUnlock(); unlockErr != nil {
			return unlockErr
		}
		// Retry after unlock
		ciphertext, nonce, err = agent.Encrypt(password)
	}
	if err != nil {
		return usererror.Wrap(usererror.ErrEncryptFailed, err)
	}

	return s.API.CreatePassword(token, passwordmanager.PasswordFullEntry{
		Name:     name,
		Username: username,
		Password: ciphertext,
		Nonce:    nonce,
	})
}

func (s *PasswordService) List() ([]passwordmanager.PasswordEntry, error) {
	token, err := s.Storage.GetAccessToken()
	if err != nil {
		return nil, usererror.ErrNotLoggedIn
	}

	passwords, err := s.API.ListPasswords(token)
	if err != nil {
		return nil, err
	}

	var results []passwordmanager.PasswordEntry
	for _, p := range passwords {
		results = append(results, passwordmanager.PasswordEntry{
			ID:       p.ID,
			Name:     p.Name,
			Username: p.Username,
		})
	}
	return results, nil
}

// DecryptedPassword holds the full password details with plaintext password
type DecryptedPassword struct {
	ID        string
	Name      string
	Username  string
	Password  string
	CreatedAt string
	UpdatedAt string
}

func (s *PasswordService) Get(id string) (*DecryptedPassword, error) {
	token, err := s.Storage.GetAccessToken()
	if err != nil {
		return nil, usererror.ErrNotLoggedIn
	}

	entry, err := s.API.GetPassword(token, id)
	if err != nil {
		return nil, err
	}

	// Decrypt the password using the agent
	plaintext, err := agent.Decrypt(entry.Password, entry.Nonce)
	if isAgentLocked(err) {
		if unlockErr := s.PromptAndUnlock(); unlockErr != nil {
			return nil, unlockErr
		}
		// Retry after unlock
		plaintext, err = agent.Decrypt(entry.Password, entry.Nonce)
	}
	if err != nil {
		return nil, usererror.Wrap(usererror.ErrDecryptFailed, err)
	}

	return &DecryptedPassword{
		ID:        entry.ID,
		Name:      entry.Name,
		Username:  entry.Username,
		Password:  plaintext,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}, nil
}

func (s *PasswordService) Update(id, name, username, password string) error {
	token, err := s.Storage.GetAccessToken()
	if err != nil {
		return usererror.ErrNotLoggedIn
	}

	ciphertext, nonce, err := agent.Encrypt(password)
	if isAgentLocked(err) {
		if unlockErr := s.PromptAndUnlock(); unlockErr != nil {
			return unlockErr
		}
		// Retry after unlock
		ciphertext, nonce, err = agent.Encrypt(password)
	}
	if err != nil {
		return usererror.Wrap(usererror.ErrEncryptFailed, err)
	}

	return s.API.UpdatePassword(token, id, passwordmanager.PasswordFullEntry{
		Name:     name,
		Username: username,
		Password: ciphertext,
		Nonce:    nonce,
	})
}

func (s *PasswordService) Delete(id string) error {
	token, err := s.Storage.GetAccessToken()
	if err != nil {
		return usererror.ErrNotLoggedIn
	}

	return s.API.DeletePassword(token, id)
}
