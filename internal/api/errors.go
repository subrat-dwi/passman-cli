package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/subrat-dwi/passman-cli/internal/usererror"
)

// wrapNetworkError converts low-level network errors into user-friendly messages
func wrapNetworkError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := strings.ToLower(err.Error())

	// Timeout errors (including Render cold start)
	if strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "deadline exceeded") ||
		strings.Contains(errMsg, "context deadline") {
		return usererror.ErrTimeout
	}

	// Connection errors
	if strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "no such host") ||
		strings.Contains(errMsg, "network is unreachable") {
		return usererror.ErrServerUnreachable
	}

	return usererror.ErrServerUnreachable
}

type APIError struct {
	Status  int
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return usererror.FromAPIError(e.Status, e.Message).Error()
}

func parseAPIError(resp *http.Response) error {
	var body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return usererror.FromAPIError(resp.StatusCode, "Server returned an empty response")
	}

	if err := json.Unmarshal(data, &body); err != nil {
		// Try to extract a readable message from raw response
		raw := strings.TrimSpace(string(data))
		if len(raw) > 100 {
			raw = raw[:100] + "..."
		}
		return usererror.FromAPIError(resp.StatusCode, raw)
	}

	status, _ := strconv.Atoi(body.Status)
	if status == 0 {
		status = resp.StatusCode
	}

	return &APIError{
		Status:  status,
		Message: body.Message,
	}
}
