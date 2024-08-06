package ai21

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	completionURL = "/chat/completions"
)

type CompletionMessage struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model          AIModel              `json:"model"`
	Messages       []*CompletionMessage `json:"messages"`
	MaxTokens      int                  `json:"max_tokens,omitempty"`
	Temperature    float32              `json:"temperature,omitempty"`
	TopP           float32              `json:"top_p,omitempty"`
	Stop           []string             `json:"stop,omitempty"`
	ResponsesCount int                  `json:"n,omitempty"`
	Stream         bool                 `json:"stream,omitempty"`
}

func defaultCompletionRequest() *CompletionRequest {
	return &CompletionRequest{
		Model:          ModelJambaInstructPreview,
		Messages:       make([]*CompletionMessage, 0),
		MaxTokens:      4096,
		Temperature:    1.0,
		TopP:           1.0,
		Stop:           make([]string, 0),
		ResponsesCount: 1,
		Stream:         false,
	}
}

type CompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Index        int               `json:"index"`
		Message      CompletionMessage `json:"message"`
		FinishReason FinishReason      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (r *CompletionResponse) Content() string {
	if len(r.Choices) == 0 {
		return ""
	}

	return r.Choices[0].Message.Content
}

type CompletionRequestOption func(r *CompletionRequest)

func WithAIModel(model AIModel) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.Model = model
	}
}

func WithInitialMessage(prompt string) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.Messages = append([]*CompletionMessage{{
			Role:    RoleSystem,
			Content: prompt,
		}}, r.Messages...)
	}
}

func WithMaxTokens(maxTokens int) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.MaxTokens = maxTokens
	}
}

func WithTemperature(temperature float32) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.Temperature = temperature
	}
}

func WithTopP(topP float32) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.TopP = topP
	}
}

func WithStopWords(stopWords []string) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.Stop = stopWords
	}
}

func WithResponsesCount(n int) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.ResponsesCount = n
	}
}

func WithStream() CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.Stream = true
	}
}

func (c *Client) CompletionRequest(
	ctx context.Context,
	prompt string,
	opts ...CompletionRequestOption,
) (*CompletionResponse, error) {
	request := defaultCompletionRequest()
	for _, optFunc := range opts {
		optFunc(request)
	}

	request.Messages = append(request.Messages, &CompletionMessage{RoleUser, prompt})

	return sendCompletionRequest(
		ctx,
		c.httpClient,
		c.baseURL,
		request,
	)
}

type Conversation struct {
	client      *Client
	history     []*CompletionMessage
	totalTokens int

	defRequest CompletionRequest
}

func NewConversation(client *Client, opts ...CompletionRequestOption) *Conversation {
	request := *defaultCompletionRequest()
	for _, optFunc := range opts {
		optFunc(&request)
	}

	request.ResponsesCount = 1

	return &Conversation{
		client:      client,
		history:     request.Messages,
		totalTokens: 0,
		defRequest:  request,
	}
}

func (c *Conversation) CompletionRequest(
	ctx context.Context,
	prompt string,
) (*CompletionResponse, error) {
	request := c.defRequest
	request.Messages = append(request.Messages, &CompletionMessage{RoleUser, prompt})

	response, err := sendCompletionRequest(
		ctx,
		c.client.httpClient,
		c.client.baseURL,
		&request,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send completion request: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response received: len(response.Choices) == 0")
	}

	c.history = append(
		c.history,
		&CompletionMessage{RoleUser, prompt},
		&response.Choices[0].Message,
	)
	
	c.totalTokens += response.Usage.TotalTokens

	return response, nil
}

func (c *Conversation) TotalTokens() int {
	return c.totalTokens
}

func sendCompletionRequest(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	request *CompletionRequest,
) (*CompletionResponse, error) {
	reqURL, err := url.JoinPath(baseURL, completionURL)
	if err != nil {
		return nil, fmt.Errorf("url.JoinPath(%q, %q): %w", baseURL, completionURL, err)
	}

	reqBody, err := encodeJSONToReader(request)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request to JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext(ctx, %q, %q, reqBody): %w", http.MethodPost, reqURL, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.httpClient.Do(req): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read HTTP error message with status %s from response: %w", resp.Status, err)
		}

		return nil, &APIError{
			Code:    resp.StatusCode,
			Message: string(body),
			URL:     reqURL,
			Method:  http.MethodPost,
		}
	}

	response, err := decodeJSON[CompletionResponse](resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &response, nil
}
