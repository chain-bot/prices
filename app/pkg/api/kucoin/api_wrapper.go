package kucoin

import (
	"fmt"
	"github.com/chain-bot/scraper/app/pkg/api/common"
	"github.com/chain-bot/scraper/app/pkg/models"
	"github.com/chain-bot/scraper/app/utils"
	"strings"
	"time"
)

//Get CandleStick data from [startTime, endTime) with pagination
func (apiClient *ApiClient) GetAllOHLCVMarketData(
	symbol models.Symbol,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCVMarketData, error) {
	if endTime.IsZero() {
		endTime = time.Now()
	}
	result := []*models.OHLCVMarketData{}
	for startTime.Before(endTime) {
		newEndTime := startTime.Add(maxLimit * interval)
		if newEndTime.After(endTime) {
			newEndTime = endTime
		}
		ohlcvMarketData, err := apiClient.GetohlcvMarketData(
			symbol,
			interval,
			startTime,
			newEndTime)
		if err != nil {
			return nil, err
		}
		result = append(result, ohlcvMarketData...)
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
			ProductID:       symbol.Symbol,
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
func (apiClient *ApiClient) GetohlcvMarketData(
	symbol models.Symbol,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCVMarketData, error) {
	candleStickResponse, err := apiClient.getKlines(
		symbol.ProductID,
		interval,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}
	ohlcvMarketData := []*models.OHLCVMarketData{}
	for i := range candleStickResponse.Data {
		candle := candleStickResponse.Data[i]
		ohlcvMarketData = append(ohlcvMarketData, &models.OHLCVMarketData{
			MarketData: models.MarketData{
				Source:        KUCOIN,
				BaseCurrency:  symbol.NormalizedBase,
				QuoteCurrency: symbol.NormalizedQuote,
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
	return utils.Reverse(ohlcvMarketData), nil
}
