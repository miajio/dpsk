package chat

// TopLogprobs 输出token对数概率信息
type TopLogprobs struct {
	Token   string `json:"token"`   // token
	Logprob int    `json:"logprob"` // 对数概率
	Bytes   []int  `json:"bytes"`   // 一个包含该 token UTF-8 字节表示的整数列表。一般在一个 UTF-8 字符被拆分成多个 token 来表示时有用。如果 token 没有对应的字节表示，则该值为 null
}

// Content 输出token对数概率信息
type Content struct {
	Token       string        `json:"token"`        // token
	Logprob     int           `json:"logprob"`      // 对token的对数概率
	Bytes       []int         `json:"bytes"`        // 一个包含该 token UTF-8 字节表示的整数列表。一般在一个 UTF-8 字符被拆分成多个 token 来表示时有用。如果 token 没有对应的字节表示，则该值为 null。
	TopLogprobs []TopLogprobs `json:"top_logprobs"` // 一个包含在该输出位置上，输出概率 top N 的 token 的列表，以及它们的对数概率。在罕见情况下，返回的 token 数量可能少于请求参数中指定的 top_logprobs 值。
}

type Logprobs struct {
	Content []Content `json:"content"` // 一个包含输出 token 对数概率信息的列表
}

// Choice 推理生成的补全选择列表
type Choice struct {
	Delta        Message  `json:"delta"`         // 输出的 delta 信息
	FinishReason string   `json:"finish_reason"` // 完成原因 [stop, length, content_filter, tool_calls, insufficient_system_resource]
	Index        int      `json:"index"`         // 输出长度达到模型上下文长度限制, 或达到了max_tokens的限制
	Message      Message  `json:"message"`       // 模型生成的补全信息
	Logprobs     Logprobs `json:"logprobs"`      // 对数概率信息
}

// CompletionTokensDetails 推理模型所产生的思维链 token 数量
type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

// usage 补全用量信息
type Usage struct {
	CompletionTokens        int                     `json:"completion_tokens"`         // 输出token数
	PromptTokens            int                     `json:"prompt_tokens"`             // 输入token数
	PromptCacheHitTokens    int                     `json:"prompt_cache_hit_tokens"`   // 上下文缓存token数
	PromptCacheMissTokens   int                     `json:"prompt_cache_miss_tokens"`  // 未命中上下文缓存的token数
	TotalTokens             int                     `json:"total_tokens"`              // 总token数
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"` // 补全使用的token详情
}

// ChatResponse API响应结构体
type ChatResponse struct {
	ID                string   `json:"id"`                 // 请求ID
	Choices           []Choice `json:"choices"`            // 聊天补全选项
	Created           int      `json:"created"`            // 创建时间戳
	Model             string   `json:"model"`              // 模型名称
	SystemFingerprint string   `json:"system_fingerprint"` // 系统指纹
	Object            string   `json:"object"`             // 对象类型
	Usage             Usage    `json:"usage"`              // 消耗的用量
}
