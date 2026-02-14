package resolvers

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"
)

type ResolverManager struct {
	resolvers []Resolver
	cache     Cache
	timeout   time.Duration
}

func NewResolverManager(cache Cache) *ResolverManager {
	return &ResolverManager{
		resolvers: []Resolver{},
		cache:     cache,
		timeout:   2 * time.Second, // Default timeout
	}
}

func (m *ResolverManager) SetTimeout(t time.Duration) {
	m.timeout = t
}

func (m *ResolverManager) Register(r Resolver) {
	m.resolvers = append(m.resolvers, r)
}

func (m *ResolverManager) ResolveMulti(ctx context.Context, urls []string) map[string]string {
	// Apply global timeout if not already set on context
	if m.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, m.timeout)
		defer cancel()
	}

	results := make(map[string]string)
	var missingURLs []string

	// 1. Check Cache
	cached := m.cache.GetMulti(urls)
	for _, u := range urls {
		if val, ok := cached[u]; ok {
			results[u] = val
		} else {
			missingURLs = append(missingURLs, u)
		}
	}

	if len(missingURLs) == 0 {
		return results
	}

	// 2. Resolve missing URLs
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, rawURL := range missingURLs {
		wg.Add(1)
		go func(raw string) {
			defer wg.Done()

			u, err := url.Parse(raw)
			if err != nil {
				log.Printf("Failed to parse URL %s: %v", raw, err)
				return
			}

			for _, r := range m.resolvers {
				if r.CanHandle(u) {
					res, err := r.Resolve(ctx, u)
					if err != nil {
						log.Printf("Resolver %s failed for %s: %v", r.Name(), raw, err)
						continue // Try next resolver if possible
					}

					if res != nil && res.Title != "" {
						mu.Lock()
						results[raw] = res.Title
						m.cache.Set(raw, res.Title)
						mu.Unlock()
						return
					}
				}
			}
		}(rawURL)
	}

	wg.Wait()
	return results
}

func (m *ResolverManager) ResolveVideoIDs(ctx context.Context, ids []string) map[string]string {
	// For backward compatibility, we convert video IDs to YouTube URLs
	urls := make([]string, len(ids))
	idMap := make(map[string]string)
	for i, id := range ids {
		u := fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
		urls[i] = u
		idMap[u] = id
	}

	urlResults := m.ResolveMulti(ctx, urls)
	results := make(map[string]string)
	for u, title := range urlResults {
		results[idMap[u]] = title
	}
	return results
}
