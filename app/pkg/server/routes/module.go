package routes

import (
	"github.com/chain-bot/prices/app/internal/repository"
	"github.com/gorilla/schema"
)

func NewHandler(
	repo repository.Repository,
) *Handler {
	return &Handler{
		repo:    repo,
		decoder: schema.NewDecoder(),
	}
}
