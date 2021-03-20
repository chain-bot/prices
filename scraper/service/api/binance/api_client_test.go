package binance

import (
	"context"
	"fmt"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/mochahub/coinprice-scraper/scraper/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBinanceClient(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	exchangeClient := NewBinanceAPIClient()
	pass := true
	pass = t.Run("TestGetCandleStickData", func(t *testing.T) {
		expectedLength := 480 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength - time.Minute)
		candleStickData, err := exchangeClient.GetOHLCMarketData(
			models.Symbol{
				RawBase:         "BTC",
				NormalizedBase:  "BTC",
				RawQuote:        "USDT",
				NormalizedQuote: "USDT",
				ProductID:       "BTCUSDT",
			},
			time.Minute,
			startTime,
			endTime,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickData)
		assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData))
	}) && pass
	pass = t.Run("TestGetExchangeInfo", func(t *testing.T) {
		exchangeInfo, err := exchangeClient.getExchangeInfo()
		assert.NoError(t, err)
		assert.NotNil(t, exchangeInfo)
		//fmt.Print(utils.PrettyJSON(exchangeInfo))
	}) && pass

	// Interface Methods
	pass = t.Run("TestGetSupportedPairs", func(t *testing.T) {
		pairs, err := exchangeClient.GetSupportedPairs()
		assert.Nil(t, err)
		assert.NotEmpty(t, pairs)
		expectedPair := models.Symbol{
			RawBase:         "BTC",
			RawQuote:        "USDT",
			NormalizedBase:  "BTC",
			NormalizedQuote: "USDT",
			ProductID:       "BTCUSDT",
		}
		assert.GreaterOrEqual(t, len(pairs), 3)
		assert.Contains(t, pairs, &expectedPair)
	}) && pass
	// Should get all prices from [start, end)
	pass = t.Run("TestGetAllOHLCMarketData", func(t *testing.T) {
		expectedLength := 12000 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)

		candleStickData, err := exchangeClient.GetAllOHLCMarketData(
			models.Symbol{
				RawBase:         "BTC",
				NormalizedBase:  "BTC",
				RawQuote:        "USDT",
				NormalizedQuote: "USDT",
				ProductID:       "BTCUSDT",
			},
			time.Minute,
			startTime,
			endTime,
		)

		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickData)
		assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData))
		assert.Equal(t, startTime.String(), candleStickData[0].StartTime.UTC().String())
		assert.Equal(t, endTime.String(), candleStickData[len(candleStickData)-1].EndTime.UTC().String())
	}) && pass

	pass = t.Run("TestWebSocket", func(t *testing.T) {
		symbol := models.Symbol{
			RawBase:         "BTC",
			NormalizedBase:  "BTC",
			RawQuote:        "USDT",
			NormalizedQuote: "USDT",
			ProductID:       "BTCUSDT",
		}
		ohlcMarketDataChannel, err := exchangeClient.GetOHLCMarketDataChannel(ctx, symbol, time.Minute)
		assert.Nil(t, err)
		for i := 0; i < 5; i += 1 {
			ohlc := <-ohlcMarketDataChannel
			fmt.Println(utils.PrettyJSON(ohlc))
		}
		cancelCtx()
		<-ctx.Done()

	}) && pass
	assert.Equal(t, true, pass)
}
