package demo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/miajio/dpsk"
)

const (
	API_KEY = "YOU DEEPSEEK API KEY"
)

func TestGetModel(t *testing.T) {
	client, err := dpsk.NewClient(dpsk.WithApiKey(API_KEY), dpsk.WithTimeout(30*time.Second))
	if err != nil {
		fmt.Println("failed to create client:", err)
		return
	}
	modelList, err := client.GetModels(context.Background())
	if err != nil {
		fmt.Println("failed to get models:", err)
		return
	}
	byteModelList, _ := json.Marshal(modelList)
	fmt.Println(string(byteModelList))

}

func TestGetBalance(t *testing.T) {
	client, err := dpsk.NewClient(dpsk.WithApiKey(API_KEY), dpsk.WithTimeout(30*time.Second))
	if err != nil {
		fmt.Println("failed to create client:", err)
		return
	}
	balance, err := client.GetBalance(context.Background())
	if err != nil {
		fmt.Println("failed to get models:", err)
		return
	}
	byteBalance, _ := json.Marshal(balance)
	fmt.Println(string(byteBalance))
}

func TestChat(t *testing.T) {
	client, err := dpsk.NewClient(dpsk.WithApiKey(API_KEY), dpsk.WithTimeout(24*time.Hour))
	if err != nil {
		fmt.Println("failed to create client:", err)
		return
	}

	chatReq := dpsk.ChatRequest{
		Model: "deepseek-chat",
		Messages: []dpsk.Message{
			{Role: "system", Content: "你是一个超级人工智能, 擅长于计算机网络领域"},
			{Role: "user", Content: "你好, 我是一名go语言程序员, 你能帮我实现一个基于go语言封装的deepseek的sdk吗? 我需要简单易用, deepseek的api网站是https://api-docs.deepseek.com"},
		},
		Stream: false,
	}

	res, err := client.Chat(context.Background(), chatReq)
	if err != nil {
		fmt.Println("failed to chat:", err)
		return
	}
	byteRes, _ := json.Marshal(res)
	fmt.Println(string(byteRes))
}

func TestChatStream(t *testing.T) {
	client, err := dpsk.NewClient(dpsk.WithApiKey(API_KEY), dpsk.WithTimeout(24*time.Hour))
	if err != nil {
		fmt.Println("failed to create client:", err)
		return
	}

	defaultMessage := []dpsk.Message{
		{Role: "system", Content: "你是一个情感ai程序"},
		{Role: "user", Content: "你好,我叫小明"},
	}

	chatReq := dpsk.ChatRequest{
		Model:    "deepseek-chat",
		Messages: defaultMessage,
		Stream:   true,
	}

	chatStream, chatStreamErr, err := client.ChatStream(context.Background(), chatReq)
	if err != nil {
		fmt.Println("failed to chat stream:", err)
		return
	}

	nextContent := ""

	for {
		select {
		case msg, ok := <-chatStream:
			if !ok {
				goto NextA
			}
			if len(msg.Choices) > 0 {
				delta := msg.Choices[0].Delta.Content
				nextContent += delta
				fmt.Print(delta)
			}
		case err, ok := <-chatStreamErr:
			if ok {
				log.Printf("Stream error: %v", err)
			}
			log.Println("Stream ended", err)
			goto NextA
		}
	}
NextA:
	defaultMessage = append(defaultMessage, dpsk.Message{Role: "assistant", Content: nextContent})
	defaultMessage = append(defaultMessage, dpsk.Message{Role: "user", Content: "你喜欢听歌么?我最近很emo,想听听让我开心的歌"})
	chatReq.Messages = defaultMessage

	chatStream, chatStreamErr, err = client.ChatStream(context.Background(), chatReq)
	if err != nil {
		fmt.Println("failed to chat stream:", err)
		return
	}

	nextContent = ""

	for {
		select {
		case msg, ok := <-chatStream:
			if !ok {
				goto NextB
			}
			if len(msg.Choices) > 0 {
				delta := msg.Choices[0].Delta.Content
				nextContent += delta
				fmt.Print(delta)
			}
		case err, ok := <-chatStreamErr:
			if ok {
				log.Printf("Stream error: %v", err)
			}
			log.Println("Stream ended", err)
			goto NextB
		}
	}
NextB:
	defaultMessage = append(defaultMessage, dpsk.Message{Role: "assistant", Content: nextContent})
	defaultMessage = append(defaultMessage, dpsk.Message{Role: "user", Content: "我是和女朋友吵架了, 她因为我在六一儿童节没有给她送礼物而不开心,我该怎么办?"})
	chatReq.Messages = defaultMessage

	chatStream, chatStreamErr, err = client.ChatStream(context.Background(), chatReq)
	if err != nil {
		fmt.Println("failed to chat stream:", err)
		return
	}

	nextContent = ""

	for {
		select {
		case msg, ok := <-chatStream:
			if !ok {
				return
			}
			if len(msg.Choices) > 0 {
				delta := msg.Choices[0].Delta.Content
				nextContent += delta
				fmt.Print(delta)
			}
		case err, ok := <-chatStreamErr:
			if ok {
				log.Printf("Stream error: %v", err)
			}
			log.Println("Stream ended", err)
			return
		}
	}
}
