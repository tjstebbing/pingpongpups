package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	pongCount int
	mu        sync.Mutex
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		mu.Lock()
		pongCount++
		mu.Unlock()
		json.NewEncoder(w).Encode(map[string]bool{"pong": true})
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count := pongCount
	mu.Unlock()
	fmt.Fprintf(w, "Pong Count: %d", count)
}

func main() {
	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/", handleIndex)

	fmt.Println("Pong server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
