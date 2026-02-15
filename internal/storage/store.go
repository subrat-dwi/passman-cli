package storage

type TokenStore interface {
	SaveAccessToken(token string) error
	GetAccessToken() (string, error)
	DeleteAccessToken() error
}

type SaltStore interface {
	SaveSalt(salt string) error
	GetSalt() ([]byte, error)
	DeleteSalt() error
}

type KeyVerifierStore interface {
	SaveKeyVerifier(ciphertext, nonce string) error
	GetKeyVerifier() (ciphertext, nonce string, err error)
}

type Store struct {
	Service string
	User    string
}
