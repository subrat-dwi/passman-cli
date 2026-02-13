package service

import (
	"github.com/subrat-dwi/passman-cli/internal/api"
	"github.com/subrat-dwi/passman-cli/internal/storage"
)

type PasswordService struct {
	API     *api.PasswordAPI
	Storage storage.TokenStore
}

func (s *PasswordService) List() ([]api.PasswordItem, error) {
	token, err := s.Storage.GetAccessToken()
	if err != nil {
		return nil, err
	}

	passwords, err := s.API.ListPasswords(token)
	if err != nil {
		return nil, err
	}

	var results []api.PasswordItem
	for _, p := range passwords {
		results = append(results, api.PasswordItem{
			ID:       p.ID,
			Name:     p.Name,
			Username: p.Username,
		})
	}
	return results, nil
}
