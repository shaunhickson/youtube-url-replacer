package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResolveRequest struct {
	VideoIDs []string `json:"videoIds"`
}

type ResolveResponse struct {
	Titles map[string]string `json:"titles"`
}

type Handler struct {
	cache     Cache
	ytService *YouTubeService
}

func NewHandler(cache Cache, ytService *YouTubeService) *Handler {
	return &Handler{
		cache:     cache,
		ytService: ytService,
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

	if len(req.VideoIDs) == 0 {
		json.NewEncoder(w).Encode(ResolveResponse{Titles: map[string]string{}})
		return
	}

	// 1. Check Cache
	cachedTitles := h.cache.GetMulti(req.VideoIDs)
	
	// 2. Identify missing IDs
	var missingIDs []string
	for _, id := range req.VideoIDs {
		if _, ok := cachedTitles[id]; !ok {
			missingIDs = append(missingIDs, id)
		}
	}

	// 3. Fetch from YouTube if needed
	if len(missingIDs) > 0 {
		log.Printf("Fetching %d videos from YouTube API", len(missingIDs))
		fetchedTitles, err := h.ytService.FetchTitles(missingIDs)
		if err != nil {
			log.Printf("Error fetching from YouTube: %v", err)
			// Return partial results (cached only) or error? 
			// Let's return what we have, maybe the UI can retry.
		} else {
			// Update cache and results
			for id, title := range fetchedTitles {
				h.cache.Set(id, title)
				cachedTitles[id] = title
			}
		}
	}

	// 4. Return combined results
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ResolveResponse{Titles: cachedTitles}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
