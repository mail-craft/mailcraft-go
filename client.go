// Package mailcraft is the official Go SDK for the MailCraft email API.
//
//	client := mailcraft.NewClient(os.Getenv("MAILCRAFT_API_KEY"))
//	_, err := client.Emails.Send(ctx, &mailcraft.SendEmailParams{
//		From:    "hello@yourdomain.com",
//		To:      []string{"person@example.com"},
//		Subject: "Welcome!",
//		HTML:    "<p>Thanks for signing up.</p>",
//	})
package mailcraft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Version is the SDK version, sent in the User-Agent header.
const Version = "1.0.0"

// DefaultBaseURL is the MailCraft API's base URL.
const DefaultBaseURL = "https://api.mailcraft.host/v1"

// Response is the API's JSON response.
type Response = map[string]any

// Client is the MailCraft API client. Create one with NewClient.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	Emails          *EmailsService
	Domains         *DomainsService
	Senders         *SendersService
	Contacts        *ContactsService
	Lists           *ListsService
	Segments        *SegmentsService
	Properties      *PropertiesService
	Templates       *TemplatesService
	TemplateFolders *TemplateFoldersService
	Campaigns       *CampaignsService
	Webhooks        *WebhooksService
	Suppressions    *SuppressionsService
	Metrics         *MetricsService
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the API base URL (staging, self-hosted).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithHTTPClient uses a custom *http.Client (timeouts, proxies, transports).
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

// NewClient creates a client. Create an API key under Settings > API keys.
// It panics if apiKey is empty, as the client can't make any request without it.
func NewClient(apiKey string, opts ...Option) *Client {
	if apiKey == "" {
		panic("mailcraft: an API key is required. Find yours under Settings > API Keys.")
	}

	c := &Client{
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Emails = &EmailsService{c}
	c.Domains = &DomainsService{c}
	c.Senders = &SendersService{c}
	c.Contacts = &ContactsService{c}
	c.Lists = &ListsService{c}
	c.Segments = &SegmentsService{c}
	c.Properties = &PropertiesService{c}
	c.Templates = &TemplatesService{c}
	c.TemplateFolders = &TemplateFoldersService{c}
	c.Campaigns = &CampaignsService{c}
	c.Webhooks = &WebhooksService{c}
	c.Suppressions = &SuppressionsService{c}
	c.Metrics = &MetricsService{c}

	return c
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any) (Response, error) {
	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("mailcraft: encoding request: %w", err)
		}
		payload = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("mailcraft: building request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "mailcraft-go/"+Version)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mailcraft: request failed: %w", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("mailcraft: reading response: %w", err)
	}

	if res.StatusCode >= 400 {
		return nil, newAPIError(res.StatusCode, raw)
	}

	if res.StatusCode == http.StatusNoContent || len(bytes.TrimSpace(raw)) == 0 {
		return nil, nil
	}

	var decoded Response
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, fmt.Errorf("mailcraft: decoding response: %w", err)
	}

	return decoded, nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values) (Response, error) {
	return c.do(ctx, http.MethodGet, path, query, nil)
}

func (c *Client) post(ctx context.Context, path string, body any) (Response, error) {
	return c.do(ctx, http.MethodPost, path, nil, body)
}

func (c *Client) patch(ctx context.Context, path string, body any) (Response, error) {
	return c.do(ctx, http.MethodPatch, path, nil, body)
}

func (c *Client) delete(ctx context.Context, path string) error {
	_, err := c.do(ctx, http.MethodDelete, path, nil, nil)
	return err
}
