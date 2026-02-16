package resolvers

import (
	"context"
	"net/url"
	"testing"
)

type MockCache struct {
	store map[string]string
}

func (m *MockCache) Get(key string) (string, bool) {
	val, ok := m.store[key]
	return val, ok
}
func (m *MockCache) Set(key string, title string) {
	m.store[key] = title
}
func (m *MockCache) GetMulti(keys []string) map[string]string {
	res := make(map[string]string)
	for _, k := range keys {
		if val, ok := m.store[k]; ok {
			res[k] = val
		}
	}
	return res
}

type MockResolver struct {
	name      string
	canHandle bool
	title     string
}

func (r *MockResolver) Name() string { return r.name }
func (r *MockResolver) CanHandle(u *url.URL) bool { return r.canHandle }
func (r *MockResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	return &Result{Title: r.title, Platform: r.name}, nil
}

func TestResolverManager(t *testing.T) {
	cache := &MockCache{store: make(map[string]string)}
	manager := NewResolverManager(cache)

	r1 := &MockResolver{name: "r1", canHandle: true, title: "Title 1"}
	manager.Register(r1)

	ctx := context.Background()
	urls := []string{"https://example.com/1"}
	
	results := manager.ResolveMulti(ctx, urls)
	
	if results["https://example.com/1"] != "Title 1" {
		t.Errorf("Expected Title 1, got %s", results["https://example.com/1"])
	}

	// Check cache
	if cache.store["https://example.com/1"] != "Title 1" {
		t.Errorf("Expected Title 1 in cache, got %s", cache.store["https://example.com/1"])
	}

	// Test ResolveVideoIDs (Legacy)
	idResults := manager.ResolveVideoIDs(ctx, []string{"abc"})
	if idResults["abc"] != "Title 1" {
		t.Errorf("Expected Title 1 for video ID abc, got %s", idResults["abc"])
	}
}
