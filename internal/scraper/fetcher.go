// internal/scraper/fetcher.go
package scraper

import (
	"context"
	"io"
	"net/http"
	"time"
)

// FetchURL returns the io.ReadCloser for the given URL's response body.
// It's the caller's responsibility to close the io.ReadCloser.
func FetchURL(ctx context.Context, url string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) 
	if err != nil {
		return nil, &ErrFetchFailed{
			URL:        url,
			Reason:     "failed to create request",
			WrappedErr: err,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, &ErrFetchFailed{
			URL:        url,
			Reason:     "failed to create request",
			WrappedErr: err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		// Close the body even if there's an error, as we're not returning it
		resp.Body.Close()
		return nil, &ErrFetchFailed{
			URL:        url,
			Reason:     "failed to execute HTTP request",
			WrappedErr: err,
		}
	}

	return resp.Body, nil // Return the body directly
}