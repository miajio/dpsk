package dpsk

import "time"

// Option 选项
type Option func(*Config)

// Config 配置
type Config struct {
	ApiKey             string        // api key
	BaseUrl            string        // base url
	ModelsUrl          string        // 分页查询模型列表的url 默认 /models
	BalanceUrl         string        // 获取账户余额的url 默认 /user/balance
	ChatCompletionsUrl string        // 对话请求的url 默认 /chat/completions
	Timeout            time.Duration // time out
}

// WithApiKey 设置api key
func WithApiKey(apiKey string) Option {
	return func(c *Config) {
		c.ApiKey = apiKey
	}
}

// WithBaseUrl 设置base url
func WithBaseUrl(baseUrl string) Option {
	return func(c *Config) {
		c.BaseUrl = baseUrl
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithModelsUrl 配置模型列表的url
func WithModelsUrl(modelsUrl string) Option {
	return func(c *Config) {
		c.ModelsUrl = modelsUrl
	}
}

// WithBalanceUrl 配置余额的url
func WithBalanceUrl(balanceUrl string) Option {
	return func(c *Config) {
		c.BalanceUrl = balanceUrl
	}
}

// WithChatCompletionsUrl 配置对话的url
func WithChatCompletionsUrl(chatCompletionsUrl string) Option {
	return func(c *Config) {
		c.ChatCompletionsUrl = chatCompletionsUrl
	}
}
