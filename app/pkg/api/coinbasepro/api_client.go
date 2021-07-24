package coinbasepro

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

func NewCoinbaseProAPIClient() *ApiClient {
	// 3 callsPerSecond
	rateLimiter := rate.NewLimiter(rate.Every(time.Second/3), 6)
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
	return COINBASEPRO
}

// Get CandleStick data from [startTime, endTime]
func (apiClient *ApiClient) getCandleStickData(
	granularity int,
	startTime time.Time,
	endTime time.Time,
	productID string,
) (candleStickResponse []*CandleStickData, err error) {
	params := url.Values{}
	params.Add("start", startTime.Format(time.RFC3339))
	params.Add("end", endTime.Format(time.RFC3339))
	params.Add("granularity", strconv.Itoa(granularity))
	path := fmt.Sprintf(getCandles, productID)
	urlString := fmt.Sprintf("%s%s?%s", baseURL, path, params.Encode())
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &candleStickResponse); err != nil {
		return nil, err
	}
	return candleStickResponse, nil
}

func (apiClient *ApiClient) getProducts() (productsResponse ProductsResponse, err error) {
	urlString := fmt.Sprintf("%s%s", baseURL, getExchangeProducts)
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &productsResponse); err != nil {
		return nil, err
	}
	return productsResponse, nil
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
