package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sph/youtube-url-replacer/backend/resolvers"
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

	// Initialize Cache (Firestore or Memory)
	var cache resolvers.Cache
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

	// Initialize Resolver Manager
	manager := resolvers.NewResolverManager(cache)

	// Register YouTube Resolver
	ytResolver, err := resolvers.NewYouTubeResolver(apiKey)
	if err != nil {
		log.Fatalf("Failed to create YouTube resolver: %v", err)
	}
	manager.Register(ytResolver)

	// Register OpenGraph Resolver (Fallback)
	manager.Register(resolvers.NewOpenGraphResolver())

	handler := NewHandler(cache, manager)

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