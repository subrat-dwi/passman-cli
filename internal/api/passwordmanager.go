package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/subrat-dwi/passman-cli/internal/passwordmanager"
)

type PasswordAPI struct {
	client *Client
}

func NewPasswordAPI(c *Client) *PasswordAPI {
	return &PasswordAPI{
		client: c,
	}
}

type ListPasswordsResponse struct {
	Passwords []passwordmanager.PasswordEntry `json:"passwords"`
}

func (p *PasswordAPI) ListPasswords(accessToken string) ([]passwordmanager.PasswordEntry, error) {
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
		return nil, wrapNetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	var out ListPasswordsResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out.Passwords, err
}

func (p *PasswordAPI) CreatePassword(accessToken string, entry passwordmanager.PasswordFullEntry) error {
	payload, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		p.client.BaseURL+"/api/passwords",
		bytes.NewReader(payload),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.HTTP.Do(req)
	if err != nil {
		return wrapNetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return parseAPIError(resp)
	}

	return nil
}

func (p *PasswordAPI) GetPassword(accessToken, id string) (*passwordmanager.PasswordFullEntry, error) {
	req, err := http.NewRequest(
		"GET",
		p.client.BaseURL+"/api/passwords/"+id,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.client.HTTP.Do(req)
	if err != nil {
		return nil, wrapNetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	var out passwordmanager.PasswordFullEntry
	err = json.NewDecoder(resp.Body).Decode(&out)
	return &out, err
}

func (p *PasswordAPI) UpdatePassword(accessToken, id string, entry passwordmanager.PasswordFullEntry) error {
	payload, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PUT",
		p.client.BaseURL+"/api/passwords/"+id,
		bytes.NewReader(payload),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.HTTP.Do(req)
	if err != nil {
		return wrapNetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseAPIError(resp)
	}

	return nil
}

func (p *PasswordAPI) DeletePassword(accessToken, id string) error {
	req, err := http.NewRequest(
		"DELETE",
		p.client.BaseURL+"/api/passwords/"+id,
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.client.HTTP.Do(req)
	if err != nil {
		return wrapNetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return parseAPIError(resp)
	}

	return nil
}
