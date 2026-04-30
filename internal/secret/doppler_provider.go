package secret

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const dopplerBaseURL = "https://api.doppler.com/v3"

type dopplerProvider struct {
	token   string
	project string
	config  string
	client  *http.Client
}

// NewDopplerProvider creates a secret provider backed by Doppler.
func NewDopplerProvider(token, project, config string) Provider {
	return &dopplerProvider{
		token:   token,
		project: project,
		config:  config,
		client:  &http.Client{},
	}
}

func (d *dopplerProvider) Get(key string) (string, error) {
	url := fmt.Sprintf("%s/configs/config/secret?project=%s&config=%s&name=%s",
		dopplerBaseURL, d.project, d.config, key)

	req, err := newHTTPRequest(http.MethodGet, url, d.token, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		Secret struct {
			RawValue struct {
				Raw string `json:"raw"`
			} `json:"raw_value"`
		} `json:"secret"`
	}

	if err := doHTTPRequest(d.client, req, &result); err != nil {
		return "", ErrNotFound{Key: key}
	}
	return result.Secret.RawValue.Raw, nil
}

func (d *dopplerProvider) Set(key, value string) error {
	url := fmt.Sprintf("%s/configs/config/secrets", dopplerBaseURL)
	body := map[string]interface{}{
		"project": d.project,
		"config":  d.config,
		"secrets": map[string]string{key: value},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := newHTTPRequest(http.MethodPost, url, d.token, b)
	if err != nil {
		return err
	}

	return doHTTPRequest(d.client, req, nil)
}

func (d *dopplerProvider) Delete(key string) error {
	url := fmt.Sprintf("%s/configs/config/secrets", dopplerBaseURL)
	body := map[string]interface{}{
		"project": d.project,
		"config":  d.config,
		"secrets": map[string]string{key: ""},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := newHTTPRequest(http.MethodDelete, url, d.token, b)
	if err != nil {
		return err
	}

	return doHTTPRequest(d.client, req, nil)
}
