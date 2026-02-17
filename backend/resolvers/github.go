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

type GitHubResolver struct {
	client  *http.Client
	token   string
	baseURL string
}

type githubRepoResponse struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	Language        string `json:"language"`
	ForksCount      int    `json:"forks_count"`
}

func NewGitHubResolver(token string) *GitHubResolver {
	return &GitHubResolver{
		client:  SafeHttpClient(2 * time.Second),
		token:   token,
		baseURL: "https://api.github.com",
	}
}

func (r *GitHubResolver) Name() string {
	return "github"
}

func (r *GitHubResolver) CanHandle(u *url.URL) bool {
	host := strings.ToLower(u.Host)
	if host != "github.com" && host != "www.github.com" {
		return false
	}

	// Pattern: /owner/repo
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		return false
	}

	// Avoid common non-repo paths
	reserved := map[string]bool{
		"settings":      true,
		"notifications": true,
		"explore":       true,
		"trending":      true,
		"marketplace":   true,
		"features":      true,
		"topics":        true,
	}
	if reserved[parts[0]] {
		return false
	}

	return true
}

func (r *GitHubResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	owner := parts[0]
	repo := parts[1]

	apiURL := fmt.Sprintf("%s/repos/%s/%s", r.baseURL, owner, repo)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "youtube-url-replacer/1.0 (+https://github.com/shaunhickson/youtube-url-replacer)")
	if r.token != "" {
		req.Header.Set("Authorization", "token "+r.token)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Let it fallback to OpenGraph
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var data githubRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	stats := fmt.Sprintf("★ %d", data.StargazersCount)
	if data.Language != "" {
		stats = fmt.Sprintf("%s | %s", stats, data.Language)
	}

	return &Result{
		Title:       data.FullName,
		Description: fmt.Sprintf("%s (%s)", data.Description, stats),
		Platform:    "github",
	}, nil
}
