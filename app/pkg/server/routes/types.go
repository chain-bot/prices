package routes

import (
	"github.com/chain-bot/prices/app/internal/repository"
	"github.com/chain-bot/prices/app/pkg/models"
	"github.com/gorilla/schema"
	"github.com/volatiletech/null/v8"
)

type Handler struct {
	repo    repository.Repository
	decoder *schema.Decoder
}

type CandleRequest struct {
	Base     string `schema:"base,required"`
	Start    uint64 `schema:"start,required"`
	End      uint64 `schema:"end,required"`
	Quote    string `schema:"quote"`
	Exchange string `schema:"exchange"`
}

type CandleResponse struct {
	OHLCV []*models.OHLCVMarketData `schema:"ohlcv,required"`
	Error null.String               `schema:"message"`
}
