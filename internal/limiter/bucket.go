package limiter

import "time"

type TokenBucket struct {
	tokens         float64
	lastRefillTime time.Time
	rate           float64
	capacity       float64
	lastReqTime    time.Time
}
