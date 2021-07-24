package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chain-bot/prices/app/internal/repository"
	"github.com/chain-bot/prices/app/pkg/models"
	"github.com/gorilla/schema"
	"github.com/volatiletech/null/v8"
)

type Routes struct {
	repo    repository.Repository
	decoder *schema.Decoder
}

func NewServer(
	repo repository.Repository,
) *Routes {
	return &Routes{
		repo:    repo,
		decoder: schema.NewDecoder(),
	}
}

func (s Routes) ping(
	w http.ResponseWriter,
	r *http.Request,
) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		log.Printf("error handling ping, url=%s, err=%s", r.URL, err.Error())
	}
}

type CandleRequest struct {
	Base     string `schema:"base,required"`
	Start    uint64 `schema:"start,required"`
	End      uint64 `schema:"end,required"`
	Quote    string `schema:"quote"`
	Exchange string `schema:"exchange"`
}

type CandleResponse struct {
	OHLCV []models.OHLCVMarketData `schema:"ohlcv,required"`
	Error string                   `schema:"message"`
}

func (s Routes) getCandles(
	w http.ResponseWriter,
	r *http.Request,
) {
	var candleRequest CandleRequest
	err := s.decoder.Decode(&candleRequest, r.URL.Query())
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
	res, err := s.repo.GetOHLCVData(
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
