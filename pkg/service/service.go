package service

import (
	"context"
	"wallet-app/pkg/repository"

	"github.com/google/uuid"
)

type Database interface {
	Close()
	NewWallet(ctx context.Context, walletID uuid.UUID, amount float64) error
	Deposit(ctx context.Context, walletID uuid.UUID, amount float64) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error)
}

type Service struct {
	Database
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Database: repo.Database,
	}
}
