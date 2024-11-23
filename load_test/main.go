package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Set up the load test parameters
	numRequests := 100000
	timeout := time.Duration(5 * time.Second)

	// Create a wait group to ensure all goroutines have finished
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Create a client with a timeout
	client := http.Client{
		Timeout: timeout,
	}

	// Start the load test
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		// Start a goroutine for each request
		go func() {
			defer wg.Done()

			id := rand.Int()
			url := fmt.Sprintf("http://localhost:8080/api/verve/accept?id=%d", id)
			method := "GET"

			req, err := http.NewRequest(method, url, nil)

			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Add("cache-control", "no-cache")
			req.Header.Add("content-type", "application/json")

			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer res.Body.Close()
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Calculate the load test duration and requests per second
	duration := time.Since(startTime)
	requestsPerSecond := float64(numRequests) / duration.Seconds()

	fmt.Printf("Load test completed in %v with %d requests per second\n", duration, int(requestsPerSecond))
}
