package protocol

import "time"

// 기본적인 메시지 구조체
type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	TimeStamp time.Time `json:"timestamp"`
	ContextID string    `json:"context_id"` // 대화의 전체 ID
}

// 대화의 전체적인 흐름
type MessageContextFlow struct {
	ID        string    `json:"id"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ModelResponse struct {
	Message   Message `json:"message"`
	ContextID string  `json:"context_id"`
}

type ModelInterface interface {
	ProcessMessage(message Message, MCF MessageContextFlow) (ModelResponse, error)
}

type MCFImpl interface {
	CreateContext() (MessageContextFlow, error)
	GetContext(id string) (MessageContextFlow, error)
	UpdateContext(context MessageContextFlow) error
	DeleteContext(id string) error
	CurrentFlowsCount() int
}
