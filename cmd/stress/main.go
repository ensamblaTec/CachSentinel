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
		totalHits   = 100000
		duration    = 15 * time.Second
		targetURL   = "http://localhost:8080/posts/1"
	)

	var wg sync.WaitGroup
	latencies := make(chan time.Duration, totalHits)

	fmt.Printf("Iniciando ataque: %d peticiones con %d workers...\n", totalHits, concurrency)

	start := time.Now()
	deadline := time.After(duration)
	ticker := time.NewTicker(concurrency * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			goto done
		case <-ticker.C:
			for range concurrency {
				wg.Go(func() {
					startReq := time.Now()
					resp, err := http.Get(targetURL)
					if err == nil {
						io.Copy(io.Discard, resp.Body)
						resp.Body.Close()
					}
					latencies <- time.Since(startReq)
				})
			}
		}
	}

done:
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
