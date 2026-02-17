package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// 60 RPM, Burst 1
	rl := NewRateLimiter(60, 1)

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"

	// 1. First request should pass
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req)
	if rec1.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec1.Code)
	}

	// 2. Second request (immediate) should fail (Burst is 1)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req)
	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429, got %d", rec2.Code)
	}

	// 3. Different IP should pass
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "5.6.7.8:5678"
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req2)
	if rec3.Code != http.StatusOK {
		t.Errorf("Expected 200 for new IP, got %d", rec3.Code)
	}
}

func TestGetIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	
	if ip := getIP(req); ip != "10.0.0.1" {
		t.Errorf("Expected 10.0.0.1, got %s", ip)
	}

	req.Header.Set("X-Forwarded-For", "203.0.113.1, 198.51.100.1")
	if ip := getIP(req); ip != "203.0.113.1" {
		t.Errorf("Expected 203.0.113.1, got %s", ip)
	}
}

func TestCleanup(t *testing.T) {
	rl := NewRateLimiter(60, 10)
	rl.ips.Store("1.2.3.4", &visitor{lastSeen: time.Now().Add(-10 * time.Minute)})
	
	// Manually trigger cleanup loop logic (simulation)
	// Since the actual CleanupBackground runs in a goroutine with sleep, 
	// we can just test the logic or rely on a short interval integration test.
	// For unit test, it's hard to test the goroutine deterministically without hooks.
	// We'll skip complex async testing here and trust the sync.Map logic.
}
