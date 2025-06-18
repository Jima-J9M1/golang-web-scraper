// main.go (modified for database integration)
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"webscraper/internal/scraper"
	"webscraper/internal/storage" // <-- NEW: Import storage package
)

func main() {
	fmt.Println("Web Scraper Project Started!")
	timeNow := time.Now()
	urls := []string{
		"http://example.com",
		"http://www.iana.org/domains/example",
		"https://www.w3.org/Consortium/fees",
		"https://www.google.com/search?q=golang",
		"http://nonexistent.invalid",
		"https://go.dev/",
		"https://pkg.go.dev/",
		"https://tour.golang.org/welcome/1",
		// "http://localhost:9999", // To test connection timeout
	}

	// --- Database Initialization ---
	db, err := storage.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close() // Ensure the database connection is closed when main exits

	// Determine the number of worker goroutines to use based on the number of CPU cores available.
	// This is done by multiplying the number of CPU cores by 2 to ensure efficient use of resources.
	numWorkers := runtime.NumCPU() * 2
	if numWorkers == 0 {
		numWorkers = 1
	}

	jobsChan := make(chan string, len(urls))
	resultsChan := make(chan scraper.ScrapedResult, len(urls))
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("Worker %d started.\n", workerID)
			for url := range jobsChan {
				fmt.Printf("Worker %d: Fetching %s\n", workerID, url)

				bodyReader, err := scraper.FetchURL(ctx, url)
				if err != nil {
					log.Printf("Worker %d: Error fetching URL %s: %v", workerID, url, err)
					resultsChan <- scraper.ScrapedResult{URL: url, Err: err}
					continue
				}
				defer bodyReader.Close()

				links, err := scraper.ParseLinks(bodyReader, url)
				if err != nil {
					log.Printf("Worker %d: Error parsing links from %s: %v", workerID, url, err)
					resultsChan <- scraper.ScrapedResult{URL: url, Err: err}
					continue
				}

				// --- Database Insertion ---
				err = storage.InsertLinks(db, url, links) // Insert the links
				if err != nil {
					log.Printf("Worker %d: Error inserting links from %s into database: %v", workerID, url, err)
					resultsChan <- scraper.ScrapedResult{URL: url, Err: fmt.Errorf("database insertion failed: %w", err)} // Wrap the db error
					continue                                                                                              // Don't send the links to the resultsChan if insertion failed
				}

				resultsChan <- scraper.ScrapedResult{URL: url, Links: links}
			}
			fmt.Printf("Worker %d finished.\n", workerID)
		}(i + 1)
	}

	// 2. Send jobs to the jobs channel
	for _, url := range urls {
		select {
		case <-ctx.Done():
			fmt.Printf("Context cancelled. Stopping job submission. Error: %v\n", ctx.Err())
			break // Exit the loop
		case jobsChan <- url:
			// Job sent successfully
		}
	}
	close(jobsChan)

	// 3. Wait for all workers to finish and then close the results channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 4. Collect results from the results channel
	for result := range resultsChan {
		if result.Err != nil {
			fmt.Printf("\n--- Error processing %s ---\n", result.URL)
			var fetchErr *scraper.ErrFetchFailed
			var parseErr *scraper.ErrParseFailed
			var dbErr error // Generic error variable for database errors

			if errors.As(result.Err, &fetchErr) {
				fmt.Printf("  Fetch Error! Reason: %s (Status: %d, Underlying: %v)\n",
					fetchErr.Reason, fetchErr.StatusCode, fetchErr.WrappedErr)
			} else if errors.As(result.Err, &parseErr) {
				fmt.Printf("  Parse Error! Reason: %s (Underlying: %v)\n",
					parseErr.Reason, parseErr.WrappedErr)
			} else if errors.As(result.Err, &dbErr) { // Check for database errors
				fmt.Printf("  Database Error! %v\n", dbErr)
			} else {
				fmt.Printf("  Unhandled Error: %v\n", result.Err)
			}
		} else {
			fmt.Printf("\n--- Scraped from %s and stored in database ---\n", result.URL) // Modified message
			// We don't need to print the links here anymore, as they are in the database.
		}
	}

	fmt.Println("\nAll scraping tasks completed and results processed!")
	fmt.Println("Operation stopped:", time.Since(timeNow))
}
