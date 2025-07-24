package service

import (
	"context"
	"wallet-app/pkg/repository"

	"github.com/google/uuid"
)

type Database interface {
	Close()
	NewWallet(ctx context.Context, walletID uuid.UUID, amount int) error
	Deposit(ctx context.Context, walletID uuid.UUID, amount int) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount int) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (int, error)
}

type Service struct {
	Database
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Database: repo.Database,
	}
}
