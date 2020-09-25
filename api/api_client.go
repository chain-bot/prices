package api

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	rateLimiter RateLimiter
	maxRetries  int
	httpClient  *retryablehttp.Client
}

func NewAPIClient(maxRetries int, callsPerSecond time.Duration) Client {
	rateLimiter := NewRateLimiter(callsPerSecond)
	httpClient := retryablehttp.Client{
		CheckRetry: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			rateLimiter.RateLimitCall()
			return true, nil
		},
	}
	return Client{
		rateLimiter,
		maxRetries,
		&httpClient,
	}
}
