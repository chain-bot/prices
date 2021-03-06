package binance

import (
	"fmt"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/common"
	"strings"
	"time"
)

// Get CandleStick data from [startTime, endTime] with pagination
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
	exchangeInfo, err := apiClient.getExchangeInfo()
	if err != nil {
		return nil, err
	}
	result := []*models.Symbol{}
	for _, symbol := range exchangeInfo.Symbols {
		result = append(result, &models.Symbol{
			RawBase:         symbol.BaseAsset,
			NormalizedBase:  strings.ToUpper(symbol.BaseAsset),
			RawQuote:        symbol.QuoteAsset,
			NormalizedQuote: strings.ToUpper(symbol.QuoteAsset),
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

// Get CandleStick data from [startTime, endTime]
func (apiClient *ApiClient) GetOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCMarketData, error) {
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
	ohlcMarketData := []*models.OHLCMarketData{}

	for i := range candleStickResponse {
		candleStart := time.Unix(int64(candleStickResponse[i].OpenTime/1000), 0)
		// We don't use the candle end time from binance because they return 59 seconds opposed to 0 seconds of next minute
		candleEnd := candleStart.Add(interval)
		ohlcMarketData = append(ohlcMarketData, &models.OHLCMarketData{
			MarketData: models.MarketData{
				Source:        BINANCE,
				BaseCurrency:  baseSymbol,
				QuoteCurrency: quoteSymbol,
			},
			StartTime:  candleStart,
			EndTime:    candleEnd,
			OpenPrice:  candleStickResponse[i].OpenPrice,
			HighPrice:  candleStickResponse[i].HighPrice,
			LowPrice:   candleStickResponse[i].LowPrice,
			ClosePrice: candleStickResponse[i].ClosePrice,
			Volume:     candleStickResponse[i].Volume,
		})
	}
	return ohlcMarketData, nil
}
