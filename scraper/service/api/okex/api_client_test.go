package okex

import (
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOkexClient(t *testing.T) {
	// TODO: Use DI instead of calling GetSecrets directly
	exchangeClient := NewOkexAPIClient()
	pass := true
	// Get Candles from [start, end]
	pass = t.Run("TestGetCandleStickData", func(t *testing.T) {
		expectedLength := 50 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)
		candleStickResponse, err := exchangeClient.getInstrumentCandles(
			"BTC-USDT",
			time.Minute,
			startTime,
			endTime,
			int64(expectedLength.Minutes()),
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickResponse)
	}) && pass
	pass = t.Run("TestGetSymbols", func(t *testing.T) {
		exchangeInfo, err := exchangeClient.getInstruments()
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
			RawBase:         "ETH",
			RawQuote:        "USDT",
			NormalizedBase:  "ETH",
			NormalizedQuote: "USDT",
			ProductID:       "ETH-USDT",
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
