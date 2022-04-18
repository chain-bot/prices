package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/distribution/context"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
)

func (h Handler) GetCandles(
	w http.ResponseWriter,
	r *http.Request,
) {
	var candleRequest CandleRequest
	w.Header().Set("Content-Type", "application/json")
	res := CandleResponse{}
	err := h.decoder.Decode(&candleRequest, r.URL.Query())
	if err != nil {
		res.Error = null.StringFrom(fmt.Sprintf("invalid Query, err=%s", err.Error()))
		log.WithField("err", err.Error()).Errorf("decode")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	endTime := time.Unix(int64(candleRequest.End), 0)
	startTime := time.Unix(int64(candleRequest.Start), 0)
	if endTime.Before(startTime) {
		res.Error = null.StringFrom(fmt.Sprintf("invalid Query, err=end before start"))
		log.Printf("end before start")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	if endTime.Sub(startTime) > time.Hour*24 {
		res.Error = null.StringFrom(fmt.Sprintf("invalid Query, err=time window greater than 1 day"))
		log.Printf("time window greater than 1 day")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	ohlcv, err := h.repo.GetOHLCVData(
		context.WithRequest(context.Background(), r),
		strings.ToUpper(candleRequest.Base),
		null.StringFrom(strings.ToUpper(candleRequest.Quote)),
		null.StringFrom(strings.ToUpper(candleRequest.Exchange)),
		startTime,
		endTime,
	)
	if err != nil {
		res.Error = null.StringFrom(fmt.Sprintf("internal server error"))
		log.WithField("err", err.Error()).Errorf("GetOHLCVData")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	res.OHLCV = ohlcv
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
