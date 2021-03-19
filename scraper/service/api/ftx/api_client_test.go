package ftx

import (
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFtxClient(t *testing.T) {
	// TODO: Use DI instead of calling GetSecrets directly
	exchangeClient := NewFtxAPIClient()
	pass := true
	// Get Candles from [start, end)
	pass = t.Run("TestGetCandleStickData", func(t *testing.T) {
		expectedLength := 50 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)
		candleStickResponse, err := exchangeClient.getHistoricalPrices(
			"BTC/USD",
			time.Minute,
			startTime,
			endTime,
			int64(expectedLength.Minutes()),
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickResponse.Result)
	}) && pass
	pass = t.Run("TestGetMarkets", func(t *testing.T) {
		markets, err := exchangeClient.getMarkets()
		assert.NoError(t, err)
		assert.NotNil(t, markets)
		//fmt.Print(utils.PrettyJSON(markets))
	}) && pass

	// Interface Methods
	pass = t.Run("TestGetSupportedPairs", func(t *testing.T) {
		pairs, err := exchangeClient.GetSupportedPairs()
		assert.Nil(t, err)
		assert.NotEmpty(t, pairs)
		// TODO: This test is not scale-able unless we use DI to inject what pairs are supported
		expectedPairs := map[string]models.Symbol{
			"BTC/USD": {
				RawBase:         "BTC",
				RawQuote:        "USD",
				NormalizedBase:  "BTC",
				NormalizedQuote: "USD",
				ProductID:       "BTC/USD",
			},
			"BTC/USDT": {
				RawBase:         "BTC",
				RawQuote:        "USDT",
				NormalizedBase:  "BTC",
				NormalizedQuote: "USDT",
				ProductID:       "BTC/USDT",
			},
			"ETH/USD": {
				RawBase:         "ETH",
				RawQuote:        "USD",
				NormalizedBase:  "ETH",
				NormalizedQuote: "USD",
				ProductID:       "ETH/USD",
			},
			"ETH/USDT": {
				RawBase:         "ETH",
				RawQuote:        "USDT",
				NormalizedBase:  "ETH",
				NormalizedQuote: "USDT",
				ProductID:       "ETH/USDT",
			},
			"ETH/BTC": {
				RawBase:         "ETH",
				RawQuote:        "BTC",
				NormalizedBase:  "ETH",
				NormalizedQuote: "BTC",
				ProductID:       "ETH/BTC",
			},
			"USDT/USD": {
				RawBase:         "USDT",
				RawQuote:        "USD",
				NormalizedBase:  "USDT",
				NormalizedQuote: "USD",
				ProductID:       "USDT/USD",
			},
		}
		assert.Equal(t, 6, len(pairs))
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
				ProductID: "BTC/USD",
			},
			time.Minute,
			startTime,
			endTime,
		)

		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickData)
		assert.Equal(t, startTime.String(), candleStickData[0].StartTime.UTC().String())
		assert.Equal(t, endTime.String(), candleStickData[len(candleStickData)-1].EndTime.UTC().String())
		assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData))
	}) && pass

	assert.Equal(t, true, pass)
}
