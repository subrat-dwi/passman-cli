package api

import (
	"encoding/json"
	"net/http"
)

type PasswordAPI struct {
	client *Client
}

func NewPasswordAPI(c *Client) *PasswordAPI {
	return &PasswordAPI{
		client: c,
	}
}

type PasswordItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListPasswordsResponse struct {
	Passwords []PasswordItem `json:"passwords"`
}

func (p *PasswordAPI) ListPasswords(accessToken string) ([]PasswordItem, error) {
	req, err := http.NewRequest(
		"GET",
		p.client.BaseURL+"/api/passwords",
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.client.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	var out ListPasswordsResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out.Passwords, err
}
