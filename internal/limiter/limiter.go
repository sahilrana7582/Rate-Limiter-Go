package limiter

import (
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	buckets        map[string]*TokenBucket
	mu             sync.Mutex
	rate           float64
	capacity       float64
	maxTimeAllowed time.Duration
	panicRecover   chan struct{}
}

func NewLimiter(rate, capacity float64, maxTimeAllowed time.Duration) *Limiter {
	return &Limiter{
		buckets:        make(map[string]*TokenBucket),
		rate:           rate,
		capacity:       capacity,
		maxTimeAllowed: maxTimeAllowed,
		panicRecover:   make(chan struct{}),
	}
}

func (l *Limiter) Allow(ip string) bool {
	l.mu.Lock()
	tokenBucket, exists := l.buckets[ip]

	if !exists {
		now := time.Now()
		tokenBucket = &TokenBucket{
			tokens:         l.capacity,
			lastRefillTime: now,
			rate:           l.rate,
			capacity:       l.capacity,
			LastReqTime:    now,
		}
		l.buckets[ip] = tokenBucket
	}
	l.mu.Unlock()

	return tokenBucket.Allow()
}

func (l *Limiter) WatchCleaner() {
	go func() {
		for {
			<-l.panicRecover
			fmt.Println("Restarting cleaner after panic...")
			l.StartCleaning()
		}
	}()
}

func (l *Limiter) StartCleaning() {
	go func() {
		fmt.Println("Cleaner Running")

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic in StartCleaning:", r)
				l.panicRecover <- struct{}{}
			}
		}()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			func() {
				l.mu.Lock()
				defer l.mu.Unlock()

				for ip, tokenB := range l.buckets {

					if time.Since(tokenB.LastReqTime) > l.maxTimeAllowed {
						fmt.Printf("Cleaning up IP: %s (inactive for %.0f seconds)\n", ip, time.Since(tokenB.LastReqTime).Seconds())
						delete(l.buckets, ip)
					}
				}
			}()

			time.Sleep(1 * time.Second)
		}
	}()
}
