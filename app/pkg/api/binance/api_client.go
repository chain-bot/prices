package binance

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
	"github.com/chain-bot/prices/app/utils"
	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type ApiClient struct {
	*retryablehttp.Client
	*rate.Limiter
}

func NewBinanceAPIClient() *ApiClient {
	// 1200 callsPerMinute:(60*1000)/1200
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
			log.WithField("err", err.Error()).Errorf("waiting for rate limit")
			return
		}
	}
	return &apiClient
}

func (apiClient *ApiClient) GetExchangeIdentifier() string {
	return BINANCE
}

// Get CandleStick data from [startTime, endTime]
func (apiClient *ApiClient) getCandleStickData(
	symbol string,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) (candleStickResponse []*CandleStickData, err error) {
	if endTime.IsZero() {
		endTime = time.Now()
	}
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("interval", string(apiClient.getBinanceIntervalFromDuration(interval)))
	params.Add("startTime", strconv.FormatInt(utils.UnixMillis(startTime), 10))
	params.Add("endTime", strconv.FormatInt(utils.UnixMillis(endTime), 10))
	params.Add("limit", strconv.Itoa(maxLimit))
	urlString := fmt.Sprintf("%s%s?%s", baseUrl, getCandleStick, params.Encode())
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &candleStickResponse); err != nil {
		var errResponse ErrorResponse
		newErr := json.Unmarshal(body, &errResponse)
		if newErr != nil {
			return nil, fmt.Errorf("err1=%s, err2=%s", err.Error(), newErr.Error())
		} else {
			return nil, fmt.Errorf("err1=%s, code=%d, msg=%s",
				err.Error(), errResponse.Code, errResponse.Msg)
		}
	}
	return candleStickResponse, nil
}

// Get ExchangeInfo (supported pairs, precision, etc)
func (apiClient *ApiClient) getExchangeInfo() (exchangeInfoResponse *ExchangeInfoResponse, err error) {
	urlString := fmt.Sprintf("%s%s", baseUrl, getExchangeInfo)
	resp, err := apiClient.sendUnAuthenticatedGetRequest(urlString)
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

func (apiClient *ApiClient) getBinanceIntervalFromDuration(
	interval time.Duration,
) Interval {
	ret := Minute
	if int(interval.Hours()) > 0 {
		ret = Hour
	}
	return ret
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
