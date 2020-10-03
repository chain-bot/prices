package binance

import (
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

func (apiClient *apiClient) GetRawMarketData() ([]*api.RawMarketData, error) {
	return nil, fmt.Errorf("not implemented")
}
