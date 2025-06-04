package chat

// Option 配置项
type ChatOption func(*ChatRequest)

var (
	// 创意写作配置
	WithCreativeWritingOption = func() ChatOption {
		return func(req *ChatRequest) {
			req.Temperature = 1.2
			req.TopP = 0.95
			req.FrequencyPenalty = 0.3
			req.PresencePenalty = 0.5
		}
	}

	// 技术问答配置
	WithTechnicalQuestionOption = func() ChatOption {
		return func(req *ChatRequest) {
			req.Temperature = 0.3
			req.TopP = 0.7
			req.FrequencyPenalty = 0.7
			req.PresencePenalty = 0.3
		}
	}

	// 多轮对话配置
	WithMultipleDialogueOption = func() ChatOption {
		return func(req *ChatRequest) {
			req.Temperature = 0.8
			req.TopP = 0.9
			req.FrequencyPenalty = 0.4
			req.PresencePenalty = 0.2
		}
	}
)

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
