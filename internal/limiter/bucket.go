package limiter

import "time"

type TokenBucket struct {
	tokens         float64
	lastRefillTime time.Time
	rate           float64
	capacity       float64
	lastReqTime    time.Time
}

func (t *TokenBucket) Allow() bool {

	now := time.Now()

	sinceLastReqTime := time.Since(t.lastReqTime).Seconds()
	tokenCap := sinceLastReqTime * t.rate

	if tokenCap > t.capacity {
		t.tokens = t.capacity
	}

	t.lastRefillTime = now

	if t.tokens >= 1 {
		t.tokens -= 1
		t.lastReqTime = now
		return true
	}

	t.lastReqTime = now
	return false
}
