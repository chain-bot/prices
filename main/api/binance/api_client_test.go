package binance

import (
	"fmt"
	"github.com/mochahub/coinprice-scraper/main/utils"
	"os"
	"testing"
	"time"

	"github.com/mochahub/coinprice-scraper/main/api"
)

func TestApiClient_GetCandleStickData(t *testing.T) {
	// TODO(Zahin): Use a config file for api key or a env variable
	binanceAPIClient := NewBinanceAPIClient(
		os.Getenv("BINANCE_API_KEY"),
		1)
	expectedLength := 480 * time.Minute
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(expectedLength - time.Minute)
	candleStickData, err := binanceAPIClient.GetOHLCMarketData(
		"BTC",
		"USDT",
		api.Minute,
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
		os.Getenv("BINANCE_API_KEY"),
		1)
	expectedLength := 1000 * time.Minute
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(expectedLength - time.Minute)
	candleStickData, err := binanceAPIClient.GetAllOHLCMarketData(
		"BTC",
		"USDT",
		api.Minute,
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

func TestApiClient_GetExchangeInfo(t *testing.T) {
	binanceAPIClient := NewBinanceAPIClient(
		os.Getenv("BINANCE_API_KEY"),
		1)
	exchangeInfo, err := binanceAPIClient.getExchangeInfo()
	if err != nil {
		t.Error(err)
	}
	if exchangeInfo == nil {
		t.Error("empty exchange info")
	}
	fmt.Print(utils.PrettyJSON(exchangeInfo))
}

func TestApiClient_GetSupportedPairs(t *testing.T) {
	binanceAPIClient := NewBinanceAPIClient(
		os.Getenv("BINANCE_API_KEY"),
		1)
	symbols, err := binanceAPIClient.GetSupportedPairs()
	if err != nil {
		t.Error(err)
	}
	if len(symbols) == 0 {
		t.Error("empty exchange info")
	}
	fmt.Print(utils.PrettyJSON(symbols))
}
