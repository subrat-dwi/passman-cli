package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type AuthAPI struct {
	client *Client
}

func NewAuthAPI(c *Client) *AuthAPI {
	return &AuthAPI{
		client: c,
	}
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"token"`
	Salt        string `json:"salt"`
}

type RegisterResponse struct {
	AccessToken string `json:"token"`
	Salt        string `json:"salt"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"token"`
	Salt        string `json:"salt"`
}

func (a *AuthAPI) Register(email, password string) (*RegisterResponse, error) {
	body, _ := json.Marshal(RegisterRequest{
		Email:    email,
		Password: password,
	})

	req, _ := http.NewRequest(
		"POST",
		a.client.BaseURL+"/api/users/register",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.HTTP.Do(req)
	if err != nil {
		return nil, wrapNetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, parseAPIError(resp)
	}

	var out RegisterResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return &out, err
}

func (a *AuthAPI) Login(email, password string) (*LoginResponse, error) {
	body, _ := json.Marshal(LoginRequest{
		Email:    email,
		Password: password,
	})

	req, _ := http.NewRequest(
		"POST",
		a.client.BaseURL+"/api/users/login",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.HTTP.Do(req)
	if err != nil {
		return nil, wrapNetworkError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	var out LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return &out, err
}
