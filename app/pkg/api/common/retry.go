package common

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultRetryMin = time.Second * 10
	MaxRetries      = 3
)

func DefaultCheckRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if resp.StatusCode >= http.StatusTooManyRequests {
		if err != nil {
			log.WithFields(
				log.Fields{
					"err":        err.Error(),
					"statusCode": resp.StatusCode,
				}).Errorf("waitinf for rate limit")
		} else {
			log.WithFields(
				log.Fields{
					"statusCode": resp.StatusCode,
				}).Errorf("waitinf for rate limit")
		}
		return true, err
	}
	return false, nil
}
