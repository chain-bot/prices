package binance

import (
	"github.com/chain-bot/prices/app/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBinanceClient(t *testing.T) {
	exchangeClient := NewBinanceAPIClient()
	pass := true

	pass = t.Run("Test GetExchangeIdentifier Pass", func(t *testing.T) {
		assert.NotEqual(t, "", exchangeClient.GetExchangeIdentifier())
	}) && pass
	pass = t.Run("Test GetCandleStickData Pass", func(t *testing.T) {
		expectedLength := 480 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength - time.Minute)
		candleStickData, err := exchangeClient.GetOHLCVMarketData(
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
	pass = t.Run("Test getCandleStickData Fail", func(t *testing.T) {
		expectedLength := 60 * time.Minute
		startTime := time.Now().Add(-expectedLength)
		_, err := exchangeClient.getCandleStickData(
			"poopee",
			time.Minute,
			startTime,
			time.Time{},
		)
		assert.Error(t, err)
	}) && pass
	pass = t.Run("Test GetExchangeInfo Pass", func(t *testing.T) {
		exchangeInfo, err := exchangeClient.getExchangeInfo()
		assert.NoError(t, err)
		assert.NotNil(t, exchangeInfo)
		//fmt.Print(utils.PrettyJSON(exchangeInfo))
	}) && pass
	// Interface Methods
	pass = t.Run("Test GetSupportedPairs Pass", func(t *testing.T) {
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
	pass = t.Run("Test GetAllOHLCVMarketData Pass", func(t *testing.T) {
		expectedLength := 12000 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)
		candleStickData, err := exchangeClient.GetAllOHLCVMarketData(
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
	assert.Equal(t, true, pass)
}
