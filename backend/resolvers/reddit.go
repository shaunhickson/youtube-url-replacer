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

type RedditResolver struct {
	client *http.Client
}

func NewRedditResolver() *RedditResolver {
	return &RedditResolver{
		client: SafeHttpClient(2 * time.Second),
	}
}

func (r *RedditResolver) Name() string {
	return "reddit"
}

func (r *RedditResolver) CanHandle(u *url.URL) bool {
	host := strings.ToLower(u.Host)
	return (host == "reddit.com" || host == "www.reddit.com" || host == "old.reddit.com") && strings.Contains(u.Path, "/comments/")
}

func (r *RedditResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	// Construct the JSON URL
	jsonURL := fmt.Sprintf("https://www.reddit.com%s.json", strings.TrimSuffix(u.Path, "/"))

	req, err := http.NewRequestWithContext(ctx, "GET", jsonURL, nil)
	if err != nil {
		return nil, err
	}
	// Reddit requires a custom User-Agent to not rate limit heavily
	req.Header.Set("User-Agent", "youtube-url-replacer/1.0")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data []struct {
		Data struct {
			Children []struct {
				Data struct {
					Title     string `json:"title"`
					Subreddit string `json:"subreddit_name_prefixed"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data) == 0 || len(data[0].Data.Children) == 0 {
		return nil, fmt.Errorf("no post data found")
	}

	post := data[0].Data.Children[0].Data
	title := post.Title
	if post.Subreddit != "" {
		title = fmt.Sprintf("%s - %s", post.Subreddit, title)
	}

	return &Result{
		Title:    title,
		Platform: "Reddit",
	}, nil
}
