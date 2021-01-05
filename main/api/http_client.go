package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func NewHTTPClient(maxRetries int, callsPerSecond time.Duration) *retryablehttp.Client {
	rateLimiter := NewRateLimiter(callsPerSecond)
	httpClient := retryablehttp.NewClient()
	httpClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp.StatusCode >= http.StatusBadRequest {
			return true, fmt.Errorf("retry eror: StatusCode %d Error %s", resp.StatusCode, err.Error())
		}
		rateLimiter.RateLimitCall()
		return false, nil
	}
	httpClient.RetryMax = maxRetries
	return httpClient
}
