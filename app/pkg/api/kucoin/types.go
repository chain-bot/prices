package kucoin

import (
	"encoding/json"
	"strconv"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// https://docs.kucoin.com/#get-klines
//////////////////////////////////////////////////////////////////////////////////////////////////////////
type CandleStickResponse struct {
	Code string            `json:"code"`
	Data []CandleStickData `json:"data"`
}
type CandleStickData struct {
	OpenTime   float64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
	Volume     float64
}

func (candleStickResponse *CandleStickData) UnmarshalJSON(
	data []byte,
) (err error) {
	var responseSlice []interface{}
	if err := json.Unmarshal(data, &responseSlice); err != nil {
		return err
	}
	// "1545904980"
	candleStickResponse.OpenTime, err = strconv.ParseFloat(responseSlice[0].(string), 64)
	if err != nil {
		return err
	}
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
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// https://docs.kucoin.com/#get-symbols-list
//////////////////////////////////////////////////////////////////////////////////////////////////////////
type SymbolsResponse struct {
	Code string   `json:"code"`
	Data []Symbol `json:"data"`
}
type Symbol struct {
	Symbol          string `json:"symbol"`
	Name            string `json:"name"`
	BaseCurrency    string `json:"baseCurrency"`
	QuoteCurrency   string `json:"quoteCurrency"`
	BaseMinSize     string `json:"baseMinSize"`
	QuoteMinSize    string `json:"quoteMinSize"`
	BaseMaxSize     string `json:"baseMaxSize"`
	QuoteMaxSize    string `json:"quoteMaxSize"`
	BaseIncrement   string `json:"baseIncrement"`
	QuoteIncrement  string `json:"quoteIncrement"`
	PriceIncrement  string `json:"priceIncrement"`
	FeeCurrency     string `json:"feeCurrency"`
	EnableTrading   bool   `json:"enableTrading"`
	IsMarginEnabled bool   `json:"isMarginEnabled"`
	PriceLimitRate  string `json:"priceLimitRate"`
}
