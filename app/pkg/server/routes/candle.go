package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/docker/distribution/context"
	"github.com/volatiletech/null/v8"
)

func (h Handler) GetCandles(
	w http.ResponseWriter,
	r *http.Request,
) {
	var candleRequest CandleRequest
	err := h.decoder.Decode(&candleRequest, r.URL.Query())
	if err != nil {
		log.Printf("decode, err=%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	endTime := time.Unix(int64(candleRequest.End), 0)
	startTime := time.Unix(int64(candleRequest.Start), 0)
	if endTime.Before(startTime) {
		log.Printf("end before start")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := h.repo.GetOHLCVData(
		context.WithRequest(context.Background(), r),
		candleRequest.Base,
		null.StringFrom(candleRequest.Quote),
		null.StringFrom(candleRequest.Exchange),
		startTime,
		endTime,
	)
	if err != nil {
		log.Printf("GetOHLCVData, err=%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("res=%s", res)
}
