package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/miajio/dpsk/errors"
)

const (
	apiUrl = "https://api.deepseek.com"
)

var (
	defaultUrlMap = map[string]string{
		"models":  "/models",
		"balance": "/user/balance",
		"chat":    "/chat/completions",
	}
)

type Client struct {
	httpClient *http.Client      // http客户端
	apiUrl     string            // api接口默认调用地址
	urlMap     map[string]string // urlMap
	apiKey     string            // apiKey
}

// NewClient 创建一个client
func NewClient(options ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{},
		apiUrl:     apiUrl,
		urlMap:     defaultUrlMap,
	}
	for _, option := range options {
		option(c)
	}
	return c, nil
}

// WithApiKey 设置apiKey
func (c *Client) WithApiKey(apiKey string) *Client {
	c.apiKey = apiKey
	return c
}

// WithHttpClient 设置httpClient
func (c *Client) WithHttpClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

// WithModelsUrl 设置模型列表url
func (c *Client) WithModelsUrl(modelsUrl string) *Client {
	c.urlMap["models"] = modelsUrl
	return c
}

// WithBalanceUrl 设置余额查询url
func (c *Client) WithBalanceUrl(balanceUrl string) *Client {
	c.urlMap["balance"] = balanceUrl
	return c
}

// WithChatCompletionsUrl 设置对话url
func (c *Client) WithChatCompletionsUrl(chatCompletionsUrl string) *Client {
	c.urlMap["chat"] = chatCompletionsUrl
	return c
}

// WithTimeout 设置超时时间
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

// makeRequest 创建请求
func (c *Client) makeRequest(ctx context.Context, method string, url string, body any) (*http.Response, error) {
	if c.apiKey == "" {
		return nil, errors.NewCodeError(http.StatusBadRequest, "api key is required")
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.httpClient.Do(req)
}
