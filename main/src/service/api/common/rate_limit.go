package common

import (
	"time"
)

// TODO(Zahin): Replace with https://pkg.go.dev/golang.org/x/time/rate
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
