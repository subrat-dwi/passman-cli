package storage

import (
	"github.com/zalando/go-keyring"
)

func (s *Store) SaveAccessToken(token string) error {
	return keyring.Set(s.Service, s.User+"_access_token", token)
}

func (s *Store) GetAccessToken() (string, error) {
	return keyring.Get(s.Service, s.User+"_access_token")
}

func (s *Store) DeleteAccessToken() error {
	return keyring.Delete(s.Service, s.User+"_access_token")
}
