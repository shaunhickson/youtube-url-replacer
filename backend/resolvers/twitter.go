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

type TwitterResolver struct {
	client      *http.Client
	bearerToken string
}

func NewTwitterResolver(bearerToken string) *TwitterResolver {
	if bearerToken == "" {
		return nil
	}
	return &TwitterResolver{
		client:      SafeHttpClient(3 * time.Second),
		bearerToken: bearerToken,
	}
}

func (r *TwitterResolver) Name() string {
	return "twitter"
}

func (r *TwitterResolver) CanHandle(u *url.URL) bool {
	host := strings.ToLower(u.Host)
	return (host == "twitter.com" || host == "x.com" || host == "www.twitter.com" || host == "www.x.com") && strings.Contains(u.Path, "/status/")
}

func (r *TwitterResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	parts := strings.Split(u.Path, "/")
	// /username/status/12345
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid twitter path")
	}

	tweetID := parts[3]
	
	apiURL := fmt.Sprintf("https://api.twitter.com/2/tweets/%s?expansions=author_id", tweetID)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+r.bearerToken)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitter api error, status: %d", resp.StatusCode)
	}

	var data struct {
		Data struct {
			Text string `json:"text"`
		} `json:"data"`
		Includes struct {
			Users []struct {
				Name string `json:"name"`
			} `json:"users"`
		} `json:"includes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	title := data.Data.Text
	if len(data.Includes.Users) > 0 {
		authorName := data.Includes.Users[0].Name
		title = fmt.Sprintf("%s - %s", authorName, title)
	}

	// Truncate if too long (tweets are up to 280-10000 chars, let's keep it clean)
	if len(title) > 100 {
		title = title[:97] + "..."
	}
	// replace newlines
	title = strings.ReplaceAll(title, "\n", " ")

	return &Result{
		Title:    title,
		Platform: "Twitter",
	}, nil
}
