package main

import (
	"sync"
)

// Cache defines the interface for storing and retrieving video titles
type Cache interface {
	Get(videoID string) (string, bool)
	Set(videoID string, title string)
	GetMulti(videoIDs []string) map[string]string
}

// InMemoryCache is a simple thread-safe map implementation of Cache
type InMemoryCache struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		store: make(map[string]string),
	}
}

func (c *InMemoryCache) Get(videoID string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.store[videoID]
	return val, ok
}

func (c *InMemoryCache) Set(videoID string, title string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[videoID] = title
}

func (c *InMemoryCache) GetMulti(videoIDs []string) map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	results := make(map[string]string)
	for _, id := range videoIDs {
		if val, ok := c.store[id]; ok {
			results[id] = val
		}
	}
	return results
}
