package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example.com/game/game/engine"
	"example.com/game/game/model"
)

// Start starts an HTTP server that accepts POST /submit and forwards requests to the engine.
// It runs the server in a new goroutine and returns the *http.Server so the caller can shutdown.
func Start(addr string, e *engine.Engine) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req model.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.UserID <= 0 {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		// forward to engine (buffered channel)
		select {
		case e.In <- req:
			// queued
		default:
			// if channel full, drop but still respond OK
			fmt.Println("engine queue full, dropped request")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
	})

	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		log.Printf("service listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	return srv
}
