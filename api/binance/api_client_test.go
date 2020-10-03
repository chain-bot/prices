package binance

import (
	"testing"
	"time"
)

func TestApiClient_GetCandleStickData(t *testing.T) {
	// TODO(Zahin): Use a config file for api key or a env variable
	binanceAPIClient := NewBinanceAPIClient(
		"ADD_API_KEY",
		1,
		1200)
	expectedLength := 480 * time.Minute
	startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(expectedLength)
	candleStickData, err := binanceAPIClient.getCandleStickData(
		"BTC",
		"USDT",
		MinuteInterval,
		startTime,
		endTime,
		0,
	)
	if err != nil {
		t.Error(err)
	}
	if candleStickData == nil {
		t.Error("empty Prices")
	}
	if float64(len(candleStickData)) != expectedLength.Minutes()+1 {
		t.Errorf("expected %f got %d", expectedLength.Minutes(), len(candleStickData))
	}
}
