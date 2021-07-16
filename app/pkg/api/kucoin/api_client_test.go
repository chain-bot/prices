package kucoin

import (
	"github.com/chain-bot/prices/app/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestKucoinClient(t *testing.T) {
	exchangeClient := NewKucoinAPIClient()
	pass := true

	pass = t.Run("TestGetExchangeIdentifier", func(t *testing.T) {
		assert.NotEqual(t, "", exchangeClient.GetExchangeIdentifier())
	}) && pass
	// Get Candles from [start, end]
	pass = t.Run("Test getKlines pass", func(t *testing.T) {
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
	pass = t.Run("Test getKlines fail", func(t *testing.T) {
		expectedLength := 480 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)
		candleStickResponse, err := exchangeClient.getKlines(
			"not-a-real-symbol",
			time.Minute,
			startTime,
			endTime,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickResponse.Msg)
		assert.NotEqual(t, 0, candleStickResponse.Code)
	}) && pass
	pass = t.Run("Test GetSymbols", func(t *testing.T) {
		exchangeInfo, err := exchangeClient.getSymbols()
		assert.NoError(t, err)
		assert.NotNil(t, exchangeInfo)
		//fmt.Print(utils.PrettyJSON(exchangeInfo))
	}) && pass
	// Interface Methods
	pass = t.Run("Test GetSupportedPairs", func(t *testing.T) {
		pairs, err := exchangeClient.GetSupportedPairs()
		assert.Nil(t, err)
		assert.NotEmpty(t, pairs)
		expectedPair := models.Symbol{
			RawBase:         "ETH",
			RawQuote:        "BTC",
			NormalizedBase:  "ETH",
			NormalizedQuote: "BTC",
			ProductID:       "ETH-BTC",
		}
		assert.GreaterOrEqual(t, len(pairs), 3)
		assert.Contains(t, pairs, &expectedPair)
	}) && pass

	// Should get all prices from [start, end)
	pass = t.Run("Test GetAllOHLCVMarketData", func(t *testing.T) {
		expectedLength := 12000 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)

		candleStickData, err := exchangeClient.GetAllOHLCVMarketData(
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
