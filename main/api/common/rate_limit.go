package common

import (
	"time"
)

// TODO(Zahin): No clue if this even works
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
