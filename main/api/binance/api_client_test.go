package binance

import (
	"github.com/mochahub/coinprice-scraper/main/api/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestBinanceClient(t *testing.T) {
	exchangeClient := NewBinanceAPIClient(os.Getenv("BINANCE_API_KEY"))
	pass := true
	pass = t.Run("TestGetCandleStickData", func(t *testing.T) {
		expectedLength := 480 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength - time.Minute)
		candleStickData, err := exchangeClient.GetOHLCMarketData(
			"BTC",
			"USDT",
			common.Minute,
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
	// TODO(Zahin): Do we even need this? exhange_clients_test will test it as well...
	pass = t.Run("TestGetAllOHLCMarketData", func(t *testing.T) {
		expectedLength := 1000 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength - time.Minute)
		candleStickData, err := exchangeClient.GetAllOHLCMarketData(
			"BTC",
			"USDT",
			common.Minute,
			startTime,
			endTime,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickData)
		assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData))
	}) && pass

	pass = t.Run("TestGetSupportedPairs", func(t *testing.T) {
		symbols, err := exchangeClient.GetSupportedPairs()
		assert.Nil(t, err)
		assert.NotEmpty(t, symbols)
		//fmt.Print(utils.PrettyJSON(symbols))
	}) && pass
	assert.Equal(t, true, pass)
}
