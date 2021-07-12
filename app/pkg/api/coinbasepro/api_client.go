package coinbasepro

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/chain-bot/prices/app/pkg/api/common"
	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
			log.Printf("ERROR WAITING FOR LIMIT: %s\n", err.Error())
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

func (apiClient *ApiClient) getProducts() (productsResponse ProductsResponse, err error) {
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

func (apiClient *ApiClient) sendAPIKeyAuthenticatedGetRequest(
	urlString string, headers ...header,
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
