package api

import (
	"time"
)

type RateLimiter struct {
	throttle <-chan time.Time
}

func NewRateLimiter(callsPerSecond time.Duration) RateLimiter {
	rateLimit := time.Second / callsPerSecond
	return RateLimiter{
		throttle: time.Tick(rateLimit),
	}
}
func (rateLimiter RateLimiter) RateLimitCall() {
	<-rateLimiter.throttle
}
