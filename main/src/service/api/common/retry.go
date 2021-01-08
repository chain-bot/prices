package common

import (
	"context"
	"log"
	"net/http"
	"time"
)

//func NewHTTPClient(
//	maxRetries int,
//	callsPerSecond time.Duration,
//) *retryablehttp.Client {
//	rateLimiter := NewRateLimiter(callsPerSecond)
//	httpClient := retryablehttp.NewClient()
//	httpClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
//		if resp.StatusCode >= http.StatusBadRequest {
//			resp.Body.Close()
//			return true, fmt.Errorf("retry eror: StatusCode %d Error %s", resp.StatusCode, err.Error())
//		}
//
//		rateLimiter.RateLimitCall()
//		return true, nil
//	}
//	//type RequestLogHook func(Logger, *http.Request, int)
//	httpClient.RetryWaitMin = time.Second * 10
//	httpClient.RetryMax = maxRetries
//	return httpClient
//}
const (
	DefaultRetryMin = time.Second * 10
	MaxRetries      = 3
)

func DefaultCheckRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if resp.StatusCode >= http.StatusBadRequest {
		if err != nil {
			log.Printf("retry error: StatusCode %d Error %s\n", resp.StatusCode, err.Error())
		} else {
			log.Printf("retry error: StatusCode %d\n", resp.StatusCode)
		}
		return true, err
	}
	return false, nil
}
