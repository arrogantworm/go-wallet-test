package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) error {
	return s.Database.Withdraw(ctx, walletID, amount)
}

func (s *Service) Deposit(ctx context.Context, walletID uuid.UUID, amount float64) error {
	return s.Database.Deposit(ctx, walletID, amount)
}

func (s *Service) GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	return s.Database.GetBalance(ctx, walletID)
}
