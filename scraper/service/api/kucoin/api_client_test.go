package kucoin

import (
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestKucoinClient(t *testing.T) {
	// TODO: Use DI instead of calling GetSecrets directly
	config.LoadEnv()
	secret, _ := config.GetSecrets()
	exchangeClient := NewKucoinAPIClient(secret)
	pass := true
	// Get Candles from [start, end]
	pass = t.Run("TestGetCandleStickData", func(t *testing.T) {
		expectedLength := 480 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)
		candleStickResponse, err := exchangeClient.getKlines(
			"BTC-USDT",
			time.Minute,
			startTime,
			endTime,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickResponse.Data)
	}) && pass
	pass = t.Run("TestGetSymbols", func(t *testing.T) {
		exchangeInfo, err := exchangeClient.getSymbols()
		assert.NoError(t, err)
		assert.NotNil(t, exchangeInfo)
		//fmt.Print(utils.PrettyJSON(exchangeInfo))
	}) && pass
	// Interface Methods
	pass = t.Run("TestGetSupportedPairs", func(t *testing.T) {
		pairs, err := exchangeClient.GetSupportedPairs()
		assert.Nil(t, err)
		assert.NotEmpty(t, pairs)
		// TODO: This test is not scale-able unless we use DI to inject what pairs are supported
		expectedPairs := map[string]models.Symbol{
			"BTC-USDT": {
				RawBase:         "BTC",
				RawQuote:        "USDT",
				NormalizedBase:  "BTC",
				NormalizedQuote: "USDT",
				ProductID:       "BTC-USDT",
			},
			"ETH-USDT": {
				RawBase:         "ETH",
				RawQuote:        "USDT",
				NormalizedBase:  "ETH",
				NormalizedQuote: "USDT",
				ProductID:       "ETH-USDT",
			},
			"ETH-BTC": {
				RawBase:         "ETH",
				RawQuote:        "BTC",
				NormalizedBase:  "ETH",
				NormalizedQuote: "BTC",
				ProductID:       "ETH-BTC",
			},
		}
		assert.Equal(t, 3, len(pairs))
		for _, pair := range pairs {
			expectedSymbol, ok := expectedPairs[pair.ProductID]
			assert.Equal(t, true, ok)
			assert.Equal(t, expectedSymbol, *pair)
		}
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
				ProductID:       "BTC-USDT",
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

	assert.Equal(t, true, pass)
}
