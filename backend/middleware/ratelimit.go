package middleware

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	ips    sync.Map
	limit  rate.Limit
	burst  int
	mu     sync.Mutex
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(rpm int, burst int) *RateLimiter {
	return &RateLimiter{
		limit: rate.Limit(rpm) / 60.0,
		burst: burst,
	}
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.ips.Load(ip)
	if !exists {
		limiter := rate.NewLimiter(rl.limit, rl.burst)
		rl.ips.Store(ip, &visitor{limiter: limiter, lastSeen: time.Now()})
		return limiter
	}

	vis := v.(*visitor)
	vis.lastSeen = time.Now()
	return vis.limiter
}

// CleanupBackground starts a goroutine to remove old entries
func (rl *RateLimiter) CleanupBackground(interval time.Duration, expiry time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			rl.ips.Range(func(key, value interface{}) bool {
				v := value.(*visitor)
				if time.Since(v.lastSeen) > expiry {
					rl.ips.Delete(key)
				}
				return true
			})
		}
	}()
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			w.Header().Set("Retry-After", "60") // Simple retry advice
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "Too many requests. Please try again later.",
			}); err != nil {
				log.Printf("Failed to encode 429 response: %v", err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	// 1. Check X-Forwarded-For (Cloud Run / Load Balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Can be "client_ip, proxy1, proxy2"
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// 2. Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
