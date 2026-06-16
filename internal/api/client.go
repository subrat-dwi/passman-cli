package api

import (
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func (c *Client) url(path string) (string, error) {
	if c.BaseURL == "" {
		return "", fmt.Errorf("backend server URL not configured.\nRun: pman config set api_base_url <your-server-url>")
	}
	return c.BaseURL + path, nil
}
