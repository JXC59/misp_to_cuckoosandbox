package misp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ClientConfig struct {
	Endpoint      string `validate:"required,url"`
	SkipTLSVerify bool
	APIKey        string `validate:"required"`
}

type Client struct {
	cfg        *ClientConfig
	endpoint   *url.URL
	httpClient *http.Client
}

func NewClient(cfg *ClientConfig) (*Client, error) {
	httpClient := &http.Client{}
	if cfg.SkipTLSVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	endpoint, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %w", err)
	}
	return &Client{cfg: cfg, httpClient: httpClient, endpoint: endpoint}, nil
}

func (c *Client) doRequest(ctx context.Context,
	method string, apiPath string, reqBody io.Reader) (*http.Response, error) {
	endpoint := c.endpoint.JoinPath(apiPath)
	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.cfg.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

type commonErrorResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

func (c *Client) doRequestExpectOK(ctx context.Context, method string, apiPath string, reqBody io.Reader) (io.ReadCloser, error) {
	resp, err := c.doRequest(ctx, method, apiPath, reqBody)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp commonErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, errResp.Message)
	}

	return resp.Body, nil
}

type ServerVersionResponse struct {
	Version      string `json:"version"`
	PermSync     bool   `json:"perm_sync"`
	PermSighting bool   `json:"perm_sighting"`
}

func (c *Client) ServerVersion(ctx context.Context) (*ServerVersionResponse, error) {
	respBody, err := c.doRequestExpectOK(ctx, http.MethodGet, "/servers/getVersion", nil)
	if err != nil {
		return nil, err
	}

	var ret ServerVersionResponse
	if err := json.NewDecoder(respBody).Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
