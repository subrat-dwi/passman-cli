package app

import (
	"net/http"
	"time"

	"github.com/subrat-dwi/passman-cli/internal/api"
	"github.com/subrat-dwi/passman-cli/internal/service"
	"github.com/subrat-dwi/passman-cli/internal/storage/keyring"
)

type App struct {
	AuthService *service.AuthService
}

func New() *App {
	baseURL := "https://shubserver.onrender.com"

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client := &api.Client{
		BaseURL: baseURL,
		HTTP:    httpClient,
	}

	authAPI := api.NewAuthAPI(client)

	authService := &service.AuthService{
		API: authAPI,
		Storage: &keyring.TokenStore{
			Service: "passman",
			User:    "default",
		},
	}

	return &App{
		AuthService: authService,
	}
}
