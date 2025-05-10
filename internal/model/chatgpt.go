package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-mcp/pkg/protocol"
	"github.com/sashabaranov/go-openai"
)

type ChatGPTModel struct {
	client *openai.Client
}

func NewChatGPTModel(apiKey string) *ChatGPTModel {
	return &ChatGPTModel{
		client: openai.NewClient(apiKey),
	}
}

var (
	SYSTEM    = openai.ChatMessageRoleSystem
	ASSISTANT = openai.ChatMessageRoleAssistant
	USER      = openai.ChatMessageRoleUser
)

func (m *ChatGPTModel) ProcessMessage(message protocol.Message, mcf protocol.MessageContextFlow) (protocol.ModelResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    SYSTEM,
			Content: "you are a helpful assistant",
		},
	}

	// 이전 메시지들을 추가
	for _, msg := range mcf.Messages {

		// 유저 답변
		role := USER

		// chat gpt 답변
		if strings.Contains(msg.ID, "response") {
			role = ASSISTANT
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 현재 메시지 추가
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    USER,
		Content: message.Content,
	})

	resp, err := m.client.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4Dot1Mini,
			Messages: messages,
		})
	if err != nil {
		return protocol.ModelResponse{}, fmt.Errorf("OpenAI API 호출 실패 : ", err)
	}

	response := protocol.Message{
		ID:        fmt.Sprintf("%s__response", message.ID),
		Content:   resp.Choices[0].Message.Content,
		TimeStamp: time.Now(),
		ContextID: mcf.ID,
	}

	return protocol.ModelResponse{
		Message:   response,
		ContextID: mcf.ID,
	}, nil
}
