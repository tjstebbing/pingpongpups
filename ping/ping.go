package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	successCount int
	failureCount int
	ppm          float64
	fpm          int
)

var target string = "http://dogebox:8082/pingpong/ping"

func pingPong() {
	ppmTicker := time.NewTicker(1 * time.Minute)
	defer ppmTicker.Stop()
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-ppmTicker.C:
			// Calculate PPM
			ppm = float64(successCount)
			fpm = failureCount
			successCount = 0
			failureCount = 0
			// Post dogebox metrics
			go postMetrics()
		case <-pingTicker.C:
			// Send ping
			resp, err := http.Post(target, "application/json", nil)
			if err == nil && resp.StatusCode == http.StatusOK {
				var result map[string]bool
				if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && result["pong"] {
					successCount++
				} else {
					failureCount++
				}
				resp.Body.Close()
			}
		}
	}
}

func postMetrics() {
	data := map[string]map[string]any{
		"ppm":    {"value": ppm},
		"fpm":    {"value": fpm},
		"target": {"value": target},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post("http://dogebox:8082/dbx/metrics", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error posting metrics:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-OK status code:", resp.Status)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Fprintf(w, "Ping Success Count: %d\nPPM: %.2f", successCount, ppm)
}

func main() {
	go pingPong()

	http.HandleFunc("/", handleIndex)
	fmt.Println("Ping server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
