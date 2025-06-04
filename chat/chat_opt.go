package chat

// Option 配置项
type ChatOption func(*ChatRequest)

func NewChatRequest(options ...ChatOption) (*ChatRequest, error) {
	chatReq := &ChatRequest{
		Stream: false,
	}
	for _, option := range options {
		option(chatReq)
	}

	if err := chatReq.Validate(); err != nil {
		return nil, err
	}

	return chatReq, nil
}

// WithMessages 设置消息
func WithMessages(messages ...Message) ChatOption {
	return func(req *ChatRequest) {
		req.Messages = messages
	}
}

// WithModel 设置模型
func WithModel(model string) ChatOption {
	return func(req *ChatRequest) {
		req.Model = model
	}
}

// WithFrequencyPenalty 设置频率惩罚
func WithFrequencyPenalty(frequencyPenalty float64) ChatOption {
	return func(req *ChatRequest) {
		req.FrequencyPenalty = frequencyPenalty
	}
}

// WithMaxTokens 设置最大token数
func WithMaxTokens(maxTokens int) ChatOption {
	return func(req *ChatRequest) {
		req.MaxTokens = maxTokens
	}
}

// WithPresencePenalty 设置存在惩罚(-2.0到2.0),正值抑制已提及的内容
func WithPresencePenalty(presencePenalty float64) ChatOption {
	return func(req *ChatRequest) {
		req.PresencePenalty = presencePenalty
	}
}

// WithResponseFormat 设置响应格式
func WithResponseFormat(tp string) ChatOption {
	return func(req *ChatRequest) {
		req.ResponseFormat = &ResponseFormat{
			Type: tp,
		}
	}
}

// WithStop 设置停止符
func WithStop(stopVal ...string) ChatOption {
	return func(req *ChatRequest) {
		if stopVal != nil {
			if len(stopVal) == 1 {
				req.Stop = stopVal[0]
			} else {
				req.Stop = stopVal
			}
		}
	}
}

// WithStream 设置流式响应
func WithStream(stream bool) ChatOption {
	return func(req *ChatRequest) {
		req.Stream = stream
	}
}

// WithStreamOptions 设置流式响应的选项
func WithStreamOptions(includeUsage bool) ChatOption {
	return func(req *ChatRequest) {
		req.StreamOptions = &StreamOptions{
			IncludeUsage: includeUsage,
		}
	}
}

// WithTemperature 设置采样温度(0-2)
func WithTemperature(temperature float64) ChatOption {
	return func(req *ChatRequest) {
		if temperature > 0 && 2 > temperature {
			req.Temperature = temperature
		}
	}
}

// WithTopP 设置核心采样(0-1)
func WithTopP(topP float64) ChatOption {
	return func(req *ChatRequest) {
		if topP > 0 && 1 > topP {
			req.TopP = topP
		}
	}
}

// WithTools 设置工具
func WithTools(tools ...Tool) ChatOption {
	return func(req *ChatRequest) {
		req.Tools = tools
	}
}

// WithToolChoice 设置工具选择策略
func WithToolChoice(toolChoice any) ChatOption {
	return func(req *ChatRequest) {
		req.ToolChoice = toolChoice
	}
}

// WithLogprobs 设置是否返回对数概率
func WithLogprobs(logprobs bool) ChatOption {
	return func(req *ChatRequest) {
		req.Logprobs = logprobs
	}
}

// WithTopLogprobs 设置返回top N的对数概率(0-20)
func WithTopLogprobs(topLogprobs int) ChatOption {
	return func(req *ChatRequest) {
		if topLogprobs > 0 && 20 > topLogprobs {
			req.TopLogprobs = topLogprobs
		}
	}
}
