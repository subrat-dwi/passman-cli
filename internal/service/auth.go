package service

import (
	"github.com/subrat-dwi/passman-cli/internal/api"
	"github.com/subrat-dwi/passman-cli/internal/storage"
)

type AuthService struct {
	API     *api.AuthAPI
	Storage storage.TokenStore
}

func (s *AuthService) Login(email, password string) error {
	resp, err := s.API.Login(email, password)
	if err != nil {
		return err
	}

	return s.Storage.SaveAccessToken(resp.AccessToken)
}
