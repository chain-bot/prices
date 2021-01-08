package coinbasepro

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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
	apiKey        string
	apiSecret     string
	apiPassPhrase string
}

func NewCoinbaseProAPIClient(
	apiKey, apiSecret, apiPassPhrase string,
) *apiClient {
	// 3 callsPerSecond
	rateLimiter := rate.NewLimiter(rate.Every(time.Second/3), 6)
	httpClient := retryablehttp.NewClient()
	httpClient.CheckRetry = common.DefaultCheckRetry
	httpClient.RetryWaitMin = common.DefaultRetryMin
	httpClient.RetryMax = common.MaxRetries
	apiClient := apiClient{
		Client:        httpClient,
		Limiter:       rateLimiter,
		apiKey:        apiKey,
		apiSecret:     apiSecret,
		apiPassPhrase: apiPassPhrase,
	}
	apiClient.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, retry int) {
		if err := apiClient.Limiter.Wait(context.Background()); err != nil {
			log.Printf("ERROR WAITING FOR LIMIT: %s\n", err.Error())
			return
		}
		apiClient.writeRequestHeaders(req)
	}
	return &apiClient
}
func (apiClient *apiClient) GetExchangeIdentifier() string {
	return COINBASE
}

// Get CandleStick data from [startTime, endTime]
func (apiClient *apiClient) getCandleStickData(
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
	resp, err := apiClient.sendAPIKeyAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &candleStickResponse); err != nil {
		log.Println(string(body))
		return nil, err
	}
	return candleStickResponse, nil
}

func (apiClient *apiClient) getProducts() (productsResponse ProductsResponse, err error) {
	urlString := fmt.Sprintf("%s%s", baseURL, getExchangeProducts)
	resp, err := apiClient.sendAPIKeyAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &productsResponse); err != nil {
		log.Println(string(body))
		return nil, err
	}
	return productsResponse, nil
}

func (apiClient *apiClient) sendAPIKeyAuthenticatedGetRequest(
	urlString string, headers ...header,
) (*http.Response, error) {
	httpReq, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	apiClient.writeRequestHeaders(httpReq, headers...)
	retryableRequest, err := retryablehttp.FromRequest(httpReq)
	if err != nil {
		return nil, err
	}
	return apiClient.Do(retryableRequest)
}

func (apiClient *apiClient) writeRequestHeaders(
	req *http.Request, headers ...header,
) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + req.Method + req.URL.Path
	sig, err := generateSig(message, apiClient.apiSecret)
	if err != nil {
		log.Printf("FAILED TO CREATE %s SIGNATURE: %s\n", apiClient.GetExchangeIdentifier(), err.Error())
		return
	}
	for i := range headers {
		req.Header.Set(headers[i].key, headers[i].value)
	}
	req.Header.Set("CB-ACCESS-KEY", apiClient.apiKey)
	req.Header.Set("CB-ACCESS-PASSPHRASE", apiClient.apiPassPhrase)
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("CB-ACCESS-SIGN", sig)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}

// Credit to https://github.com/preichenberger/go-coinbasepro
func generateSig(message, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	signature := hmac.New(sha256.New, key)
	_, err = signature.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}
