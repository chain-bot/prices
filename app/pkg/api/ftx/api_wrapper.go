package ftx

import (
	"fmt"
	"github.com/chain-bot/scraper/app/pkg/api/common"
	"github.com/chain-bot/scraper/app/pkg/models"
	"strings"
	"time"
)

//Get CandleStick data from [startTime, endTime) with pagination
func (apiClient *ApiClient) GetAllOHLCMarketData(
	symbol models.Symbol,
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
			symbol,
			interval,
			startTime,
			newEndTime.Add(-interval))
		if err != nil {
			return nil, err
		}
		result = append(result, ohlcMarketData...)
		startTime = newEndTime
	}
	return result, nil
}

func (apiClient *ApiClient) GetSupportedPairs() ([]*models.Symbol, error) {
	marketResponse, err := apiClient.getMarkets()
	if err != nil {
		return nil, err
	}
	result := []*models.Symbol{}
	for _, symbol := range marketResponse.Result {
		// Only Spot Markets Supported Rn
		if symbol.Type != "spot" {
			continue
		}
		quote := symbol.QuoteCurrency
		normalizedQuote := strings.ToUpper(quote)
		base := symbol.BaseCurrency
		normalizedBase := strings.ToUpper(base)
		result = append(result, &models.Symbol{
			RawBase:         base,
			NormalizedBase:  normalizedBase,
			RawQuote:        quote,
			NormalizedQuote: normalizedQuote,
			ProductID:       symbol.Name,
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

//Get CandleStick data from [startTime, endTime)
func (apiClient *ApiClient) GetOHLCMarketData(
	symbol models.Symbol,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCMarketData, error) {
	historicalPriceResponse, err := apiClient.getHistoricalPrices(
		symbol.ProductID,
		interval,
		startTime,
		endTime,
		int64(endTime.Sub(startTime))/int64(interval),
	)
	if err != nil {
		return nil, err
	}
	ohlcMarketData := []*models.OHLCMarketData{}
	for i := range historicalPriceResponse.Result {
		candle := historicalPriceResponse.Result[i]
		ohlcMarketData = append(ohlcMarketData, &models.OHLCMarketData{
			MarketData: models.MarketData{
				Source:        FTX,
				BaseCurrency:  symbol.NormalizedBase,
				QuoteCurrency: symbol.NormalizedQuote,
			},
			StartTime:  candle.StartTime,
			EndTime:    candle.StartTime.Add(interval),
			OpenPrice:  candle.Open,
			HighPrice:  candle.High,
			LowPrice:   candle.Low,
			ClosePrice: candle.Close,
			Volume:     candle.Volume,
		})
	}
	return ohlcMarketData, nil
}
