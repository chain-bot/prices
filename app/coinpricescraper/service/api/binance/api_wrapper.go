package binance

import (
	"fmt"
	"github.com/mochahub/coinprice-scraper/app/coinpricescraper/service/api/common"
	"strings"
	"time"
)

// Get CandleStick data from [startTime, endTime] with pagination
func (apiClient *ApiClient) GetAllOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval common.Interval,
	startTime time.Time,
	endTime time.Time,
) ([]*common.OHLCMarketData, error) {
	var durationFromInterval time.Duration
	switch interval {
	case common.Day:
		durationFromInterval = time.Hour * 24
	case common.Hour:
		durationFromInterval = time.Hour
	case common.Minute:
		durationFromInterval = time.Minute
	default:
		return nil, fmt.Errorf("unknown interval: %s", interval)
	}
	if endTime.IsZero() {
		endTime = time.Now()
	}
	result := []*common.OHLCMarketData{}
	for startTime.Before(endTime) {
		newEndTime := startTime.Add(maxLimit * durationFromInterval)
		if newEndTime.After(endTime) {
			newEndTime = endTime
		}
		ohlcMarketData, err := apiClient.GetOHLCMarketData(
			baseSymbol,
			quoteSymbol,
			interval,
			startTime,
			newEndTime)
		if err != nil {
			return nil, err
		}
		result = append(result, ohlcMarketData...)
		startTime = newEndTime
	}
	return result, nil
}

func (apiClient *ApiClient) GetSupportedPairs() ([]*common.Symbol, error) {
	exchangeInfo, err := apiClient.getExchangeInfo()
	if err != nil {
		return nil, err
	}
	result := []*common.Symbol{}
	for _, symbol := range exchangeInfo.Symbols {
		result = append(result, &common.Symbol{
			RawBase:         symbol.BaseAsset,
			NormalizedBase:  strings.ToUpper(symbol.BaseAsset),
			RawQuote:        symbol.QuoteAsset,
			NormalizedQuote: strings.ToUpper(symbol.QuoteAsset),
		})
	}
	return common.FilterSupportedAssets(result), nil
}

func (apiClient *ApiClient) GetRawMarketData() ([]*common.RawMarketData, error) {
	return nil, fmt.Errorf("not implemented")
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helpers
//////////////////////////////////////////////////////////////////////////////////////////////////////////

// Get CandleStick data from [startTime, endTime]
func (apiClient *ApiClient) GetOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval common.Interval,
	startTime time.Time,
	endTime time.Time,
) ([]*common.OHLCMarketData, error) {
	candleStickResponse, err := apiClient.getCandleStickData(
		baseSymbol,
		quoteSymbol,
		interval,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}
	ohlcMarketData := []*common.OHLCMarketData{}

	for i := range candleStickResponse {
		ohlcMarketData = append(ohlcMarketData, &common.OHLCMarketData{
			MarketData: common.MarketData{
				Source:        BINANCE,
				BaseCurrency:  baseSymbol,
				QuoteCurrency: quoteSymbol,
			},
			StartTime:  time.Unix(int64(candleStickResponse[i].OpenTime/1000), 0),
			EndTime:    time.Unix(int64(candleStickResponse[i].CloseTime/1000), 0),
			OpenPrice:  candleStickResponse[i].OpenPrice,
			HighPrice:  candleStickResponse[i].HighPrice,
			LowPrice:   candleStickResponse[i].LowPrice,
			ClosePrice: candleStickResponse[i].ClosePrice,
			Volume:     candleStickResponse[i].Volume,
		})
	}
	return ohlcMarketData, nil
}
