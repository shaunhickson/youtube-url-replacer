package resolvers

import (
	"context"
	"net/url"
	"testing"
)

func TestYouTubeResolver_CanHandle(t *testing.T) {
	r, _ := NewYouTubeResolver("") // Mock mode
	tests := []struct {
		urlStr string
		want   bool
	}{
		{"https://youtube.com/watch?v=123", true},
		{"https://www.youtube.com/watch?v=123", true},
		{"https://youtu.be/123", true},
		{"https://youtube.com/shorts/123", true},
		{"https://google.com", false},
	}

	for _, tt := range tests {
		u, _ := url.Parse(tt.urlStr)
		if got := r.CanHandle(u); got != tt.want {
			t.Errorf("CanHandle(%s) = %v, want %v", tt.urlStr, got, tt.want)
		}
	}
}

func TestYouTubeResolver_Resolve(t *testing.T) {
	r, _ := NewYouTubeResolver("") // Mock mode
	ctx := context.Background()

	u, _ := url.Parse("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	res, err := r.Resolve(ctx, u)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if res.Title != "Mock Title for Video dQw4w9WgXcQ" {
		t.Errorf("Unexpected title: %s", res.Title)
	}

	u2, _ := url.Parse("https://youtu.be/dQw4w9WgXcQ")
	res2, err := r.Resolve(ctx, u2)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if res2.Title != "Mock Title for Video dQw4w9WgXcQ" {
		t.Errorf("Unexpected title: %s", res2.Title)
	}
	
	u3, _ := url.Parse("https://www.youtube.com/shorts/abc")
	res3, err := r.Resolve(ctx, u3)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if res3.Title != "Mock Title for Video abc" {
		t.Errorf("Unexpected title: %s", res3.Title)
	}
}
