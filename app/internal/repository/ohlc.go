package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/chain-bot/prices/app/pkg/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/volatiletech/null/v8"
)

func (repo *RepositoryImpl) UpsertOHLCVData(
	ohlcvMarketData []*models.OHLCVMarketData,
	exchange string,
	pair *models.Symbol,
) {
	tags := map[string]string{
		"quote":    pair.NormalizedQuote,
		"exchange": exchange,
	}
	for index := range ohlcvMarketData {
		ohlcv := ohlcvMarketData[index]
		fields := map[string]interface{}{
			"open":   ohlcv.OpenPrice,
			"high":   ohlcv.HighPrice,
			"low":    ohlcv.LowPrice,
			"close":  ohlcv.ClosePrice,
			"volume": ohlcv.Volume,
		}
		p := influxdb2.NewPoint(
			pair.NormalizedBase,
			tags,
			fields,
			ohlcv.StartTime)
		repo.writeAPI.WritePoint(p)
	}
}

func (repo *RepositoryImpl) GetOHLCVData(
	ctx context.Context,
	base string,
	quote null.String,
	exchange null.String,
	start time.Time,
	end time.Time,
) ([]*models.OHLCVMarketData, error) {
	tags := map[string]string{}
	if quote.Valid {
		tags["quote"] = quote.String
	}
	if exchange.Valid {
		tags["exchange"] = exchange.String
	}
	query := fmt.Sprintf("from(bucket:\"%s\")", repo.ohlcvBucket)
	query = fmt.Sprintf("%s|> range(start: %s, stop: %s)",
		query, start.Format(time.RFC3339), end.Format(time.RFC3339))
	query = fmt.Sprintf("%s|> filter(fn: (r) => r[\"_measurement\"] == \"%s\")", query, base)
	for key, value := range tags {
		query = fmt.Sprintf("%s|> filter(fn: (r) => r[\"%s\"] == \"%s\")", query, key, value)
	}
	query = fmt.Sprintf("%s|> pivot( rowKey:[\"_time\"], columnKey: [\"_field\"],valueColumn: \"_value\")", query)
	// log.Print(query)
	result, err := repo.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	res := []*models.OHLCVMarketData{}
	for result.Next() {
		startTime := result.Record().Start()
		endTime := startTime.Add(time.Minute)
		ohlc := &models.OHLCVMarketData{
			MarketData: models.MarketData{
				Source:        result.Record().ValueByKey("exchange").(string),
				BaseCurrency:  base,
				QuoteCurrency: result.Record().ValueByKey("quote").(string),
			},
			StartTime:  startTime,
			EndTime:    endTime,
			OpenPrice:  result.Record().ValueByKey("open").(float64),
			HighPrice:  result.Record().ValueByKey("high").(float64),
			LowPrice:   result.Record().ValueByKey("low").(float64),
			ClosePrice: result.Record().ValueByKey("close").(float64),
			Volume:     result.Record().ValueByKey("volume").(float64),
		}
		res = append(res, ohlc)
	}
	return res, nil
}
