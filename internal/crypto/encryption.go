package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

// Encrypt encrypts the plaintext using AES-GCM with the provided key.
func Encrypt(plaintext, key []byte) ([]byte, []byte, error) {
	// Create a new AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// create a new GCM cipher mode instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the plaintext using the GCM cipher (don't prepend nonce)
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// Decrypt decrypts the ciphertext using AES-GCM with the provided key and nonce.
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	return DecryptWithNonce(ciphertext, nil, key)
}

// DecryptWithNonce decrypts ciphertext with explicit nonce. If nonce is nil, expects nonce prepended to ciphertext.
func DecryptWithNonce(ciphertext, nonce, key []byte) ([]byte, error) {
	// Create a new AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// create a new GCM cipher mode instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// If nonce not provided, extract from ciphertext prefix
	if nonce == nil {
		if len(ciphertext) < gcm.NonceSize() {
			return nil, errors.New("ciphertext too short")
		}
		nonce = ciphertext[:gcm.NonceSize()]
		ciphertext = ciphertext[gcm.NonceSize():]
	}

	// Decrypt the ciphertext using the GCM cipher
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}
