package engine

import (
	"net/http"
	"time"
)

type Option func(*Client)

// WithApiUrl 设置deepseek api url
func WithApiUrl(apiUrl string) Option {
	return func(c *Client) {
		c.apiUrl = apiUrl
	}
}

// WithModelsUrl 设置模型列表url
func WithModelsUrl(modelsUrl string) Option {
	return func(c *Client) {
		c.urlMap["models"] = modelsUrl
	}
}

// WithBalanceUrl 设置余额查询url
func WithBalanceUrl(balanceUrl string) Option {
	return func(c *Client) {
		c.urlMap["balance"] = balanceUrl
	}
}

// WithChatCompletionsUrl 设置对话url
func WithChatCompletionsUrl(chatCompletionsUrl string) Option {
	return func(c *Client) {
		c.urlMap["chat"] = chatCompletionsUrl
	}
}

// WithApiKey 设置apiKey
func WithApiKey(apiKey string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHttpClient 配置httpClient
func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
