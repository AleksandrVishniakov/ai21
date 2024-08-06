package ai21

import "net/http"

type ClientConfigs struct {
	// Default: "https://api.ai21.com/studio/v1/"
	BaseURL string

	HTTPClient *http.Client
}

func DefaultConfigs() *ClientConfigs {
	return &ClientConfigs{
		BaseURL:    "https://api.ai21.com/studio/v1/",
		HTTPClient: &http.Client{},
	}
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return NewClientWithConfigs(apiKey, DefaultConfigs())
}

func NewClientWithConfigs(apiKey string, cfg *ClientConfigs) *Client {
	client := cfg.HTTPClient
	client.Transport = &authTransport{
		transport: cfg.HTTPClient.Transport,
		token:     apiKey,
	}
	return &Client{
		baseURL:    cfg.BaseURL,
		httpClient: client,
	}
}

type authTransport struct {
	transport http.RoundTripper
	token     string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)

	if t.transport == nil {
		t.transport = http.DefaultTransport
	}
	return t.transport.RoundTrip(req)
}
