package binance

import (
	"encoding/json"
	"strconv"
)

// CandleStickData https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#klinecandlestick-data
type CandleStickData struct {
	OpenTime                 float64
	OpenPrice                float64
	HighPrice                float64
	LowPrice                 float64
	ClosePrice               float64
	Volume                   float64
	CloseTime                float64
	QuoteAssetVolume         float64
	NumTrades                float64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}

func (candleStickResponse *CandleStickData) UnmarshalJSON(
	data []byte,
) (err error) {
	var responseSlice []interface{}
	if err := json.Unmarshal(data, &responseSlice); err != nil {
		return err
	}
	// 1501545600000
	candleStickResponse.OpenTime = responseSlice[0].(float64)
	// "4261.48000000"
	candleStickResponse.OpenPrice, err = strconv.ParseFloat(responseSlice[1].(string), 64)
	if err != nil {
		return err
	}
	// "4745.42000000"
	candleStickResponse.HighPrice, err = strconv.ParseFloat(responseSlice[2].(string), 64)
	if err != nil {
		return err
	}
	// "3400.00000000"
	candleStickResponse.LowPrice, err = strconv.ParseFloat(responseSlice[3].(string), 64)
	if err != nil {
		return err
	}
	// "4724.89000000"
	candleStickResponse.ClosePrice, err = strconv.ParseFloat(responseSlice[4].(string), 64)
	if err != nil {
		return err
	}
	// "10015.64027200"
	candleStickResponse.Volume, err = strconv.ParseFloat(responseSlice[5].(string), 64)
	if err != nil {
		return err
	}
	// 1504223999999
	candleStickResponse.CloseTime = responseSlice[6].(float64)
	// "42538297.66482722"
	candleStickResponse.QuoteAssetVolume, err = strconv.ParseFloat(responseSlice[7].(string), 64)
	if err != nil {
		return err
	}
	// 69180
	candleStickResponse.NumTrades = responseSlice[8].(float64)
	// "4610.01943100"
	candleStickResponse.TakerBuyBaseAssetVolume, err = strconv.ParseFloat(responseSlice[9].(string), 64)
	if err != nil {
		return err
	}
	// "19419232.11660334"
	candleStickResponse.TakerBuyQuoteAssetVolume, err = strconv.ParseFloat(responseSlice[10].(string), 64)
	if err != nil {
		return err
	}
	return nil
}

type ErrorResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

// ExchangeInfoResponse https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#exchange-information
type ExchangeInfoResponse struct {
	Timezone   string    `json:"timezone"`
	ServerTime int64     `json:"serverTime"`
	Symbols    []Symbols `json:"symbols"`
}
type Symbols struct {
	Symbol                     string   `json:"symbol"`
	Status                     string   `json:"status"`
	BaseAsset                  string   `json:"baseAsset"`
	BaseAssetPrecision         int      `json:"baseAssetPrecision"`
	QuoteAsset                 string   `json:"quoteAsset"`
	QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
	OrderTypes                 []string `json:"orderTypes"`
	IcebergAllowed             bool     `json:"icebergAllowed"`
	OcoAllowed                 bool     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
	Permissions                []string `json:"permissions"`
}

type Interval string

const (
	Minute Interval = "1m"
	Hour   Interval = "1h"
)
