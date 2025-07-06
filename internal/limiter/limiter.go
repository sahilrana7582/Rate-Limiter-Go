package limiter

import (
	"sync"
	"time"
)

type Limiter struct {
	buckets map[string]*TokenBucket
	mu      sync.Mutex
}

func (l *Limiter) Allow(ip string) bool {

	tokenBucket, exists := l.buckets[ip]
	if !exists {
		buket := &TokenBucket{
			tokens:         1,
			lastRefillTime: time.Now(),
			rate:           1,
			capacity:       10,
		}

		tokenBucket = buket

		l.buckets[ip] = tokenBucket

	}

	return tokenBucket.Allow()
}
