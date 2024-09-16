package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	successCount      int
	mu                sync.Mutex
	ppm               float64
	lastPingTimestamp time.Time
)

func pingPong() {
	ppmTicker := time.NewTicker(1 * time.Minute)
	defer ppmTicker.Stop()
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-ppmTicker.C:
			// Calculate PPM
			mu.Lock()
			interval := time.Since(lastPingTimestamp).Minutes()
			if interval > 0 {
				ppm = float64(successCount) / interval
			} else {
				ppm = 0
			}
			successCount = 0
			lastPingTimestamp = time.Now()
			mu.Unlock()
			// Post PPM to dogebox metrics endpoint in a new goroutine
			go postPPMToMetrics()
		case <-pingTicker.C:
			// Send ping
			resp, err := http.Post("http://dogebox:8082/pingpong/ping", "application/json", nil)
			if err == nil && resp.StatusCode == http.StatusOK {
				var result map[string]bool
				if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && result["pong"] {
					mu.Lock()
					successCount++
					lastPingTimestamp = time.Now()
					mu.Unlock()
				}
				resp.Body.Close()
			}
		}
	}
}

func postPPMToMetrics() {
	data := map[string]map[string]float64{
		"ppm": {"value": ppm},
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
