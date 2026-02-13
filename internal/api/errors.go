package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type APIError struct {
	Status  int
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error (%d): %s", e.Status, e.Message)
}

func parseAPIError(resp *http.Response) error {
	var body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return &APIError{
			Status:  resp.StatusCode,
			Message: "unknown api error",
		}
	}

	if err := json.Unmarshal(data, &body); err != nil {
		return &APIError{
			Status:  resp.StatusCode,
			Message: string(data),
		}
	}

	status, _ := strconv.Atoi(body.Status)

	return &APIError{
		Status:  status,
		Message: body.Message,
	}
}
