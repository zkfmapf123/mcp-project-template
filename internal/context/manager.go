package context

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-mcp/pkg/protocol"
	"github.com/google/uuid"
)

type Manager struct {
	mcfContextManager map[string]protocol.MessageContextFlow
	mu                sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		mcfContextManager: map[string]protocol.MessageContextFlow{},
	}
}

func (m *Manager) CreateContext() (protocol.MessageContextFlow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	context := protocol.MessageContextFlow{
		ID:        generateID(),
		Messages:  make([]protocol.Message, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}

	m.mcfContextManager[context.ID] = context
	return context, nil
}

func (m *Manager) GetContext(id string) (protocol.MessageContextFlow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	content, isExists := m.mcfContextManager[id]
	if !isExists {
		return protocol.MessageContextFlow{}, fmt.Errorf("%s context is not found", id)
	}

	return content, nil
}

func (m *Manager) UpdateContext(mcf protocol.MessageContextFlow) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.mcfContextManager[mcf.ID]; !exists {
		return fmt.Errorf("%s context is not found", mcf.ID)
	}

	mcf.UpdatedAt = time.Now()
	m.mcfContextManager[mcf.ID] = mcf
	return nil
}

func (m *Manager) DeleteContext(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.mcfContextManager, id)
	return nil
}

func (m *Manager) CurrentFlowsCount() int {
	return len(m.mcfContextManager)
}

func generateID() string {
	return uuid.NewString()
}

func (m *Manager) GetMessages() map[string][]string {
	msg := map[string][]string{}

	for _, contextManager := range m.mcfContextManager {
		for _, context := range contextManager.Messages {
			msg[contextManager.ID] = append(msg[contextManager.ID], context.Content)
		}
	}

	return msg
}
