package kucoin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chain-bot/prices/app/pkg/api/common"
	"golang.org/x/time/rate"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type ApiClient struct {
	*retryablehttp.Client
	*rate.Limiter
}

func NewKucoinAPIClient() *ApiClient {
	rateLimiter := rate.NewLimiter(rate.Every(time.Minute/callsPerMinute), 1)
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
	return KUCOIN
}

//Get CandleStick data from [startTime, endTime)
func (apiClient *ApiClient) getKlines(
	symbol string,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) (candleStickResponse *CandleStickResponse, err error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("startAt", strconv.FormatInt(startTime.Unix(), 10))
	params.Add("endAt", strconv.FormatInt(endTime.Unix(), 10))
	params.Add("type", apiClient.intervalQueryParamFromDuration(interval))
	urlString := fmt.Sprintf("%s%s?%s", baseUrl, getKlines, params.Encode())
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

// Get ExchangeInfo (supported pairs, precision, etc)
func (apiClient *ApiClient) getSymbols() (symbolsResponse *SymbolsResponse, err error) {
	urlString := fmt.Sprintf("%s%s", baseUrl, getSymbols)
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &symbolsResponse); err != nil {
		return nil, err
	}
	return symbolsResponse, nil
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

func (apiClient *ApiClient) intervalQueryParamFromDuration(intervalDuration time.Duration) (interval string) {
	return fmt.Sprintf("%dmin", int(intervalDuration.Minutes()))
}
