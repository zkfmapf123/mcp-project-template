package model

import (
	"time"

	"github.com/go-mcp/pkg/protocol"
)

type SimpleModel struct {
	apiKey string
}

func NewSimpleModel(apiKey string) *SimpleModel {
	return &SimpleModel{
		apiKey: apiKey,
	}
}

func (m *SimpleModel) ProcessMessage(message protocol.Message, mcf protocol.MessageContextFlow) (protocol.ModelResponse, error) {
	// 여기에 메시지 처리 로직 구현
	response := protocol.ModelResponse{
		Message: protocol.Message{
			ID:        message.ID,
			Content:   "Simple model response >> " + message.Content,
			TimeStamp: time.Now(),
			ContextID: mcf.ID,
		},
		ContextID: mcf.ID,
	}

	return response, nil
}
