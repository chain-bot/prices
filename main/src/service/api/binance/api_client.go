package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mochahub/coinprice-scraper/main/src/service/api/common"
	"golang.org/x/time/rate"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type apiClient struct {
	*retryablehttp.Client
	*rate.Limiter
	apiKey string
}

func NewBinanceAPIClient(
	apiKey string,
) *apiClient {
	// 1200 callsPerMinute:(60*1000)/1200
	rateLimiter := rate.NewLimiter(rate.Every(50*time.Millisecond), 2)
	httpClient := retryablehttp.NewClient()
	httpClient.CheckRetry = common.DefaultCheckRetry
	httpClient.RetryWaitMin = common.DefaultRetryMin
	httpClient.RetryMax = common.MaxRetries
	apiClient := apiClient{
		Client:  httpClient,
		Limiter: rateLimiter,
		apiKey:  apiKey,
	}
	apiClient.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, retry int) {
		if err := apiClient.Limiter.Wait(context.Background()); err != nil {
			log.Printf("ERROR WAITING FOR LIMIT: %s\n", err.Error())
			return
		}
	}
	return &apiClient
}

func (apiClient *apiClient) GetExchangeIdentifier() string {
	return BINANCE
}

// Get CandleStick data from [startTime, endTime]
func (apiClient *apiClient) getCandleStickData(
	baseSymbol string,
	quoteSymbol string,
	interval common.Interval,
	startTime time.Time,
	endTime time.Time,
) (candleStickResponse []*CandleStickData, err error) {
	if endTime.IsZero() {
		endTime = time.Now()
	}
	params := url.Values{}
	params.Add("symbol", baseSymbol+quoteSymbol)
	params.Add("interval", string(interval))
	params.Add("startTime", strconv.FormatInt(UnixMillis(startTime), 10))
	params.Add("endTime", strconv.FormatInt(UnixMillis(endTime), 10))
	params.Add("limit", strconv.Itoa(maxLimit))
	urlString := fmt.Sprintf("%s%s?%s", baseUrl, getCandleStick, params.Encode())
	resp, err := apiClient.sendAPIKeyAuthenticatedGetRequest(urlString)
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

// Get ExchangeInfo (supported pairs, percision, etc)
func (apiClient *apiClient) getExchangeInfo() (exchangeInfoResponse *ExchangeInfoResponse, err error) {
	urlString := fmt.Sprintf("%s%s", baseUrl, getExchangeInfo)
	resp, err := apiClient.sendAPIKeyAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &exchangeInfoResponse); err != nil {
		return nil, err
	}
	return exchangeInfoResponse, nil
}

func (apiClient *apiClient) sendAPIKeyAuthenticatedGetRequest(
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
	retryableRequest.Header.Add("X-MBX-APIKEY", apiClient.apiKey)
	return apiClient.Do(retryableRequest)
}
