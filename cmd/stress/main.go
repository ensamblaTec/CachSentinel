package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	const (
		concurrency = 50
		totalHits   = 500
		targetURL   = "http://localhost:8080/posts/1"
	)

	var wg sync.WaitGroup
	start := time.Now()

	latencies := make(chan time.Duration, totalHits)

	fmt.Printf("Iniciando ataque: %d peticiones con %d workers...\n", totalHits, concurrency)

	semaphore := make(chan struct{}, concurrency)

	for index := range totalHits {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			reqStart := time.Now()
			resp, err := http.Get(targetURL)
			if err != nil {
				fmt.Printf("error req %d: %v\n", id, err)
				return
			}

			_, _ = io.Copy(io.Discard, resp.Body)
			defer resp.Body.Close()

			latencies <- time.Since(reqStart)

			time.Sleep(10 * time.Millisecond)
		}(index)
	}

	wg.Wait()
	close(latencies)

	var total time.Duration
	var count int
	for lat := range latencies {
		total += lat
		count++
	}

	fmt.Println("\n--- Reporte de Rendimiento ---")
	fmt.Printf("Tiempo Total:    %v\n", time.Since(start))
	fmt.Printf("Latencia Media:  %v\n", total/time.Duration(count))
	fmt.Printf("Requests/sec:    %.2f\n", float64(count)/time.Since(start).Seconds())
}
