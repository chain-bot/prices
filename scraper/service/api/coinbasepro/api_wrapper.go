package coinbasepro

import (
	"fmt"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/common"
	"github.com/mochahub/coinprice-scraper/scraper/utils"
	"strings"
	"time"
)

// Get CandleStick data from [startTime, endTime) with pagination
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
	products, err := apiClient.getProducts()
	if err != nil {
		return nil, err
	}
	result := []*models.Symbol{}

	for i := range products {
		product := products[i]
		quote := product.QuoteCurrency
		normalizedQuote := GetCoinpriceSymbolFromCoinbasePro(quote)
		base := product.BaseCurrency
		normalizedBase := GetCoinpriceSymbolFromCoinbasePro(base)
		newPair := &models.Symbol{
			RawBase:         base,
			NormalizedBase:  strings.ToUpper(normalizedBase),
			RawQuote:        quote,
			NormalizedQuote: strings.ToUpper(normalizedQuote),
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
func (apiClient *ApiClient) GetOHLCMarketData(
	symbol models.Symbol,
	durationInterval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCMarketData, error) {
	candleStickData, err := apiClient.getCandleStickData(
		int(durationInterval.Seconds()), startTime, endTime, symbol.ProductID)
	if err != nil {
		return nil, err
	}
	result := []*models.OHLCMarketData{}
	for i := range candleStickData {
		candle := candleStickData[i]
		candleEnd := time.Unix(int64(candle.CloseTime), 0)
		result = append(result, &models.OHLCMarketData{
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
