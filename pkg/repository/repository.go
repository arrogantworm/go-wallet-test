package repository

import (
	"context"

	"github.com/google/uuid"
)

type Database interface {
	Close()
	NewWallet(ctx context.Context, walletID uuid.UUID, amount int) error
	Deposit(ctx context.Context, walletID uuid.UUID, amount int) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount int) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (int, error)
}

type Repository struct {
	Database
}

func NewRepository(pg *postgresDB) *Repository {
	return &Repository{
		Database: pg,
	}
}
