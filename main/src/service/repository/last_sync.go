package repository

import "github.com/jmoiron/sqlx"

type LastSyncRepository struct {
	db *sqlx.DB
}

func NewLastSyncRepository(db *sqlx.DB) *LastSyncRepository {
	return &LastSyncRepository{db: db}
}

func (repo *LastSyncRepository) GetAll() {
	print("TODO")
}
