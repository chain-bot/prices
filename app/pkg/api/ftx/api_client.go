package ftx

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chain-bot/prices/app/pkg/api/common"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/hashicorp/go-retryablehttp"
)

type ApiClient struct {
	*retryablehttp.Client
	*rate.Limiter
}

func NewFtxAPIClient() *ApiClient {
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
	return FTX
}

//Get CandleStick data from [startTime, endTime)
func (apiClient *ApiClient) getHistoricalPrices(
	market string,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
	limit int64,
) (historicalPricesResponse *HistoricalPricesResponse, err error) {
	if endTime.IsZero() {
		endTime = time.Now()
	}
	uri := fmt.Sprintf(getHistoricalPrices, market, int(interval.Seconds()), limit, startTime.Unix(), endTime.Unix())
	urlString := fmt.Sprintf("%s%s", baseUrl, uri)
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &historicalPricesResponse); err != nil {
		return nil, err
	}
	return historicalPricesResponse, nil
}

func (apiClient *ApiClient) getMarkets() (marketResponse *MarketsResponse, err error) {
	urlString := fmt.Sprintf("%s%s", baseUrl, getMarkets)
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &marketResponse); err != nil {
		return nil, err
	}
	return marketResponse, nil
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
