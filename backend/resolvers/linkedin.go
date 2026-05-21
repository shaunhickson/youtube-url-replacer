package resolvers

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type LinkedInResolver struct {
	client      *http.Client
	accessToken string
}

func NewLinkedInResolver(accessToken string) *LinkedInResolver {
	if accessToken == "" {
		return nil
	}
	return &LinkedInResolver{
		client:      SafeHttpClient(3 * time.Second),
		accessToken: accessToken,
	}
}

func (r *LinkedInResolver) Name() string {
	return "linkedin"
}

func (r *LinkedInResolver) CanHandle(u *url.URL) bool {
	host := strings.ToLower(u.Host)
	return (host == "linkedin.com" || host == "www.linkedin.com") && (strings.Contains(u.Path, "/posts/") || strings.Contains(u.Path, "/in/"))
}

func (r *LinkedInResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	// Note: LinkedIn's APIs are highly restricted.
	// This acts as a placeholder that would be expanded with actual
	// LinkedIn Community Management API calls using the accessToken.
	
	title := "LinkedIn Profile/Post"
	if strings.Contains(u.Path, "/in/") {
		parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
		if len(parts) >= 2 {
			title = "LinkedIn Profile: " + parts[1]
		}
	} else if strings.Contains(u.Path, "/posts/") {
		title = "LinkedIn Post"
	}

	return &Result{
		Title:    title,
		Platform: "LinkedIn",
	}, nil
}
