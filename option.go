package dpsk

import "time"

// Option 选项
type Option func(*Config)

// Config 配置
type Config struct {
	ApiKey  string        // api key
	BaseUrl string        // base url
	Timeout time.Duration // time out
}

func WithApiKey(apiKey string) Option {
	return func(c *Config) {
		c.ApiKey = apiKey
	}
}

func WithBaseUrl(baseUrl string) Option {
	return func(c *Config) {
		c.BaseUrl = baseUrl
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}
