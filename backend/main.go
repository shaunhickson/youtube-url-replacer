package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	var ytService *YouTubeService
	var err error

	if apiKey == "" {
		log.Println("Warning: YOUTUBE_API_KEY is not set. Using MOCK YouTube service.")
		ytService = &YouTubeService{service: nil} // Nil service indicates mock mode
	} else {
		ytService, err = NewYouTubeService(apiKey)
		if err != nil {
			log.Fatalf("Failed to create YouTube service: %v", err)
		}
	}

	// Initialize Cache (Firestore or Memory)
	var cache Cache
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID != "" {
		log.Printf("Initializing Firestore Cache for project: %s", projectID)
		fsCache, err := NewFirestoreCache(projectID)
		if err != nil {
			log.Fatalf("Failed to initialize Firestore: %v", err)
		}
		cache = fsCache
	} else {
		log.Println("Initializing In-Memory Cache (Non-persistent)")
		cache = NewInMemoryCache()
	}

	handler := NewHandler(cache, ytService)

	// Set up routes
	http.Handle("/resolve", handler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}