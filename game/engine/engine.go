package engine

import (
	"fmt"
	"sync"
	"sync/atomic"

	"example.com/game/game/model"
)

type Engine struct {
	In       chan model.UserRequest
	once     sync.Once
	winnerID atomic.Int64
}

func New(buffer int) *Engine {
	e := &Engine{
		In: make(chan model.UserRequest, buffer),
	}
	go e.loop()
	return e
}

func (e *Engine) loop() {
	for req := range e.In {
		if req.Correct {
			e.once.Do(func() {
				e.winnerID.Store(int64(req.UserID))
				fmt.Printf("Winner found: User %d\n", req.UserID)
			})
		}
	}
}

func (e *Engine) WinnerID() int {
	return int(e.winnerID.Load())
}
