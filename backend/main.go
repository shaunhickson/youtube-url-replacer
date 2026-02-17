package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sph/youtube-url-replacer/backend/middleware"
	"github.com/sph/youtube-url-replacer/backend/resolvers"
)

func getEnvInt(key string, defaultVal int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

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

	// Initialize Rate Limiter
	rpm := getEnvInt("RATE_LIMIT_RPM", 60)
	burst := getEnvInt("RATE_LIMIT_BURST", 20)
	rateLimiter := middleware.NewRateLimiter(rpm, burst)
	// Clean up old visitors every minute, expire after 3 minutes
	rateLimiter.CleanupBackground(1*time.Minute, 3*time.Minute)

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

	// Configure Timeout
	if timeoutStr := os.Getenv("RESOLVER_TIMEOUT_MS"); timeoutStr != "" {
		if ms, err := strconv.Atoi(timeoutStr); err == nil {
			manager.SetTimeout(time.Duration(ms) * time.Millisecond)
		}
	}

	enabledResolvers := os.Getenv("ENABLED_RESOLVERS")
	isEnabled := func(name string) bool {
		if enabledResolvers == "" {
			return true
		}
		for _, r := range strings.Split(enabledResolvers, ",") {
			if strings.TrimSpace(r) == name {
				return true
			}
		}
		return false
	}

	// Register YouTube Resolver
	if isEnabled("youtube") {
		ytResolver, err := resolvers.NewYouTubeResolver(apiKey)
		if err != nil {
			log.Fatalf("Failed to create YouTube resolver: %v", err)
		}
		manager.Register(ytResolver)
	}

	// Register OpenGraph Resolver (Fallback)
	if isEnabled("opengraph") {
		manager.Register(resolvers.NewOpenGraphResolver())
	}

	handler := NewHandler(cache, manager)
	handler.MaxItems = getEnvInt("MAX_ITEMS_PER_REQUEST", 50)
	handler.MaxBodyBytes = int64(getEnvInt("MAX_BODY_BYTES", 10240))

	// Set up routes
	http.Handle("/resolve", rateLimiter.Middleware(handler))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Health check write failed: %v", err)
		}
	})

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}