package agent

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/subrat-dwi/passman-cli/internal/ipc"
)

// Unlock tells the agent to store the provided base64-encoded key and start its TTL countdown.
func Unlock(keyBase64 string, ttlSeconds int) error {
	_, err := send(Request{Cmd: "unlock", Key: keyBase64, TTL: ttlSeconds})
	return err
}

// Lock wipes the key from the agent.
func Lock() error {
	_, err := send(Request{Cmd: "lock"})
	return err
}

// Status returns whether the agent is unlocked and the remaining seconds until auto-lock.
func Status() (unlocked bool, expiresSeconds int, err error) {
	resp, err := send(Request{Cmd: "status"})
	if err != nil {
		return false, 0, err
	}

	m, ok := resp.Data.(map[string]any)
	if !ok {
		return false, 0, fmt.Errorf("unexpected status payload: %T", resp.Data)
	}

	if v, ok := m["unlocked"].(bool); ok {
		unlocked = v
	}
	if v, ok := m["expires"].(float64); ok {
		expiresSeconds = int(v)
	}
	return unlocked, expiresSeconds, nil
}

// Encrypt asks the agent to encrypt the provided plaintext using the stored key and returns the base64-encoded ciphertext and nonce.
func Encrypt(plaintext string) (ciphertextBase64, nonceBase64 string, err error) {
	resp, err := send(Request{Cmd: "encrypt", Plaintext: plaintext})
	if err != nil {
		return "", "", err
	}

	m, ok := resp.Data.(map[string]any)
	if !ok {
		return "", "", fmt.Errorf("unexpected encrypt payload: %T", resp.Data)
	}

	return m["ciphertext"].(string), m["nonce"].(string), nil
}

// Decrypt asks the agent to decrypt the provided base64 ciphertext using the stored key and nonce.
func Decrypt(ciphertextBase64, nonceBase64 string) (string, error) {
	resp, err := send(Request{Cmd: "decrypt", Ciphertext: ciphertextBase64, Nonce: nonceBase64})
	if err != nil {
		return "", err
	}

	plaintextBase64, ok := resp.Data.(string)
	if !ok {
		return "", fmt.Errorf("unexpected decrypt payload: %T", resp.Data)
	}

	// Decode base64 plaintext
	plaintext, err := base64.RawStdEncoding.DecodeString(plaintextBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode plaintext: %w", err)
	}
	return string(plaintext), nil
}

// send connects to the agent via IPC, sends the provided request, and waits for a response. It returns the response or an error if the operation fails.
func send(req Request) (Response, error) {
	conn, err := ipc.Dial()
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return Response{}, err
	}

	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return Response{}, err
	}

	if !resp.OK {
		return resp, errors.New(resp.Error)
	}
	return resp, nil
}
