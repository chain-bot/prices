package ftx

import (
	"github.com/chain-bot/prices/app/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFtxClient(t *testing.T) {
	exchangeClient := NewFtxAPIClient()
	pass := true

	pass = t.Run("TestGetExchangeIdentifier", func(t *testing.T) {
		assert.NotEqual(t, "", exchangeClient.GetExchangeIdentifier())
	}) && pass
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
		expectedPair := models.Symbol{
			RawBase:         "ETH",
			RawQuote:        "BTC",
			NormalizedBase:  "ETH",
			NormalizedQuote: "BTC",
			ProductID:       "ETH/BTC",
		}
		assert.GreaterOrEqual(t, len(pairs), 3)
		assert.Contains(t, pairs, &expectedPair)
	}) && pass
	// Should get all prices from [start, end)
	pass = t.Run("TestGetAllOHLCVMarketData", func(t *testing.T) {
		expectedLength := 12000 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)
		candleStickData, err := exchangeClient.GetAllOHLCVMarketData(
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
