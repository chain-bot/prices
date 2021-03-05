package kucoin

import (
	"fmt"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/common"
	"strings"
	"time"
)

//Get CandleStick data from [startTime, endTime] with pagination
func (apiClient *ApiClient) GetAllOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval common.Interval,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCMarketData, error) {
	// TODO: We should just change the interface signature to use duration instead of a custom type
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
	// Kucoin returns data as [start, end)
	endTime = endTime.Add(durationFromInterval)

	result := []*models.OHLCMarketData{}
	for startTime.Before(endTime) {
		newEndTime := startTime.Add(maxLimit * durationFromInterval)
		if newEndTime.After(endTime) {
			newEndTime = endTime
		}
		ohlcMarketData, err := apiClient.GetOHLCMarketData(
			baseSymbol,
			quoteSymbol,
			durationFromInterval,
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

func (apiClient *ApiClient) GetSupportedPairs() ([]*models.Symbol, error) {
	exchangeSymbols, err := apiClient.getSymbols()
	if err != nil {
		return nil, err
	}
	result := []*models.Symbol{}
	for _, symbol := range exchangeSymbols.Data {
		result = append(result, &models.Symbol{
			RawBase:         symbol.BaseCurrency,
			NormalizedBase:  strings.ToUpper(symbol.BaseCurrency),
			RawQuote:        symbol.QuoteCurrency,
			NormalizedQuote: strings.ToUpper(symbol.QuoteCurrency),
		})
	}
	return common.FilterSupportedAssets(result), nil
}

func (apiClient *ApiClient) GetRawMarketData() ([]*models.RawMarketData, error) {
	return nil, fmt.Errorf("not implemented")
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helpers
//////////////////////////////////////////////////////////////////////////////////////////////////////////

//Get CandleStick data from [startTime, endTime]
func (apiClient *ApiClient) GetOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCMarketData, error) {
	candleStickResponse, err := apiClient.getKlines(
		fmt.Sprintf("%s-%s", baseSymbol, quoteSymbol),
		interval,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}
	ohlcMarketData := []*models.OHLCMarketData{}
	for i := range candleStickResponse.Data {
		ohlcMarketData = append(ohlcMarketData, &models.OHLCMarketData{
			MarketData: models.MarketData{
				Source:        KUCOIN,
				BaseCurrency:  baseSymbol,
				QuoteCurrency: quoteSymbol,
			},
			StartTime:  time.Unix(int64(candleStickResponse.Data[i].OpenTime), 0),
			EndTime:    time.Unix(int64(candleStickResponse.Data[i].OpenTime+interval.Seconds()), 0),
			OpenPrice:  candleStickResponse.Data[i].OpenPrice,
			HighPrice:  candleStickResponse.Data[i].HighPrice,
			LowPrice:   candleStickResponse.Data[i].LowPrice,
			ClosePrice: candleStickResponse.Data[i].ClosePrice,
			Volume:     candleStickResponse.Data[i].Volume,
		})
	}
	return ohlcMarketData, nil
}
