package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/subrat-dwi/passman-cli/internal/usererror"
)

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
