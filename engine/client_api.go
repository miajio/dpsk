package engine

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/miajio/dpsk/chat"
	"github.com/miajio/dpsk/errors"
	"github.com/miajio/dpsk/model"
)

// GetModels 获取模型列表
func (c *Client) GetModels(ctx context.Context) (*model.ModelList, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, c.apiUrl+c.urlMap["models"], nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewCodeErrorF(resp.StatusCode, "failed to get models: %s", resp.Status)
	}

	var modelList model.ModelList
	if err := json.NewDecoder(resp.Body).Decode(&modelList); err != nil {
		return nil, err
	}
	return &modelList, nil
}

// GetBalance 获取账户余额
func (c *Client) GetBalance(ctx context.Context) (*model.Balance, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, c.apiUrl+c.urlMap["balance"], nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewCodeErrorF(resp.StatusCode, "failed to get balance: %s", resp.Status)
	}

	var balance model.Balance
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, err
	}
	return &balance, nil
}

// Chat 发送消息到模型
func (c *Client) Chat(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	if req.Stream {
		return nil, errors.NewCodeError(http.StatusBadRequest, "streaming is not supported, use ChatStream instead")
	}
	resp, err := c.makeRequest(ctx, http.MethodPost, c.apiUrl+c.urlMap["chat"], req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewCodeErrorF(resp.StatusCode, "failed to chat: %s", resp.Status)
	}
	var completion chat.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return nil, err
	}
	return &completion, nil
}

// ChatStream 发送流式请求
func (c *Client) ChatStream(ctx context.Context, req *chat.ChatRequest) (<-chan chat.ChatResponse, <-chan error, error) {
	if !req.Stream {
		return nil, nil, errors.NewCodeError(http.StatusBadRequest, "stream is not enabled")
	}

	resp, err := c.makeRequest(ctx, http.MethodPost, c.apiUrl+c.urlMap["chat"], req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.NewCodeErrorF(resp.StatusCode, "failed to chat stream: %s", resp.Status)
	}

	errChan := make(chan error, 1)
	resChain := make(chan chat.ChatResponse)

	go func() {
		defer resp.Body.Close()
		defer close(resChain)
		defer close(errChan)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}
			if line == "data: [DONE]" {
				return
			}
			if !strings.HasPrefix(line, "data: ") {
				errChan <- errors.NewCodeErrorF(http.StatusBadRequest, "invalid response: %s", line)
				continue
			}

			jsonData := strings.TrimPrefix(line, "data: ")
			var event chat.ChatResponse
			if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
				log.Printf("failed to parse response: %s, error: %v", jsonData, err)
				errChan <- errors.NewCodeErrorF(http.StatusBadRequest, "failed to parse response: %v", err)
				continue
			}
			resChain <- event
		}
		if err := scanner.Err(); err != nil {
			errChan <- errors.NewCodeErrorF(http.StatusBadRequest, "scanner error: %v", err)
		}
	}()
	return resChain, errChan, nil
}
