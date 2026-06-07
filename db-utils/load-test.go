package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := "http://localhost:8000/incidents"
	concurrency := 10
	iterations := 50

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("Starting load test on %s with concurrency %d...\n", url, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				payload := map[string]interface{}{
					"reporter_id":             fmt.Sprintf("load-test-%d", id),
					"subject":                 "Load Test Incident",
					"description":             "Automated load test grievance submission.",
					"assumed_good_intentions": true,
				}
				jsonData, _ := json.Marshal(payload)
				resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					// Silent fail or log
					continue
				}
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Load test completed in %v. Total requests: %d\n", duration, concurrency*iterations)
}
