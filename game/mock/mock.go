package mock

import (
	"fmt"
	"sync"

	"example.com/game/game/model"
)

// MockEngine is a lightweight engine replacement for tests.
type MockEngine struct {
	In     chan model.UserRequest
	winner int
	mu     sync.Mutex
}

func NewMock(buffer int) *MockEngine {
	me := &MockEngine{In: make(chan model.UserRequest, buffer)}
	go me.loop()
	return me
}

func (m *MockEngine) loop() {
	for req := range m.In {
		if req.Correct {
			m.mu.Lock()
			if m.winner == 0 {
				m.winner = req.UserID
				fmt.Printf("[mock] Winner found: User %d\n", req.UserID)
			}
			m.mu.Unlock()
		}
	}
}

func (m *MockEngine) WinnerID() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.winner
}
