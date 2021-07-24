package okex

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chain-bot/prices/app/pkg/api/common"
	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type ApiClient struct {
	*retryablehttp.Client
	*rate.Limiter
}

func NewOkexAPIClient() *ApiClient {
	rateLimiter := rate.NewLimiter(rate.Every(time.Second/callsPerSecond), 1)
	httpClient := retryablehttp.NewClient()
	httpClient.CheckRetry = common.DefaultCheckRetry
	httpClient.RetryWaitMin = common.DefaultRetryMin
	httpClient.RetryMax = common.MaxRetries
	apiClient := ApiClient{
		Client:  httpClient,
		Limiter: rateLimiter,
	}
	apiClient.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, retry int) {
		if err := apiClient.Limiter.Wait(context.Background()); err != nil {
			log.WithField("err", err.Error()).Errorf("waiting for rate limit")
			return
		}
	}
	return &apiClient
}

func (apiClient *ApiClient) GetExchangeIdentifier() string {
	return OKEX
}

//Get CandleStick data from [startTime, endTime)
func (apiClient *ApiClient) getInstrumentCandles(
	symbol string,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
	limit int64,
) (candleStickData []*CandleStickData, err error) {
	uri := fmt.Sprintf(getCandles, symbol)
	// For okex end < start...
	urlString := fmt.Sprintf("%s%s?granularity=%d&limit=%d&start=%s&end=%s",
		baseUrl,
		uri,
		int(interval.Seconds()),
		limit,
		endTime.Add(-interval).UTC().Format(time.RFC3339),
		startTime.UTC().Format(time.RFC3339),
	)
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &candleStickData); err != nil {
		return nil, err
	}
	return candleStickData, nil
}

func (apiClient *ApiClient) getInstruments() (instruments []*Instrument, err error) {
	urlString := fmt.Sprintf("%s%s", baseUrl, getInstruments)
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &instruments); err != nil {
		return nil, err
	}
	return instruments, nil
}

func (apiClient *ApiClient) intervalQueryParamFromDuration(intervalDuration time.Duration) (interval string) {
	return fmt.Sprintf("%dmin", int(intervalDuration.Minutes()))
}

func (apiClient *ApiClient) sendUnAuthenticatedGetRequest(
	urlString string,
) (*http.Response, error) {
	httpReq, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	retryableRequest, err := retryablehttp.FromRequest(httpReq)
	if err != nil {
		return nil, err
	}
	return apiClient.Do(retryableRequest)
}
