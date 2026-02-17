package resolvers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sph/youtube-url-replacer/backend/transport"
)

func TestUnshortenerResolver(t *testing.T) {
	transport.AllowLocalIPs = true
	defer func() { transport.AllowLocalIPs = false }()

	mux := http.NewServeMux()

	// 1. Simple redirect
	mux.HandleFunc("/short", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/final", http.StatusFound)
	})
	mux.HandleFunc("/final", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write([]byte("<html><head><title>Final Destination</title></head></html>")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// 2. Loop
	mux.HandleFunc("/loop1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/loop2", http.StatusFound)
	})
	mux.HandleFunc("/loop2", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/loop1", http.StatusFound)
	})

	// 3. YouTube link
	mux.HandleFunc("/to-yt", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	cache := &MockCache{store: make(map[string]string)}
	manager := NewResolverManager(cache)
	
	unshortener := NewUnshortenerResolver(manager)
	// Override domains for testing
	u, _ := url.Parse(ts.URL)
	unshortener.domains = []string{u.Host}

	og := NewOpenGraphResolver()
	yt, _ := NewYouTubeResolver("") // Mock YT

	manager.Register(yt)
	manager.Register(unshortener)
	manager.Register(og)

	ctx := context.Background()

	t.Run("Simple Redirect", func(t *testing.T) {
		u, _ := url.Parse(ts.URL + "/short")
		res, err := unshortener.Resolve(ctx, u)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Title != "Final Destination" {
			t.Errorf("Expected 'Final Destination', got '%s'", res.Title)
		}
	})

	t.Run("Redirect Loop", func(t *testing.T) {
		u, _ := url.Parse(ts.URL + "/loop1")
		_, err := unshortener.Resolve(ctx, u)
		if err == nil {
			t.Fatal("Expected error for redirect loop, got nil")
		}
	})

	t.Run("Redirect to YouTube", func(t *testing.T) {
		u, _ := url.Parse(ts.URL + "/to-yt")
		res, err := unshortener.Resolve(ctx, u)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		// YouTube mock returns "Mock Title for Video ..."
		if res.Title != "Mock Title for Video dQw4w9WgXcQ" {
			t.Errorf("Expected YT mock title, got '%s'", res.Title)
		}
	})
}
