package binance

import (
	"testing"
	"time"

	"github.com/mochahub/coinprice-scraper/api"
)

func TestApiClient_GetCandleStickData(t *testing.T) {
	// TODO(Zahin): Use a config file for api key or a env variable
	binanceAPIClient := NewBinanceAPIClient(
		"ADD_API_KEY",
		1,
		1200)
	expectedLength := 480 * time.Minute
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(expectedLength - time.Minute)
	candleStickData, err := binanceAPIClient.GetOHLCMarketData(
		"BTC",
		"USDT",
		api.BinanceMinuteInterval,
		startTime,
		endTime,
	)
	if err != nil {
		t.Error(err)
	}
	if candleStickData == nil {
		t.Error("empty Prices")
	}
	if len(candleStickData) != int(expectedLength.Minutes()) {
		t.Errorf("expected %d got %d", int(expectedLength.Minutes()), len(candleStickData))
	}
}

func TestApiClient_GetAllOHLCMarketData(t *testing.T) {
	binanceAPIClient := NewBinanceAPIClient(
		"ADD_API_KEY",
		1,
		1200)
	expectedLength := 1001 * time.Minute
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(expectedLength - time.Minute)
	candleStickData, err := binanceAPIClient.GetAllOHLCMarketData(
		"BTC",
		"USDT",
		api.BinanceMinuteInterval,
		startTime,
		endTime,
	)
	if err != nil {
		t.Error(err)
	}
	if candleStickData == nil {
		t.Error("empty Prices")
	}
	if len(candleStickData) != int(expectedLength.Minutes()) {
		t.Errorf("expected %d got %d", int(expectedLength.Minutes()), len(candleStickData))
	}
}
