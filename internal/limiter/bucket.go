package limiter

import "time"

type TokenBucket struct {
	tokens         float64
	lastRefillTime time.Time
	rate           float64
	capacity       float64
	LastReqTime    time.Time
}

func (t *TokenBucket) Allow() bool {
	now := time.Now()

	elapsed := now.Sub(t.lastRefillTime).Seconds()
	t.tokens += elapsed * t.rate
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}

	t.lastRefillTime = now

	if t.tokens >= 1 {
		t.tokens -= 1
		t.LastReqTime = now
		return true
	}

	t.LastReqTime = now
	return false
}
