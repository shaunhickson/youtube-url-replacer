package resolvers

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeResolver struct {
	service *youtube.Service
}

func NewYouTubeResolver(apiKey string) (*YouTubeResolver, error) {
	if apiKey == "" {
		return &YouTubeResolver{service: nil}, nil // Mock mode
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating youtube service: %v", err)
	}
	return &YouTubeResolver{service: service}, nil
}

func (r *YouTubeResolver) Name() string {
	return "youtube"
}

func (r *YouTubeResolver) CanHandle(u *url.URL) bool {
	host := strings.ToLower(u.Host)
	return host == "youtube.com" || host == "www.youtube.com" || host == "youtu.be"
}

func (r *YouTubeResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	videoID := ""

	if strings.ToLower(u.Host) == "youtu.be" {
		videoID = strings.TrimPrefix(u.Path, "/")
	} else {
		videoID = u.Query().Get("v")
		if videoID == "" && strings.HasPrefix(u.Path, "/shorts/") {
			videoID = strings.TrimPrefix(u.Path, "/shorts/")
		}
		if videoID == "" && strings.HasPrefix(u.Path, "/live/") {
			videoID = strings.TrimPrefix(u.Path, "/live/")
		}
	}

	if videoID == "" {
		return nil, fmt.Errorf("could not extract video ID from YouTube URL: %s", u.String())
	}

	// Handle mock mode
	if r.service == nil {
		return &Result{
			Title:    fmt.Sprintf("Mock Title for Video %s", videoID),
			Platform: "YouTube",
		}, nil
	}

	call := r.service.Videos.List([]string{"snippet"}).Id(videoID)
	response, err := call.Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("error calling youtube api: %v", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found: %s", videoID)
	}

	item := response.Items[0]
	return &Result{
		Title:    item.Snippet.Title,
		Platform: "YouTube",
	}, nil
}
