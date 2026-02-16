package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sph/youtube-url-replacer/backend/resolvers"
)

type ResolveRequest struct {
	VideoIDs []string `json:"videoIds"`
	URLs     []string `json:"urls"`
}

type ResolveResponse struct {
	Titles map[string]string `json:"titles"`
}

type Handler struct {
	cache    resolvers.Cache
	manager  *resolvers.ResolverManager
}

func NewHandler(cache resolvers.Cache, manager *resolvers.ResolverManager) *Handler {
	return &Handler{
		cache:   cache,
		manager: manager,
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

	var req ResolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.VideoIDs) == 0 && len(req.URLs) == 0 {
		json.NewEncoder(w).Encode(ResolveResponse{Titles: map[string]string{}})
		return
	}

	results := make(map[string]string)

	// 1. Resolve URLs
	if len(req.URLs) > 0 {
		urlResults := h.manager.ResolveMulti(r.Context(), req.URLs)
		for u, title := range urlResults {
			results[u] = title
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
	if err := json.NewEncoder(w).Encode(ResolveResponse{Titles: results}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
