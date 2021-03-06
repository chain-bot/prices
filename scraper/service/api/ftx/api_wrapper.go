package ftx

import (
	"fmt"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/common"
	"strings"
	"time"
)

//Get CandleStick data from [startTime, endTime) with pagination
func (apiClient *ApiClient) GetAllOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval common.Interval,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCMarketData, error) {
	// TODO: We should just change the interface signature to use duration instead of a custom type
	ftxBaseSymbol := GetFtxSymbolFromCoinprice(baseSymbol)
	ftxQuoteSymbol := GetFtxSymbolFromCoinprice(quoteSymbol)
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
	result := []*models.OHLCMarketData{}
	for startTime.Before(endTime) {
		newEndTime := startTime.Add(maxLimit * durationFromInterval)
		if newEndTime.After(endTime) {
			newEndTime = endTime
		}
		ohlcMarketData, err := apiClient.GetOHLCMarketData(
			ftxBaseSymbol,
			ftxQuoteSymbol,
			durationFromInterval,
			startTime,
			newEndTime.Add(-durationFromInterval))
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
		// Hack, assumed that USDT is the same as USD
		// TODO: We should revisit this assumption (same assumption with coinbasepro)
		if symbol.BaseCurrency == "USD" || symbol.QuoteCurrency == "USD" {
			continue
		}
		quote := symbol.QuoteCurrency
		normalizedQuote := GetCoinpriceSymbolFtx(quote)
		base := symbol.BaseCurrency
		normalizedBase := GetCoinpriceSymbolFtx(base)
		result = append(result, &models.Symbol{
			RawBase:         base,
			NormalizedBase:  strings.ToUpper(normalizedBase),
			RawQuote:        quote,
			NormalizedQuote: strings.ToUpper(normalizedQuote),
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
	baseSymbol string,
	quoteSymbol string,
	interval time.Duration,
	startTime time.Time,
	endTime time.Time,
) ([]*models.OHLCMarketData, error) {
	historicalPriceResponse, err := apiClient.getHistoricalPrices(
		fmt.Sprintf("%s/%s", baseSymbol, quoteSymbol),
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
				BaseCurrency:  baseSymbol,
				QuoteCurrency: quoteSymbol,
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
