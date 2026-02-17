package transport

import (
	"context"
	"net"
	"testing"
)

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		ip      string
		private bool
	}{
		{"127.0.0.1", true},
		{"10.0.0.1", true},
		{"192.168.1.1", true},
		{"172.16.0.1", true},
		{"172.31.255.255", true}, // Private range
		{"169.254.169.254", true}, // Link-local
		{"8.8.8.8", false},       // Google DNS
		{"1.1.1.1", false},       // Cloudflare
		{"::1", true},            // IPv6 Loopback
		{"fc00::1", true},        // IPv6 Unique Local
		{"2001:4860:4860::8888", false}, // IPv6 Public
	}

	for _, tc := range tests {
		ip := net.ParseIP(tc.ip)
		if got := isPrivateIP(ip); got != tc.private {
			t.Errorf("isPrivateIP(%s) = %v; want %v", tc.ip, got, tc.private)
		}
	}
}

func TestSafeDialer_Blocked(t *testing.T) {
	// Create a safe dialer
	dialer := &net.Dialer{}
	safeDial := SafeDialer(dialer)

	// Try to dial localhost (should fail)
	// We use a random port that is likely closed, but the blocking happens BEFORE connection
	_, err := safeDial(context.Background(), "tcp", "127.0.0.1:1234")
	if err == nil {
		t.Error("Expected error for 127.0.0.1, got nil")
	} else if err.Error() != "blocked: resolves to private/local IP" {
		t.Errorf("Expected blocked error, got: %v", err)
	}

	// Try to dial private IP
	_, err = safeDial(context.Background(), "tcp", "192.168.1.1:80")
	if err == nil {
		t.Error("Expected error for 192.168.1.1, got nil")
	}
}
