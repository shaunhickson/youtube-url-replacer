package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sph/youtube-url-replacer/backend/resolvers"
)

type ResolveRequest struct {
	VideoIDs []string `json:"videoIds"`
	URLs     []string `json:"urls"`
}

type ResolveResponse struct {
	Titles  map[string]string            `json:"titles"`
	Details map[string]*resolvers.Result `json:"details,omitempty"`
}

type Handler struct {
	cache        resolvers.Cache
	manager      *resolvers.ResolverManager
	MaxBodyBytes int64
	MaxItems     int
}

func NewHandler(cache resolvers.Cache, manager *resolvers.ResolverManager) *Handler {
	return &Handler{
		cache:        cache,
		manager:      manager,
		MaxBodyBytes: 10 * 1024, // Default 10KB
		MaxItems:     50,        // Default 50 items
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for development
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Limit Body Size
	r.Body = http.MaxBytesReader(w, r.Body, h.MaxBodyBytes)

	var req ResolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}
		return
	}

	// 2. Limit Item Count
	totalItems := len(req.VideoIDs) + len(req.URLs)
	if totalItems > h.MaxItems {
		http.Error(w, "Too many items in request", http.StatusRequestEntityTooLarge)
		return
	}

	if totalItems == 0 {
		if err := json.NewEncoder(w).Encode(ResolveResponse{Titles: map[string]string{}}); err != nil {
			slog.Error("Error encoding empty response", "error", err)
		}
		return
	}

	results := make(map[string]string)
	details := make(map[string]*resolvers.Result)

	// 1. Resolve URLs
	if len(req.URLs) > 0 {
		urlResults := h.manager.ResolveMulti(r.Context(), req.URLs)
		for u, res := range urlResults {
			results[u] = res.Title
			details[u] = res
		}
	}

	// 2. Resolve Video IDs (Legacy)
	if len(req.VideoIDs) > 0 {
		idResults := h.manager.ResolveVideoIDs(r.Context(), req.VideoIDs)
		for id, title := range idResults {
			results[id] = title
		}
	}

	// 3. Return combined results
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ResolveResponse{
		Titles:  results,
		Details: details,
	}); err != nil {
		slog.Error("Error encoding response", "error", err)
	}
}
