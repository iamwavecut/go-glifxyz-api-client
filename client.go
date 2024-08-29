package glifclient

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"log/slog"

	"golang.org/x/time/rate"
)

const (
	EndpointRun       = "/api/v1/run/"
	EndpointAddresses = "/api/v1/addresses"
	EndpointGlifs     = "/api/glifs"
	EndpointRuns      = "/api/runs"
	EndpointUser      = "/api/user"
	EndpointMe        = "/api/me"
	EndpointSpheres   = "/api/spheres"
	DefaultBaseURL    = "https://glif.app"
)

type GlifClient struct {
	BaseURL     string
	Client      *http.Client
	Logger      *slog.Logger
	RateLimiter *rate.Limiter
	APIToken    string
}

type AddressList struct {
	Addresses []string `json:"addresses"`
}

type RunResponse struct {
	ID     string      `json:"id"`
	Inputs interface{} `json:"inputs"`
	Output string      `json:"output"`
}

type GlifInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Add other fields as needed
}

type GlifRun struct {
	ID     string      `json:"id"`
	GlifID string      `json:"glifId"`
	Inputs interface{} `json:"inputs"`
	Output interface{} `json:"output"`
	// Add other fields as needed
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	// Add other fields as needed
}

type SphereInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	// Add other fields as needed
}

type ClientOption func(*GlifClient)

func WithBaseURL(url string) ClientOption {
	return func(c *GlifClient) {
		c.BaseURL = url
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *GlifClient) {
		c.Client = client
	}
}

func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *GlifClient) {
		c.Logger = logger
	}
}

func WithRateLimit(r rate.Limit, b int) ClientOption {
	return func(c *GlifClient) {
		c.RateLimiter = rate.NewLimiter(r, b)
	}
}

func WithAPIToken(token string) ClientOption {
	return func(c *GlifClient) {
		c.APIToken = token
	}
}

func NewGlifClient(opts ...ClientOption) *GlifClient {
	c := &GlifClient{
		BaseURL:     DefaultBaseURL,
		Client:      &http.Client{Timeout: 30 * time.Second},
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		RateLimiter: rate.NewLimiter(rate.Limit(10), 1), // Default rate limit: 10 requests per second
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *GlifClient) RunSimple(ctx context.Context, modelID string, args interface{}) (*RunResponse, error) {
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	endpoint := fmt.Sprintf("%s%s%s", c.BaseURL, EndpointRun, modelID)

	payload, err := json.Marshal(map[string]interface{}{
		"id":     modelID,
		"inputs": args,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	c.Logger.InfoContext(ctx, "Sending request", "endpoint", endpoint, "modelID", modelID)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var runResponse RunResponse
	if err := json.NewDecoder(resp.Body).Decode(&runResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &runResponse, nil
}

func (c *GlifClient) GetAddresses(ctx context.Context) (*AddressList, error) {
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	endpoint := fmt.Sprintf("%s%s", c.BaseURL, EndpointAddresses)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.Logger.InfoContext(ctx, "Sending request", "endpoint", endpoint)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var addresses AddressList
	if err := json.NewDecoder(resp.Body).Decode(&addresses); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &addresses, nil
}

func (c *GlifClient) GetGlifs(ctx context.Context, params url.Values) ([]GlifInfo, error) {
	endpoint := fmt.Sprintf("%s%s?%s", c.BaseURL, EndpointGlifs, params.Encode())
	var result []GlifInfo
	err := c.getJSON(ctx, endpoint, &result)
	return result, err
}

func (c *GlifClient) GetGlifRuns(ctx context.Context, glifID string, params url.Values) ([]GlifRun, error) {
	endpoint := fmt.Sprintf("%s%s?glifId=%s&%s", c.BaseURL, EndpointRuns, glifID, params.Encode())
	var result []GlifRun
	err := c.getJSON(ctx, endpoint, &result)
	return result, err
}

func (c *GlifClient) GetUserInfo(ctx context.Context, usernameOrID string) (*UserInfo, error) {
	endpoint := fmt.Sprintf("%s%s?%s=%s", c.BaseURL, EndpointUser, getIdentifierType(usernameOrID), usernameOrID)
	var result UserInfo
	err := c.getJSON(ctx, endpoint, &result)
	return &result, err
}

func (c *GlifClient) GetMyInfo(ctx context.Context) (*UserInfo, error) {
	endpoint := fmt.Sprintf("%s%s", c.BaseURL, EndpointMe)
	var result UserInfo
	err := c.getJSON(ctx, endpoint, &result)
	return &result, err
}

func (c *GlifClient) GetSpheres(ctx context.Context, params url.Values) ([]SphereInfo, error) {
	endpoint := fmt.Sprintf("%s%s?%s", c.BaseURL, EndpointSpheres, params.Encode())
	var result []SphereInfo
	err := c.getJSON(ctx, endpoint, &result)
	return result, err
}

func (c *GlifClient) StreamRunSimple(ctx context.Context, modelID string, args interface{}, callback func([]byte) error) error {
	endpoint := fmt.Sprintf("%s%s%s", c.BaseURL, EndpointRun, modelID)
	payload, err := json.Marshal(map[string]interface{}{
		"id":     modelID,
		"inputs": args,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading stream: %w", err)
		}

		if err := callback(line); err != nil {
			return fmt.Errorf("error in callback: %w", err)
		}
	}

	return nil
}

func (c *GlifClient) getJSON(ctx context.Context, endpoint string, v interface{}) error {
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit exceeded: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	c.Logger.InfoContext(ctx, "Sending request", "endpoint", endpoint)

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func getIdentifierType(identifier string) string {
	if strings.HasPrefix(identifier, "cl") {
		return "id"
	}
	return "username"
}
