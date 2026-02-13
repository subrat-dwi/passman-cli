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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"token"`
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
		return nil, err
	}
	defer resp.Body.Close()

	var out LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return &out, err
}
