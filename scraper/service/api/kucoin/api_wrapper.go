package kucoin

import (
	"fmt"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/common"
	"github.com/mochahub/coinprice-scraper/scraper/utils"
	"strings"
	"time"
)

//Get CandleStick data from [startTime, endTime) with pagination
func (apiClient *ApiClient) GetAllOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCMarketData, error) {
	if endTime.IsZero() {
		endTime = time.Now()
	}
	result := []*models.OHLCMarketData{}
	for startTime.Before(endTime) {
		newEndTime := startTime.Add(maxLimit * interval)
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
		candle := candleStickResponse.Data[i]
		ohlcMarketData = append(ohlcMarketData, &models.OHLCMarketData{
			MarketData: models.MarketData{
				Source:        KUCOIN,
				BaseCurrency:  baseSymbol,
				QuoteCurrency: quoteSymbol,
			},
			StartTime:  time.Unix(int64(candle.OpenTime), 0),
			EndTime:    time.Unix(int64(candle.OpenTime+interval.Seconds()), 0),
			OpenPrice:  candle.OpenPrice,
			HighPrice:  candle.HighPrice,
			LowPrice:   candle.LowPrice,
			ClosePrice: candle.ClosePrice,
			Volume:     candle.Volume,
		})
	}
	// Return ascending time
	return utils.Reverse(ohlcMarketData), nil
}
