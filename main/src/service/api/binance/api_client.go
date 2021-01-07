package binance

import (
	"encoding/json"
	"fmt"
	"github.com/mochahub/coinprice-scraper/main/src/service/api/common"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type apiClient struct {
	*retryablehttp.Client
	apiKey string
}

func NewBinanceAPIClient(
	apiKey string,
) *apiClient {
	return &apiClient{
		Client: common.NewHTTPClient(maxRetries, time.Duration(rateLimit)),
		apiKey: apiKey,
	}
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
