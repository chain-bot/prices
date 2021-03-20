package binance

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"log"
	"strings"
	"time"
)

// https://github.com/rootpd/go-binance/blob/master/service_websocket.go
func (apiClient *ApiClient) GetOHLCMarketDataChannel(
	symbol models.Symbol,
	interval time.Duration,
) (chan *models.OHLCMarketData, error) {
	url := fmt.Sprintf(klineSocketStream, strings.ToLower(symbol.ProductID), string(getBinanceIntervalFromDuration(interval)))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	ohlcMarketDataChannel := make(chan *models.OHLCMarketData)

	go func() {
		defer c.Close()
		for {
			select {
			case <-apiClient.Context.Done():
				log.Println("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("wsRead", err)
					return
				}
				klineResponse := KlineResponse{}
				if err := json.Unmarshal(message, &klineResponse); err != nil {
					log.Println("wsUnmarshal", err, "body", string(message))
					return
				}
				candleStart := time.Unix(int64(klineResponse.Kline.KlineStart/1000), 0)
				// We don't use the candle end time from binance because they return 59 seconds opposed to 0 seconds of next minute
				candleEnd := candleStart.Add(interval)
				ohlcv := &models.OHLCMarketData{
					MarketData: models.MarketData{
						Source:        BINANCE,
						BaseCurrency:  symbol.NormalizedBase,
						QuoteCurrency: symbol.NormalizedQuote,
					},
					StartTime:  candleStart,
					EndTime:    candleEnd,
					OpenPrice:  klineResponse.Kline.Open,
					HighPrice:  klineResponse.Kline.High,
					LowPrice:   klineResponse.Kline.Low,
					ClosePrice: klineResponse.Kline.Close,
					Volume:     klineResponse.Kline.Volume,
				}
				ohlcMarketDataChannel <- ohlcv
			}
		}
	}()
	return ohlcMarketDataChannel, nil
}
