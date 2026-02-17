package resolvers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// UnshortenerResolver follows redirect chains to find the final URL
type UnshortenerResolver struct {
	client  *http.Client
	manager *ResolverManager
	domains []string
}

// NewUnshortenerResolver creates a new resolver for shortened URLs
func NewUnshortenerResolver(manager *ResolverManager) *UnshortenerResolver {
	// Known shortener domains
	domains := []string{
		"bit.ly", "t.co", "tinyurl.com", "is.gd", "buff.ly",
		"goo.gl", "bit.do", "ow.ly", "t.ly", "shorturl.at",
	}

	return &UnshortenerResolver{
		client:  SafeHttpClient(2 * time.Second),
		manager: manager,
		domains: domains,
	}
}

func (r *UnshortenerResolver) Name() string {
	return "unshortener"
}

func (r *UnshortenerResolver) CanHandle(u *url.URL) bool {
	host := strings.ToLower(u.Host)
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	for _, d := range r.domains {
		if host == d {
			return true
		}
	}
	return false
}

func (r *UnshortenerResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	currentURL := u.String()
	hops := 0
	maxHops := 5
	seen := make(map[string]bool)

	for hops < maxHops {
		if seen[currentURL] {
			return nil, fmt.Errorf("redirect loop detected at %s", currentURL)
		}
		seen[currentURL] = true

		req, err := http.NewRequestWithContext(ctx, "HEAD", currentURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "youtube-url-replacer/1.0 (+https://github.com/shaunhickson/youtube-url-replacer)")

		// We use a client that DOES NOT automatically follow redirects so we can track them
		resp, err := r.client.Transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location := resp.Header.Get("Location")
			if location == "" {
				break // Redirect without location? Stop here.
			}

			// Handle relative URLs
			nextURL, err := u.Parse(location)
			if err != nil {
				break
			}
			currentURL = nextURL.String()
			u = nextURL // Update u for relative parsing in next hop
			hops++
		} else {
			// Not a redirect, we've reached the end
			break
		}
	}

	finalURL, err := url.Parse(currentURL)
	if err != nil {
		return nil, err
	}

	// Now that we have the final URL, let the manager resolve it properly
	// This allows it to hit YouTube, OpenGraph, etc.
	// We call Resolve on the manager but we must be careful of recursion
	// The manager already has a list of resolvers.
	return r.manager.resolveRecursively(ctx, finalURL, r.Name())
}
