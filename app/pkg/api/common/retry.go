package common

import (
	"context"
	"log"
	"net/http"
	"time"
)

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
