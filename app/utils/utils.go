package utils

import (
	"encoding/json"
	"github.com/chain-bot/prices/app/pkg/models"
	"log"
	"time"
)

func PrettyJSON(obj interface{}) string {
	json, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(json)
}

func UnixMillis(t time.Time) int64 {
	t.UTC()
	return t.UnixNano() / int64(time.Millisecond)
}

func Reverse(s []*models.OHLCVMarketData) []*models.OHLCVMarketData {
	a := make([]*models.OHLCVMarketData, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}
