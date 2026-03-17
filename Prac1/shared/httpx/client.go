package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client – обёртка над http.Client с поддержкой таймаута и прокидывания request-id.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// DoRequest выполняет HTTP-запрос, автоматически добавляя заголовки и request-id из контекста.
func (c *Client) DoRequest(ctx context.Context, method, path string, headers map[string]string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if body != nil {
		req.Header.Set("content-type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Прокидываем request-id из контекста, если он есть.
	if requestID, ok := ctx.Value("requestID").(string); ok && requestID != "" {
		req.Header.Set("x-request-id", requestID)
	}

	return c.httpClient.Do(req)
}
