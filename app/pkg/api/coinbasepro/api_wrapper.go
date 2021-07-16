package coinbasepro

import (
	"fmt"
	"github.com/chain-bot/prices/app/pkg/api/common"
	"github.com/chain-bot/prices/app/pkg/models"
	"github.com/chain-bot/prices/app/utils"
	"strings"
	"time"
)

// Get CandleStick data from [startTime, endTime) with pagination
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
		ohlcvMarketData, err := apiClient.GetOHLCVMarketData(
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
	products, err := apiClient.getProducts()
	if err != nil {
		return nil, err
	}
	result := []*models.Symbol{}

	for i := range products {
		product := products[i]
		quote := product.QuoteCurrency
		normalizedQuote := strings.ToUpper(quote)
		base := product.BaseCurrency
		normalizedBase := strings.ToUpper(base)
		if _, ok := coinbaseProInstrumentFilter[normalizedBase]; ok {
			continue
		}
		if _, ok := coinbaseProInstrumentFilter[normalizedQuote]; ok {
			continue
		}
		newPair := &models.Symbol{
			RawBase:         base,
			NormalizedBase:  normalizedBase,
			RawQuote:        quote,
			NormalizedQuote: normalizedQuote,
			ProductID:       product.ID,
		}

		result = append(result, newPair)
	}
	return common.FilterSupportedAssets(result), nil
}

func (apiClient *ApiClient) GetRawMarketData() ([]*models.RawMarketData, error) {
	return nil, fmt.Errorf("not implemented")
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helpers
//////////////////////////////////////////////////////////////////////////////////////////////////////////

// Get CandleStick data from [startTime, endTime]
func (apiClient *ApiClient) GetOHLCVMarketData(
	symbol models.Symbol,
	durationInterval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCVMarketData, error) {
	candleStickData, err := apiClient.getCandleStickData(
		int(durationInterval.Seconds()), startTime, endTime, symbol.ProductID)
	if err != nil {
		return nil, err
	}
	result := []*models.OHLCVMarketData{}
	for i := range candleStickData {
		candle := candleStickData[i]
		candleEnd := time.Unix(int64(candle.CloseTime), 0)
		result = append(result, &models.OHLCVMarketData{
			MarketData: models.MarketData{
				Source:        apiClient.GetExchangeIdentifier(),
				BaseCurrency:  symbol.NormalizedBase,
				QuoteCurrency: symbol.NormalizedQuote,
			},
			StartTime:  candleEnd.Add(-durationInterval),
			EndTime:    candleEnd,
			OpenPrice:  candle.OpenPrice,
			HighPrice:  candle.ClosePrice,
			LowPrice:   candle.LowPrice,
			ClosePrice: candle.ClosePrice,
			Volume:     candle.Volume,
		})
	}
	return utils.Reverse(result), nil
}
