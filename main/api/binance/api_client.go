package binance

import (
	"encoding/json"
	"fmt"
	"github.com/mochahub/coinprice-scraper/main/api/common"
	"io/ioutil"
	"net/http"
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
		Client: common.NewHTTPClient(
			maxRetries,
			time.Duration(rateLimit)),
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
	urlString := fmt.Sprintf("%s%s?symbol=%s%s&interval=%s&startTime=%d&endTime=%d&limit=%d",
		baseUrl,
		getCandleStick,
		baseSymbol,
		quoteSymbol,
		interval,
		startTime.UTC().Unix()*1000,
		endTime.UTC().Unix()*1000,
		maxLimit,
	)
	httpReq, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	retryableRequest, err := retryablehttp.FromRequest(httpReq)
	if err != nil {
		return nil, err
	}
	retryableRequest.Header.Add("X-MBX-APIKEY", apiClient.apiKey)
	httpResp, err := apiClient.Do(retryableRequest)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err = json.Unmarshal(body, &candleStickResponse); err != nil {
		return nil, err
	}
	return candleStickResponse, nil
}

// Get CandleStick data from [startTime, endTime]
func (apiClient *apiClient) getExchangeInfo() (candleStickResponse *ExchangeInfoResponse, err error) {
	urlString := fmt.Sprintf("%s%s", baseUrl, getExchangeInfo)
	httpReq, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	retryableRequest, err := retryablehttp.FromRequest(httpReq)
	if err != nil {
		return nil, err
	}
	retryableRequest.Header.Add("X-MBX-APIKEY", apiClient.apiKey)
	httpResp, err := apiClient.Do(retryableRequest)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err = json.Unmarshal(body, &candleStickResponse); err != nil {
		return nil, err
	}
	return candleStickResponse, nil
}
