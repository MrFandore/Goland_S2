package authclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Prac1/shared/httpx"
)

type Client struct {
	httpClient *httpx.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: httpx.NewClient(baseURL, timeout),
	}
}

type verifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject"`
	Error   string `json:"error"`
}

func (c *Client) Verify(ctx context.Context, token, requestID string) (bool, error) {
	headers := map[string]string{
		"authorization": "bearer " + token,
		"x-request-id":  requestID,
	}

	resp, err := c.httpClient.DoRequest(ctx, http.MethodGet, "/v1/auth/verify", headers, nil)
	if err != nil {
		return false, fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var vResp verifyResponse
		if err := json.NewDecoder(resp.Body).Decode(&vResp); err != nil {
			return false, fmt.Errorf("decode auth response: %w", err)
		}
		return vResp.Valid, nil
	case http.StatusUnauthorized:

		return false, nil
	default:
		return false, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}
}
