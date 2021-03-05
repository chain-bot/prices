package okex

import (
	"encoding/json"
	"strconv"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// https://www.okex.com/docs/en/#spot-some
//////////////////////////////////////////////////////////////////////////////////////////////////////////
type CandleStickData struct {
	OpenTime   float64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
	Volume     float64
}

func (candleStickResponse *CandleStickData) UnmarshalJSON(
	data []byte,
) (err error) {
	var responseSlice []interface{}
	if err := json.Unmarshal(data, &responseSlice); err != nil {
		return err
	}
	// "2020-01-01T00:00:00.000Z""
	temp, err := time.Parse(time.RFC3339, responseSlice[0].(string))
	if err != nil {
		return err
	}
	candleStickResponse.OpenTime = float64(temp.Unix())
	// "7195.0"
	candleStickResponse.OpenPrice, err = strconv.ParseFloat(responseSlice[1].(string), 64)
	if err != nil {
		return err
	}
	// "7195.0"
	candleStickResponse.HighPrice, err = strconv.ParseFloat(responseSlice[2].(string), 64)
	if err != nil {
		return err
	}
	// "7195.0"
	candleStickResponse.LowPrice, err = strconv.ParseFloat(responseSlice[3].(string), 64)
	if err != nil {
		return err
	}
	// "7195.0"
	candleStickResponse.ClosePrice, err = strconv.ParseFloat(responseSlice[4].(string), 64)
	if err != nil {
		return err
	}
	// "30.32773632"
	candleStickResponse.Volume, err = strconv.ParseFloat(responseSlice[5].(string), 64)
	if err != nil {
		return err
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// https://www.okex.com/docs/en/#spot-some
//////////////////////////////////////////////////////////////////////////////////////////////////////////
type Instrument struct {
	BaseCurrency  string `json:"base_currency"`
	Category      string `json:"category"`
	InstrumentID  string `json:"instrument_id"`
	MinSize       string `json:"min_size"`
	QuoteCurrency string `json:"quote_currency"`
	SizeIncrement string `json:"size_increment"`
	TickSize      string `json:"tick_size"`
}
