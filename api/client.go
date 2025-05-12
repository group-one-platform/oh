package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/edvin/oh/cache"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

// Fetch performs an HTTP request with the given method, URL, headers, body
// The Bearer token from configuration is included into the Authorization header
func Fetch[T any](
	method, relativePath string,
	body any,
	cacheKey cache.CacheKey,
) (T, error) {
	var zero T

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	token := viper.GetString("token")
	if token == "" {
		return zero, fmt.Errorf("auth token is not set in the configuration")
	}
	headers["Authorization"] = "Bearer " + token

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return zero, fmt.Errorf("error marshalling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(b)
	}

	baseURL := viper.GetString("base_url")
	if baseURL == "" {
		return zero, fmt.Errorf("base_url is not set in the configuration")
	}
	url := baseURL + relativePath

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return zero, fmt.Errorf("error creating request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return zero, fmt.Errorf("%s %s request failed: %w", method, relativePath, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		// Some POST requests are StatusOK, so we deem 200/201 OK
	default:
		return zero, NewAPIError(resp, fmt.Sprintf("%s %s", method, relativePath))
	}

	var wrapper struct {
		Data T `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		io.Copy(io.Discard, resp.Body)
		return zero, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if viper.GetBool("no-cache") || cacheKey != cache.NoCache {
		cache.Store(cacheKey, &wrapper.Data)
	}
	return wrapper.Data, nil
}
