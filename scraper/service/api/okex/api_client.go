package okex

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/common"
	"golang.org/x/time/rate"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
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
			log.Printf("ERROR WAITING FOR LIMIT: %s\n", err.Error())
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
	if endTime.IsZero() {
		endTime = time.Now()
	}
	uri := fmt.Sprintf(getCandles, symbol)
	urlString := fmt.Sprintf("%s%s?start=%s&granularity=%d&limit=%d",
		baseUrl,
		uri,
		startTime.UTC().Format(time.RFC3339),
		int(interval.Seconds()),
		limit,
	)
	println(urlString)
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
