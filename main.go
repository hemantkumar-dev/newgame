package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
	// internal packages removed so this client can run standalone
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// create engine and start service (commented out so this client is standalone)
	// e := engine.New(16384)
	// srv := service.Start(":8080", e)
	// defer srv.Close()

	// // small delay to let server start
	// time.Sleep(100 * time.Millisecond)

	N := 2000
	var wg sync.WaitGroup

	client := &http.Client{Timeout: 2 * time.Second}

	// local UserRequest type to make this client standalone
	type UserRequest struct {
		UserID  int  `json:"userId"`
		Correct bool `json:"correct"`
	}

	for i := 1; i <= N; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			delay := time.Duration(rand.Intn(1000)) * time.Millisecond
			time.Sleep(delay)

			res := UserRequest{
				UserID:  id,
				Correct: rand.Intn(10) == 0,
			}

			data, err := json.Marshal(res)
			if err != nil {
				fmt.Printf("Error marshaling request for user %d: %v\n", id, err)
				return
			}

			_, err = client.Post("http://localhost:8080/submit", "application/json", bytes.NewBuffer(data))
			if err != nil {
				fmt.Printf("Error sending request for user %d: %v\n", id, err)
			}
			// _, err = client.Post("http://localhost:8080/metrics", "application/json", bytes.NewBuffer(data))
		}(i)
	}
	wg.Wait()
	fmt.Println("All requests sent")

	// fetch final metrics from server and print sum
	resp, err := client.Get("http://localhost:8080/metrics")
	if err != nil {
		fmt.Printf("Error fetching final metrics: %v\n", err)
		return
	}
	defer resp.Body.Close()
	var metrics map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		fmt.Printf("Error decoding metrics response: %v\n", err)
		return
	}
	correct := int64(metrics["correct"].(float64))
	incorrect := int64(metrics["incorrect"].(float64))
	fmt.Printf("Final metrics from server: correct=%d incorrect=%d total=%d\n", correct, incorrect, correct+incorrect)
}
