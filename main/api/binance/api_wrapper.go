package binance

import (
	"fmt"
	"github.com/mochahub/coinprice-scraper/main/api/common"
	"github.com/mochahub/coinprice-scraper/main/config"
	"strings"
	"time"
)

// Get CandleStick data from [startTime, endTime] with pagination
func (apiClient *apiClient) GetAllOHLCMarketData(
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
	for startTime.Before(endTime) || startTime.Equal(endTime) {
		newEndTime := startTime.Add(maxLimit * durationFromInterval)
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
		startTime = newEndTime.Add(durationFromInterval)
	}
	return result, nil
}

func (apiClient *apiClient) GetSupportedPairs() ([]*common.Symbol, error) {
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
	return filterSupportedAssets(result), nil
}

func (apiClient *apiClient) GetRawMarketData() ([]*common.RawMarketData, error) {
	return nil, fmt.Errorf("not implemented")
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helpers
//////////////////////////////////////////////////////////////////////////////////////////////////////////

// Get CandleStick data from [startTime, endTime]
func (apiClient *apiClient) GetOHLCMarketData(
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
			StartTime: time.Unix(int64(candleStickResponse[i].OpenTime/1000), 0),
			EndTime:   time.Unix(int64(candleStickResponse[i].CloseTime/1000), 0),
			HighPrice: candleStickResponse[i].HighPrice,
			LowPrice:  candleStickResponse[i].LowPrice,
			Volume:    candleStickResponse[i].Volume,
		})
	}
	return ohlcMarketData, nil
}

func filterSupportedAssets(symbols []*common.Symbol) []*common.Symbol {
	result := []*common.Symbol{}
	supportedAssets := config.GetSupportedAssets()
	for index := range symbols {
		pair := symbols[index]
		_, ok := supportedAssets[pair.NormalizedBase]
		if !ok {
			continue
		}
		_, ok = supportedAssets[pair.NormalizedQuote]
		if !ok {
			continue
		}
		result = append(result, pair)
	}
	return result
}
