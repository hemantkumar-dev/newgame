package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

type UserRequest struct {
	UserID  int  `json:"userId"`
	Correct bool `json:"correct"`
}

var (
	winnerFound bool
	winnerID    int
	mu          sync.Mutex
	// metrics
	correctCount   int64
	incorrectCount int64
)

func submitHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// update metrics (count every received submission)
	if req.Correct {
		atomic.AddInt64(&correctCount, 1)
	} else {
		atomic.AddInt64(&incorrectCount, 1)
	}
	mu.Lock()
	defer mu.Unlock()
	if winnerFound {
		return
	}

	if req.Correct {
		winnerFound = true
		winnerID = req.UserID
		// print winner and final metrics exactly once
		fmt.Printf("Winner found: User %d\n", winnerID)
	}
}

// metricsHandler returns JSON with current metrics and winner info
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// read winner state under lock
	mu.Lock()
	wf := winnerFound
	wid := winnerID
	mu.Unlock()

	m := map[string]interface{}{
		"correct":     atomic.LoadInt64(&correctCount),
		"incorrect":   atomic.LoadInt64(&incorrectCount),
		"winnerFound": wf,
		"winnerID":    wid,
	}
	json.NewEncoder(w).Encode(m)
}

func main() {
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/metrics", metricsHandler)
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
