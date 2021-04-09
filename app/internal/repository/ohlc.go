package repository

import (
	"github.com/chain-bot/scraper/app/pkg/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func (repo *RepositoryImpl) UpsertOHLCData(
	ohlcData []*models.OHLCMarketData,
	exchange string,
	pair *models.Symbol,
) {
	writeAPI := (*repo.influxClient).WriteAPI(repo.influxOrg, repo.ohlcBucket)
	tags := map[string]string{
		"quote":    pair.NormalizedQuote,
		"exchange": exchange,
	}
	for index := range ohlcData {
		ohlc := ohlcData[index]
		fields := map[string]interface{}{
			"open":   ohlc.OpenPrice,
			"high":   ohlc.HighPrice,
			"low":    ohlc.LowPrice,
			"close":  ohlc.ClosePrice,
			"volume": ohlc.Volume,
		}
		p := influxdb2.NewPoint(
			pair.NormalizedBase,
			tags,
			fields,
			ohlc.StartTime)
		writeAPI.WritePoint(p)
	}
}
