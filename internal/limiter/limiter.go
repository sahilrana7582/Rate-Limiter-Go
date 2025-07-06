package limiter

import (
	"sync"
	"time"
)

type Limiter struct {
	buckets  map[string]*TokenBucket
	mu       sync.Mutex
	rate     float64
	capacity float64
}

func NewLimiter(rate, capacity float64) *Limiter {
	return &Limiter{
		buckets:  make(map[string]*TokenBucket),
		rate:     rate,
		capacity: capacity,
	}
}

func (l *Limiter) Allow(ip string) bool {

	tokenBucket, exists := l.buckets[ip]
	if !exists {
		now := time.Now()
		buket := &TokenBucket{
			tokens:         1,
			lastRefillTime: now,
			rate:           1,
			capacity:       10,
			lastReqTime:    now,
		}

		tokenBucket = buket
		l.mu.Lock()
		l.buckets[ip] = tokenBucket
		l.mu.Unlock()
	}
	return tokenBucket.Allow()
}

// func (l *Limiter) StartCleaning() {

// }
