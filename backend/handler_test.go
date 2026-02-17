package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sph/youtube-url-replacer/backend/resolvers"
)

func TestHandler_Limits(t *testing.T) {
	cache := NewInMemoryCache()
	manager := resolvers.NewResolverManager(cache)
	h := NewHandler(cache, manager)
	
	// Configure limits for test
	h.MaxItems = 2
	h.MaxBodyBytes = 100 // Small limit for testing

	t.Run("Valid Request", func(t *testing.T) {
		body := `{"videoIds": ["123"]}`
		req := httptest.NewRequest("POST", "/resolve", strings.NewReader(body))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Valid Request with URLs", func(t *testing.T) {
		body := `{"urls": ["https://example.com"]}`
		req := httptest.NewRequest("POST", "/resolve", strings.NewReader(body))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Too Many Items", func(t *testing.T) {
		body := `{"videoIds": ["1", "2", "3"]}` // 3 items > MaxItems 2
		req := httptest.NewRequest("POST", "/resolve", strings.NewReader(body))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected 413, got %d", w.Code)
		}
	})

	t.Run("Body Too Large", func(t *testing.T) {
		// Create a body larger than 100 bytes
		largeID := strings.Repeat("a", 150)
		body := `{"videoIds": ["` + largeID + `"]}`
		req := httptest.NewRequest("POST", "/resolve", strings.NewReader(body))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected 413, got %d", w.Code)
		}
	})
}
