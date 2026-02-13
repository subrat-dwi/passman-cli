package api

import (
	"net/http"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}
