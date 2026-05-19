package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sph/youtube-url-replacer/backend/resolvers"
	"github.com/sph/youtube-url-replacer/backend/transport"
)

func setupTestServer(t *testing.T) (*httptest.Server, *Handler) {
	// Setup in-memory cache and manager
	cache := NewInMemoryCache()
	manager := resolvers.NewResolverManager(cache)

	// Add generic OpenGraph resolver
	manager.Register(resolvers.NewOpenGraphResolver())
	// Note: We don't add YouTube resolver here because it uses Google API client, 
	// which is harder to mock via simple URL redirection without overriding endpoints.
	// We'll focus integration on the generic OpenGraph resolver which makes raw HTTP calls.

	handler := NewHandler(cache, manager)
	handler.MaxItems = 10
	handler.MaxBodyBytes = 10240

	ts := httptest.NewServer(handler)
	return ts, handler
}

func TestIntegration_ResolveOpenGraph(t *testing.T) {
	transport.AllowLocalIPs = true // Allow test to hit local mock server
	defer func() { transport.AllowLocalIPs = false }()

	// Create a mock target server that returns a valid OpenGraph HTML response
	mockTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`
			<html>
			<head>
				<meta property="og:title" content="Integration Test Title">
			</head>
			<body></body>
			</html>
		`))
	}))
	defer mockTarget.Close()

	// Setup our API server
	apiServer, _ := setupTestServer(t)
	defer apiServer.Close()

	// Prepare request
	reqBody, _ := json.Marshal(map[string][]string{
		"urls": {mockTarget.URL},
	})

	resp, err := http.Post(apiServer.URL+"/resolve", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result struct {
		Titles map[string]string `json:"titles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	title, ok := result.Titles[mockTarget.URL]
	if !ok {
		t.Fatalf("Expected title for URL %s, but not found in response", mockTarget.URL)
	}
	if title != "Integration Test Title" {
		t.Errorf("Expected title 'Integration Test Title', got '%s'", title)
	}
}

func TestIntegration_SSRFProtection(t *testing.T) {
	// Ensure AllowLocalIPs is false (default behavior)
	transport.AllowLocalIPs = false

	// Setup our API server
	apiServer, _ := setupTestServer(t)
	defer apiServer.Close()

	// List of URLs that should be blocked by SSRF protection
	ssrfURLs := []string{
		"http://127.0.0.1:8080/admin",
		"http://169.254.169.254/latest/meta-data/",
		"http://localhost/secret",
	}

	reqBody, _ := json.Marshal(map[string][]string{
		"urls": ssrfURLs,
	})

	resp, err := http.Post(apiServer.URL+"/resolve", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 (since partial failures return 200), got %d", resp.StatusCode)
	}

	var result struct {
		Titles map[string]string `json:"titles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Because all URLs were SSRF attempts, none should have resolved successfully.
	// The titles map should be empty or not contain these URLs.
	for _, u := range ssrfURLs {
		if _, ok := result.Titles[u]; ok {
			t.Errorf("Expected URL %s to be blocked by SSRF protection, but it was resolved", u)
		}
	}
}

func TestIntegration_InvalidURL(t *testing.T) {
	// Setup our API server
	apiServer, _ := setupTestServer(t)
	defer apiServer.Close()

	reqBody, _ := json.Marshal(map[string][]string{
		"urls": {"not-a-valid-url", "ftp://example.com/file"},
	})

	resp, err := http.Post(apiServer.URL+"/resolve", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result struct {
		Titles map[string]string `json:"titles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result.Titles) != 0 {
		t.Errorf("Expected 0 resolved titles for invalid URLs, got %d", len(result.Titles))
	}
}
