package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	successCount int
	mu           sync.Mutex
)

func pingPong() {
	for {
		resp, err := http.Post("http://dogebox:8082/pingpong/ping", "application/json", nil)
		if err == nil && resp.StatusCode == http.StatusOK {
			var result map[string]bool
			if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && result["pong"] {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
			resp.Body.Close()
		}
		time.Sleep(5 * time.Second)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count := successCount
	mu.Unlock()
	fmt.Fprintf(w, "Ping Success Count: %d", count)
}

func main() {
	go pingPong()

	http.HandleFunc("/", handleIndex)
	fmt.Println("Ping server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
