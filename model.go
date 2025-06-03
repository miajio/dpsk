package dpsk

// Model 模型结构体
type Model struct {
	ID      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

//  ModelList 模型列表结构体
type ModelList struct {
	Data   []Model `json:"data,omitempty"`
	Object string  `json:"object,omitempty"`
}

// Balance 余额结构体
type Balance struct {
	IsAvailable  bool          `json:"is_available,omitempty"`
	BalanceInfos []BalanceInfo `json:"balance_infos,omitempty"`
}

//  BalanceInfo 余额信息结构体
type BalanceInfo struct {
	Currency        string `json:"currency,omitempty"`          // 货币,人民币或美元
	TotalBalance    string `json:"total_balance,omitempty"`     // 总的可用余额,包括赠金和充值余额
	GrantedBalance  string `json:"granted_balance,omitempty"`   // 未过期的赠金余额
	ToppedUpBalance string `json:"topped_up_balance,omitempty"` // 充值余额
}

// Function 方法结构体
type Function struct {
	Name      string `json:"name,omitempty"`      // 模型调用的方法名
	Arguments string `json:"arguments,omitempty"` // 调用的参数
}

// ToolCall 工具调用结构体
type ToolCall struct {
	ID       string   `json:"id,omitempty"`       // id
	Type     string   `json:"type,omitempty"`     //  工具类型,目前仅支持function
	Function Function `json:"function,omitempty"` // 调用方法
}

// Message 消息结构体
type Message struct {
	Content          string     `json:"content"`                     // 内容
	Role             string     `json:"role"`                        // 角色
	Name             string     `json:"name,omitempty"`              // 名称
	Prefix           bool       `json:"prefix,omitempty"`            // 是否为前缀
	ReasoningContent string     `json:"reasoning_content,omitempty"` // 对话前续写功能下作为assistant思维链内容的输入
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`        // 工具调用
	ToolCallId       string     `json:"tool_call_id,omitempty"`      // 此消息所响应的tool call id
}

// ResponseFormat 响应格式
// 一个 object,指定模型必须输出的格式.设置为 { "type": "json_object" } 以启用 JSON 模式,该模式保证模型生成的消息是有效的 JSON.
// 注意: 使用 JSON 模式时,你还必须通过系统或用户消息指示模型生成 JSON.
// 否则,模型可能会生成不断的空白字符,直到生成达到令牌限制,从而导致请求长时间运行并显得“卡住”.
// 此外,如果 finish_reason="length",这表示生成超过了 max_tokens 或对话超过了最大上下文长度,消息内容可能会被部分截断.
type ResponseFormat struct {
	Type string `json:"type,omitempty"`
}

// StreamOptions 对话流式请求结构体
// 流式输出相关选项.只有在 stream 参数为 true 时,才可设置此参数
// 如果设置为 true,在流式消息最后的 data: [DONE] 之前将会传输一个额外的块.
// 此块上的 usage 字段显示整个请求的 token 使用统计信息,而 choices 字段将始终是一个空数组.
// 所有其他块也将包含一个 usage 字段,但其值为 null.
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}

// ToolFunction 工具信息
type ToolFunction struct {
	Description string         `json:"description"`          // 描述
	Name        string         `json:"name"`                 // 工具名称
	Parameters  map[string]any `json:"parameters,omitempty"` // 参数信息
}

// Tool 功能工具结构体
type Tool struct {
	Type     string       `json:"type"`     // 工具类型, 目前仅支持function
	Function ToolFunction `json:"function"` // 工具信息
}

// ChatRequest 对话请求结构体
type ChatRequest struct {
	Messages []Message `json:"messages"` // 消息
	Model    string    `json:"model"`    // 使用的模型
	// 介于 -2.0 和 2.0 之间的数字.如果该值为正, 那么新 token 会根据其在已有文本中的出现频率受到相应的惩罚,降低模型重复相同内容的可能性.
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	// 介于 1 到 8192 间的整数,限制一次请求中模型生成 completion 的最大 token 数.输入 token 和输出 token 的总长度受模型的上下文长度的限制.如未指定 max_tokens参数,默认使用 4096.
	MaxTokens int `json:"max_tokens,omitempty"`
	// 介于 -2.0 和 2.0 之间的数字.如果该值为正,那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚,从而增加模型谈论新主题的可能性.
	PresencePenalty float64         `json:"presence_penalty,omitempty"`
	ResponseFormat  *ResponseFormat `json:"response_format,omitempty"` // 输出格式
	Stop            any             `json:"stop,omitempty"`            // 一个 string 或最多包含 16 个 string 的 list,在遇到这些词时,API 将停止生成更多的 token.
	Stream          bool            `json:"stream,omitempty"`          // 如果设置为 True,将会以 SSE（server-sent events）的形式以流式发送消息增量.消息流以 data: [DONE] 结尾.
	StreamOptions   *StreamOptions  `json:"stream_options,omitempty"`  // 对话流式请求配置
	Temperature     float64         `json:"temperature,omitempty"`     // 采样温度,介于 0 和 2 之间.更高的值,如 0.8,会使输出更随机,而更低的值,如 0.2,会使其更加集中和确定. 我们通常建议可以更改这个值或者更改 top_p,但不建议同时对两者进行修改.
	TopP            float64         `json:"top_p,omitempty"`           // 作为调节采样温度的替代方案,模型会考虑前 top_p 概率的 token 的结果.所以 0.1 就意味着只有包括在最高 10% 概率中的 token 会被考虑. 我们通常建议修改这个值或者更改 temperature,但不建议同时对两者进行修改.
	Tools           []Tool          `json:"tools,omitempty"`
	ToolChoice      any             `json:"tool_choice,omitempty"`
	Logprobs        bool            `json:"logprobs,omitempty"`     // 是否返回所输出 token 的对数概率.如果为 true,则在 message 的 content 中返回每个输出 token 的对数概率.
	TopLogprobs     int             `json:"top_logprobs,omitempty"` // 一个介于 0 到 20 之间的整数 N,指定每个输出位置返回输出概率 top N 的 token,且返回这些 token 的对数概率.指定此参数时,logprobs 必须为 true.
}

// TopLogprobs 输出token对数概率信息
type TopLogprobs struct {
	Token   string `json:"token,omitempty"`   // token
	Logprob int    `json:"logprob,omitempty"` // 对数概率
	Bytes   []int  `json:"bytes,omitempty"`   // 一个包含该 token UTF-8 字节表示的整数列表。一般在一个 UTF-8 字符被拆分成多个 token 来表示时有用。如果 token 没有对应的字节表示，则该值为 null
}

// Content 输出token对数概率信息
type Content struct {
	Token       string        `json:"token,omitempty"`        // token
	Logprob     int           `json:"logprob,omitempty"`      // 对token的对数概率
	Bytes       []int         `json:"bytes,omitempty"`        // 一个包含该 token UTF-8 字节表示的整数列表。一般在一个 UTF-8 字符被拆分成多个 token 来表示时有用。如果 token 没有对应的字节表示，则该值为 null。
	TopLogprobs []TopLogprobs `json:"top_logprobs,omitempty"` // 一个包含在该输出位置上，输出概率 top N 的 token 的列表，以及它们的对数概率。在罕见情况下，返回的 token 数量可能少于请求参数中指定的 top_logprobs 值。
}

type Logprobs struct {
	Content []Content `json:"content,omitempty"` // 一个包含输出 token 对数概率信息的列表
}

// Choice 推理生成的补全选择列表
type Choice struct {
	Delta        Message  `json:"delta,omitempty"`         //  输出的 delta 信息
	FinishReason string   `json:"finish_reason,omitempty"` // 结束原因 [stop, length, content_filter, tool_calls, insufficient_system_resource]
	Index        int      `json:"index,omitempty"`         // 输出长度达到模型上下文长度限制, 或达到了max_tokens的限制
	Message      Message  `json:"message,omitempty"`       // 模型生成的补全信息
	Logprobs     Logprobs `json:"logprobs,omitempty"`      // 对 choice的对数概率信息
}

// CompletionTokensDetails 推理模型所产生的思维链 token 数量
type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
}

// usage 补全用量信息
type Usage struct {
	CompletionTokens        int                     `json:"completion_tokens,omitempty"`         // 模型补全消耗的token数
	PromptTokens            int                     `json:"prompt_tokens,omitempty"`             // 包含token数 =
	PromptCacheHitTokens    int                     `json:"prompt_cache_hit_tokens,omitempty"`   // 上下文缓存token数
	PromptCacheMissTokens   int                     `json:"prompt_cache_miss_tokens,omitempty"`  // 未命中上下文缓存的token数
	TotalTokens             int                     `json:"total_tokens,omitempty"`              // 总token数
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details,omitempty"` // 补全使用的token详情
}

// ChatCompletion 聊天补全结构体
type ChatCompletion struct {
	ID                string   `json:"id,omitempty"`                 // id
	Choices           []Choice `json:"choices,omitempty"`            // 聊天补全选项
	Created           int      `json:"created,omitempty"`            // 创建聊天完成的时间戳(秒)
	Model             string   `json:"model,omitempty"`              // 使用的模型
	SystemFingerprint string   `json:"system_fingerprint,omitempty"` // 系统指纹
	Object            string   `json:"object,omitempty"`             // 对象 值为:chat.completion
	Usage             Usage    `json:"usage,omitempty"`              // 消耗的用量
}
