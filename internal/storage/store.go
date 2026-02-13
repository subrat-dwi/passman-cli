package storage

type TokenStore interface {
	SaveAccessToken(token string) error
	GetAccessToken() (string, error)
}

type CryptoStore interface {
	SaveSalt(salt []byte) error
	GetSalt() ([]byte, error)
}
