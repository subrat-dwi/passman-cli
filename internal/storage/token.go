package storage

func (s *Store) SaveAccessToken(token string) error {
	return secureSet(s.Service, s.User+"_access_token", token)
}

func (s *Store) GetAccessToken() (string, error) {
	return secureGet(s.Service, s.User+"_access_token")
}

func (s *Store) DeleteAccessToken() error {
	return secureDelete(s.Service, s.User+"_access_token")
}
