package agent

import (
	"sync"
	"time"

	"github.com/subrat-dwi/passman-cli/internal/crypto"
)

// State holds the agent's state, including whether it's unlocked, the encryption key, last activity time, and TTL.
type State struct {
	unlocked     bool
	key          []byte
	lastActivity time.Time
	ttl          time.Duration
	mu           sync.Mutex
}

func NewState(ttl time.Duration) *State {
	return &State{
		unlocked: false,
		ttl:      ttl,
	}
}

// Handle processes incoming requests and returns appropriate responses based on the command.
func (s *State) Handle(req Request) Response {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch req.Cmd {

	case "status":
		return s.Status()

	case "unlock":
		return s.Unlock(req)

	case "lock":
		s.Lock()
		return ok()

	case "encrypt":
		return s.Encrpt(req)

	case "decrypt":
		return s.Decrypt(req)

	default:
		return errResp("unknown command")
	}
}

// status returns the current status of the agent, including whether it's unlocked and the remaining time until it auto-locks.
func (s *State) Status() Response {
	if !s.unlocked {
		return data(map[string]any{"unlocked": false})
	}
	remaining := int(s.ttl.Seconds() - time.Since(s.lastActivity).Seconds())
	return data(map[string]any{
		"unlocked": true,
		"expires":  remaining,
	})
}

// unlock attempts to unlock the agent using the provided key. It decodes the key from base64 and updates the state accordingly.
func (s *State) Unlock(req Request) Response {
	key, err := crypto.DecodeBase64(req.Key)
	if err != nil {
		return errResp("bad key")
	}

	s.key = key
	s.unlocked = true
	s.ttl = time.Duration(req.TTL) * time.Second
	s.lastActivity = time.Now()
	return ok()
}

// lock securely wipes the encryption key from memory and updates the state to reflect that the agent is now locked.
func (s *State) Lock() {
	if s.key != nil {
		for i := range s.key {
			s.key[i] = 0
		}
	}
	s.key = nil
	s.unlocked = false
}

func (s *State) Encrpt(req Request) Response {
	if !s.unlocked {
		return errResp("agent locked")
	}

	ciphertext, nonce, err := crypto.Encrypt([]byte(req.Plaintext), []byte(s.key))
	if err != nil {
		return errResp("encrypt failed")
	}

	s.lastActivity = time.Now()
	return data(map[string]any{
		"ciphertext": crypto.EncodeBase64(ciphertext),
		"nonce":      crypto.EncodeBase64(nonce),
	})
}

// decrypt attempts to decrypt the provided ciphertext using the stored key. It checks if the agent is unlocked, decodes the ciphertext and nonce from base64, and returns the decrypted plaintext or an error response if decryption fails.
func (s *State) Decrypt(req Request) Response {
	if !s.unlocked {
		return errResp("agent locked")
	}

	ciphertextBytes, _ := crypto.DecodeBase64(req.Ciphertext)

	var nonce []byte
	if req.Nonce != "" {
		nonce, _ = crypto.DecodeBase64(req.Nonce)
	}

	s.lastActivity = time.Now()
	plaintext, err := crypto.DecryptWithNonce(ciphertextBytes, nonce, s.key)
	if err != nil {
		return errResp("decrypt failed")
	}

	s.lastActivity = time.Now()
	return data(crypto.EncodeBase64(plaintext))
}

// autoLock runs in a separate goroutine and periodically checks if the agent should be automatically locked based on the last activity time and the configured TTL. If the agent has been unlocked for longer than the TTL, it locks the agent.
func (s *State) AutoLock() {
	for {
		time.Sleep(5 * time.Second)
		s.mu.Lock()
		if s.unlocked && time.Since(s.lastActivity) > s.ttl {
			s.Lock()
		}
		s.mu.Unlock()
	}
}
