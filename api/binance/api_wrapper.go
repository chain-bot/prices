package binance

import (
	"errors"
	"fmt"
	"time"

	"github.com/mochahub/coinprice-scraper/api"
)

// Get CandleStick data from [startTime, endTime]
func (apiClient *apiClient) GetOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval api.Interval,
	startTime time.Time,
	endTime time.Time,
) ([]*api.OHLCMarketData, error) {
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
	ohlcMarketData := []*api.OHLCMarketData{}
	for i := range candleStickResponse {
		ohlcMarketData = append(ohlcMarketData, &api.OHLCMarketData{
			MarketData: api.MarketData{
				Source:        InstitutionIdentifier,
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

func (apiClient *apiClient) GetAllOHLCMarketData(
	baseSymbol string,
	quoteSymbol string,
	interval api.Interval,
	startTime time.Time,
	endTime time.Time,
) ([]*api.OHLCMarketData, error) {
	var durationFromInterval time.Duration
	switch interval {
	case api.BinanceDayInterval:
		durationFromInterval = time.Hour * 24
	case api.BinanceHourInterval:
		durationFromInterval = time.Hour
	case api.BinanceMinuteInterval:
		durationFromInterval = time.Minute
	default:
		return nil, errors.New(fmt.Sprintf("unknown interval: %s", interval))
	}
	result := []*api.OHLCMarketData{}
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

func (apiClient *apiClient) GetRawMarketData() ([]*api.RawMarketData, error) {
	return nil, fmt.Errorf("not implemented")
}
