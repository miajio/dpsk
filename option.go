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

func WithModelsUrl(modelsUrl string) Option {
	return func(c *Config) {
		c.ModelsUrl = modelsUrl
	}
}

func WithBalanceUrl(balanceUrl string) Option {
	return func(c *Config) {
		c.BalanceUrl = balanceUrl
	}
}

func WithChatCompletionsUrl(chatCompletionsUrl string) Option {
	return func(c *Config) {
		c.ChatCompletionsUrl = chatCompletionsUrl
	}
}
