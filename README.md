# DeepSeek Go SDK

一个基于 Go 语言开发的 DeepSeek API 客户端 SDK，提供简单易用的接口与 DeepSeek 进行交互。

## 功能特性

✅ 支持 DeepSeek 所有模型查询  
✅ 支持账户余额查询  
✅ 支持普通聊天模式  
✅ 支持流式聊天模式  
✅ 支持上下文对话管理  
✅ 可配置超时时间  

## 安装

```bash
go get github.com/miajio/dpsk
```

## 快速开始

### 1. 初始化客户端

```go
import "github.com/miajio/dpsk/engine"

client, err := engine.NewClient(
    engine.WithApiKey("YOUR_DEEPSEEK_API_KEY"),
    engine.WithTimeout(30*time.Second),
)
if err != nil {
    log.Fatal("初始化客户端失败:", err)
}
```

### 2. 查询可用模型

```go
models, err := client.GetModels(context.Background())
if err != nil {
    log.Fatal("获取模型失败:", err)
}
fmt.Printf("可用模型: %+v\n", models)
```

### 3. 查询账户余额

```go
balance, err := client.GetBalance(context.Background())
if err != nil {
    log.Fatal("获取余额失败:", err)
}
fmt.Printf("账户余额: %+v\n", balance)
```

### 4. 普通聊天模式

```go
import "github.com/miajio/dpsk/chat"

req, err := chat.NewChatRequest(
    chat.WithModel("deepseek-chat"),
    chat.WithMessages([]chat.Message{
        {Role: "system", Content: "你是一个情感ai程序"},
        {Role: "user", Content: "你好,我叫小明"},
    }...),
)
if err != nil {
    log.Fatal("创建请求失败:", err)
}

resp, err := client.Chat(context.Background(), req)
if err != nil {
    log.Fatal("聊天失败:", err)
}
fmt.Printf("AI回复: %s\n", resp.Choices[0].Message.Content)
```

### 5. 流式聊天模式

```go
req, err := chat.NewChatRequest(
    chat.WithModel("deepseek-chat"),
    chat.WithMessages([]chat.Message{
        {Role: "system", Content: "你是一个情感ai程序"},
        {Role: "user", Content: "你好,我叫小明"},
    }...),
    chat.WithStream(true),
)

stream, streamErr, err := client.ChatStream(context.Background(), req)
if err != nil {
    log.Fatal("开启流式聊天失败:", err)
}

for {
    select {
    case msg, ok := <-stream:
        if !ok {
            return
        }
        if len(msg.Choices) > 0 {
            fmt.Print(msg.Choices[0].Delta.Content)
        }
    case err, ok := <-streamErr:
        if ok {
            log.Printf("流错误: %v", err)
        }
        return
    }
}
```

## 高级用法

### 上下文对话管理

```go
// 初始化对话
req, _ := chat.NewChatRequest(
    chat.WithModel("deepseek-chat"),
    chat.WithMessages([]chat.Message{
        {Role: "system", Content: "你是一个情感ai程序"},
        {Role: "user", Content: "你好,我叫小明"},
    }...),
)

// 获取AI回复后添加到上下文
resp, _ := client.Chat(context.Background(), req)
req.AddMessage("assistant", resp.Choices[0].Message.Content)

// 继续对话
req.AddMessage("user", "你喜欢听歌么?")
resp, _ = client.Chat(context.Background(), req)
```

## 配置选项

| 选项 | 描述 | 默认值 |
|------|------|--------|
| `WithApiKey` | 设置DeepSeek API密钥 | 必填 |
| `WithTimeout` | 设置请求超时时间 | 30秒 |
| `WithBaseUrl` | 设置API基础URL | DeepSeek官方API地址 |

## 贡献

欢迎提交 Pull Request 或 Issue 来改进这个项目！

## 许可证

MIT License
