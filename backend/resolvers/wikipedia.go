package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WikipediaResolver struct {
	client *http.Client
}

func NewWikipediaResolver() *WikipediaResolver {
	return &WikipediaResolver{
		client: SafeHttpClient(2 * time.Second),
	}
}

func (r *WikipediaResolver) Name() string {
	return "wikipedia"
}

func (r *WikipediaResolver) CanHandle(u *url.URL) bool {
	host := strings.ToLower(u.Host)
	return strings.HasSuffix(host, "wikipedia.org") && strings.HasPrefix(u.Path, "/wiki/")
}

func (r *WikipediaResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	// Extract the article title from the URL path (/wiki/Title)
	parts := strings.Split(u.Path, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid wikipedia path")
	}
	articleTitle := parts[2]
	
	// Wikipedia API: https://en.wikipedia.org/api/rest_v1/page/summary/{title}
	// Note: We extract the language subdomain from the host (e.g. en.wikipedia.org -> en)
	hostParts := strings.Split(strings.ToLower(u.Host), ".")
	lang := "en"
	if len(hostParts) == 3 {
		lang = hostParts[0]
	}

	apiURL := fmt.Sprintf("https://%s.wikipedia.org/api/rest_v1/page/summary/%s", lang, articleTitle)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "youtube-url-replacer/1.0 (https://github.com/shaunhickson/youtube-url-replacer)")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	title := data.Title
	if data.Description != "" {
		title = fmt.Sprintf("%s - %s", title, data.Description)
	}

	return &Result{
		Title:    title,
		Platform: "Wikipedia",
	}, nil
}
