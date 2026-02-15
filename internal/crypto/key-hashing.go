package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

// GenerateSalt generates a random salt of the specified length.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16) // 16 bytes salt
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey derives an encryption key from the password and salt using Argon2id
func DeriveKey(masterPassword, salt string) []byte {

	saltBytes, err := DecodeBase64(salt)
	if err != nil {
		panic(err)
	}
	return argon2.IDKey([]byte(masterPassword), saltBytes, 1, 64*1024, 4, 32)
}

// Encoder/Decoder helpers for storing the salt
func EncodeBase64(data []byte) string {
	return base64.RawStdEncoding.EncodeToString(data)
}

func DecodeBase64(s string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(s)
}
