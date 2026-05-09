package adapter

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HTTPFetcher struct {
	BaseURL string
	Client  *http.Client
}

func (fetcher *HTTPFetcher) Fetch(ctx context.Context, key string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", fetcher.BaseURL, key)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := fetcher.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream_status_%d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream_status_%d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
