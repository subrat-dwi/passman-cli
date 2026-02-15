package app

import (
	"net/http"
	"time"

	"github.com/subrat-dwi/passman-cli/internal/api"
	"github.com/subrat-dwi/passman-cli/internal/config"
	"github.com/subrat-dwi/passman-cli/internal/service"
	"github.com/subrat-dwi/passman-cli/internal/storage"
)

// App is the main application struct that holds all services and configurations.
type App struct {
	AuthService     *service.AuthService
	PasswordService *service.PasswordService
}

func New() *App {
	// Initialize configuration
	config.Init()
	baseURL := config.Get("api_base_url")

	// Set up HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Initialize API clients
	client := &api.Client{
		BaseURL: baseURL,
		HTTP:    httpClient,
	}

	authAPI := api.NewAuthAPI(client)
	passwordAPI := api.NewPasswordAPI(client)

	// Initialize services
	storage := storage.Store{
		Service: "passman",
		User:    "default",
	}

	authService := &service.AuthService{
		API:     authAPI,
		Storage: &storage,
	}

	passwordService := &service.PasswordService{
		API:     passwordAPI,
		Storage: &storage,
	}

	// Return the initialized App instance
	return &App{
		AuthService:     authService,
		PasswordService: passwordService,
	}
}

func (a *App) ResetState() {
	a.AuthService = nil
	a.PasswordService = nil
}
