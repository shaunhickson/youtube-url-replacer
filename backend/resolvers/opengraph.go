package resolvers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OpenGraphResolver struct {
	client *http.Client
}

func NewOpenGraphResolver() *OpenGraphResolver {
	return &OpenGraphResolver{
		client: SafeHttpClient(2 * time.Second),
	}
}

func (r *OpenGraphResolver) Name() string {
	return "opengraph"
}

func (r *OpenGraphResolver) CanHandle(u *url.URL) bool {
	// Generic fallback handles everything that looks like a valid http/https URL
	return u.Scheme == "http" || u.Scheme == "https"
}

func (r *OpenGraphResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Set a User-Agent to avoid some basic blocks
	req.Header.Set("User-Agent", "youtube-url-replacer/1.0 (+https://github.com/shaunhickson/youtube-url-replacer)")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res, err := ExtractMetadata(resp.Body)
	if err != nil {
		return nil, err
	}

	res.Platform = "Generic"
	return res, nil
}
