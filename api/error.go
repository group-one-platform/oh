package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Error struct {
	Action  string
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details struct {
		Reason        string            `json:"reason"`
		InvalidFields map[string]string `json:"invalid_fields"`
	} `json:"details"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s (%s)", e.Code, e.Message, e.Details.Reason)
}

func NewAPIError(resp *http.Response, action string) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body while executing %s: %w", action, err)
	}

	var apiErr Error
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return fmt.Errorf("failed to parse API error while executing %s: %w", action, err)
	}

	if apiErr.Code == 0 {
		apiErr.Code = resp.StatusCode
	}

	apiErr.Action = action

	return &apiErr
}
