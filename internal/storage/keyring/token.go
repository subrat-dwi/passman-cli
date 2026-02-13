package keyring

import (
	"github.com/zalando/go-keyring"
)

type TokenStore struct {
	Service string
	User    string
}

func (t *TokenStore) SaveAccessToken(token string) error {
	return keyring.Set(t.Service, t.User+"_access", token)
}

func (t *TokenStore) GetAccessToken() (string, error) {
	return keyring.Get(t.Service, t.User+"_access")
}
