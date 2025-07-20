package repository

import (
	"context"

	"github.com/google/uuid"
)

type Database interface {
	Close()
	NewWallet(ctx context.Context, walletID uuid.UUID, amount float64) error
	Deposit(ctx context.Context, walletID uuid.UUID, amount float64) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error)
}

type Repository struct {
	Database
}

func NewRepository(pg *postgresDB) *Repository {
	return &Repository{
		Database: pg,
	}
}
