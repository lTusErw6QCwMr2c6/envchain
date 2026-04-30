package secret

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// newHTTPRequest builds an HTTP request with Vault token authentication.
func newHTTPRequest(method, url, token, body string) (*http.Request, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// doHTTPRequest executes a request and checks for the expected status code.
func doHTTPRequest(req *http.Request, expectedStatus int) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedStatus {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

// vaultKVResponse models the Vault KV v2 GET response envelope.
type vaultKVResponse struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}

// doHTTPRequestValue executes a GET request and extracts the "value" field
// from the Vault KV v2 response payload.
func doHTTPRequestValue(req *http.Request) (string, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("secret not found")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var kvResp vaultKVResponse
	if err := json.NewDecoder(resp.Body).Decode(&kvResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	val, ok := kvResp.Data.Data["value"]
	if !ok {
		return "", fmt.Errorf("response missing 'value' field")
	}
	return val, nil
}
