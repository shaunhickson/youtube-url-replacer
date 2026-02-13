package main

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeService struct {
	service *youtube.Service
}

func NewYouTubeService(apiKey string) (*YouTubeService, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating youtube service: %v", err)
	}
	return &YouTubeService{service: service}, nil
}

// FetchTitles retrieves titles for a list of video IDs
// It respects the API limit of 50 items per request, but we will chunk outside if needed.
// For now, let's assume the caller chunks or we handle small batches.
func (s *YouTubeService) FetchTitles(videoIDs []string) (map[string]string, error) {
	// MOCK MODE: If service is nil, return fake titles
	if s.service == nil {
		mockTitles := make(map[string]string)
		for _, id := range videoIDs {
			mockTitles[id] = fmt.Sprintf("Mock Title for Video %s", id)
		}
		return mockTitles, nil
	}

	if len(videoIDs) == 0 {
		return map[string]string{}, nil
	}

	call := s.service.Videos.List([]string{"snippet"}).Id(videoIDs...)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error calling youtube api: %v", err)
	}

	titles := make(map[string]string)
	for _, item := range response.Items {
		titles[item.Id] = item.Snippet.Title
	}

	return titles, nil
}
