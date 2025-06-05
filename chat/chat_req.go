package chat

import (
	"net/http"

	"github.com/miajio/dpsk/errors"
)

// ChatReq 对话请求
type ChatRequest struct {
	Messages         []Message       `json:"messages"`                    // 消息历史列表,包含用户输入和AI回复的对话上下文
	Model            string          `json:"model"`                       // 指定使用的模型名称
	FrequencyPenalty float64         `json:"frequency_penalty,omitempty"` // 频率惩罚(-2.0到2.0),正值抑制重复内容
	MaxTokens        int             `json:"max_tokens,omitempty"`        // 生成的最大token数量(控制回复长度)
	PresencePenalty  float64         `json:"presence_penalty,omitempty"`  // 存在惩罚(-2.0到2.0),正值抑制已提及的内容
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`   // 输出格式
	Stop             any             `json:"stop,omitempty"`              // 停止词(string或[]string)
	Stream           bool            `json:"stream,omitempty"`            // 是否使用流式传输
	StreamOptions    *StreamOptions  `json:"stream_options,omitempty"`    // 流式选项
	Temperature      float64         `json:"temperature,omitempty"`       // 采样温度(0-2)
	TopP             float64         `json:"top_p,omitempty"`             // 核心采样(0-1)
	Tools            []Tool          `json:"tools,omitempty"`             // 可用工具列表
	ToolChoice       any             `json:"tool_choice,omitempty"`       // 工具选择策略
	Logprobs         bool            `json:"logprobs,omitempty"`          // 是否返回对数概率
	TopLogprobs      int             `json:"top_logprobs,omitempty"`      // 返回top N的对数概率(0-20)
}

// Validate 验证请求参数
func (cr *ChatRequest) Validate() error {
	if cr.Messages != nil && len(cr.Messages) > 0 {
		for _, msg := range cr.Messages {
			if err := msg.Validate(); err != nil {
				return err
			}
		}
	}

	if cr.Model == "" {
		return errors.NewCodeError(http.StatusBadRequest, "model is required")
	}

	return nil
}

// AddMessage 添加消息
func (cr *ChatRequest) AddMessage(role, content string) error {
	msg := Message{Role: role, Content: content}
	if err := msg.Validate(); err != nil {
		return err
	}
	cr.Messages = append(cr.Messages, msg)
	return nil
}

// AddMessageObj 添加消息对象
func (cr *ChatRequest) AddMessages(msgs ...Message) error {
	for _, msg := range msgs {
		if err := msg.Validate(); err != nil {
			return err
		}
	}
	cr.Messages = append(cr.Messages, msgs...)
	return nil
}

// MessageReq 对话消息请求
type Message struct {
	Role             string `json:"role"`                        // 角色
	Content          string `json:"content"`                     // 内容
	Name             string `json:"name,omitempty"`              // 名称
	Prefix           bool   `json:"prefix,omitempty"`            // 是否为前缀
	ReasoningContent string `json:"reasoning_content,omitempty"` // 对话前续写下用于assistant思维链内容输入
	ToolClassId      string `json:"tool_class_id,omitempty"`     // 此消息所响应的 tool call 的 ID
}

// Validate 验证消息参数
func (m *Message) Validate() error {
	if m.Role == "" {
		return errors.NewCodeError(http.StatusBadRequest, "role is required")
	}
	if m.Content == "" {
		return errors.NewCodeError(http.StatusBadRequest, "content is required")
	}
	return nil
}

// ResponseFormat 输出格式结构体
type ResponseFormat struct {
	Type string `json:"type,omitempty"` // 输出格式类型，如 "text" 或 "json_object"
}

// StreamOptions 流式选项结构体
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"` // 是否在流中包含使用情况
}

// Tool 工具结构体
type Tool struct {
	Type     string   `json:"type"`     // 工具类型(通常为"function")
	Function Function `json:"function"` // 函数定义
}

// Function 函数结构体
type Function struct {
	Name        string `json:"name"`        // 函数名称
	Description string `json:"description"` // 函数描述
	Parameters  any    `json:"parameters"`  // 函数参数(JSON Schema格式)
}
