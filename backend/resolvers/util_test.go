package resolvers

import (
	"strings"
	"testing"
	"time"
)

func TestSafeHttpClient_BlocksPrivate(t *testing.T) {
	client := SafeHttpClient(1 * time.Second)
	_, err := client.Get("http://127.0.0.1:12345")
	if err == nil {
		t.Error("Expected error for 127.0.0.1, got nil")
	} else if !strings.Contains(err.Error(), "blocked") {
		t.Errorf("Expected 'blocked' error, got: %v", err)
	}
}

func TestExtractMetadata(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		want     string
		wantDesc string
	}{
		{
			name: "OG Title",
			html: `<html><head><meta property="og:title" content="OG Title"><title>Fallback Title</title></head></html>`,
			want: "OG Title",
		},
		{
			name: "Standard Title",
			html: `<html><head><title>Standard Title</title></head></html>`,
			want: "Standard Title",
		},
		{
			name: "With Entities",
			html: `<html><head><title>A &amp; B &quot;C&quot;</title></head></html>`,
			want: `A & B "C"`,
		},
		{
			name: "OG Description",
			html: `<html><head><meta property="og:title" content="Title"><meta property="og:description" content="Desc"></head></html>`,
			want: "Title",
			wantDesc: "Desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ExtractMetadata(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("ExtractMetadata failed: %v", err)
			}
			if res.Title != tt.want {
				t.Errorf("Got title %q, want %q", res.Title, tt.want)
			}
			if tt.wantDesc != "" && res.Description != tt.wantDesc {
				t.Errorf("Got description %q, want %q", res.Description, tt.wantDesc)
			}
		})
	}
}
