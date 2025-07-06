package middleware

import (
	"net/http"
	"strings"

	"github.com/sahilrana7582/rate-limiter-go/internal/limiter"
)

type RateLimiterMiddleware struct {
	limiter *limiter.Limiter
}

func NewRateLimiter(l *limiter.Limiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter: l,
	}
}

func (rl *RateLimiterMiddleware) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)

		if !rl.limiter.Allow(ip) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded. Try again later."))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}

	ipPort := strings.Split(r.RemoteAddr, ":")
	return ipPort[0]
}
