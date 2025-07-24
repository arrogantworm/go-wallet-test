package service

import (
	"context"
	"strconv"

	"github.com/google/uuid"
)

func (s *Service) Withdraw(ctx context.Context, walletID uuid.UUID, amount int) error {
	return s.Database.Withdraw(ctx, walletID, amount)
}

func (s *Service) Deposit(ctx context.Context, walletID uuid.UUID, amount int) error {
	return s.Database.Deposit(ctx, walletID, amount)
}

func (s *Service) GetBalance(ctx context.Context, walletID uuid.UUID) (string, error) {
	balance, err := s.Database.GetBalance(ctx, walletID)
	if err != nil {
		return "", err
	}
	balanceStr := strconv.Itoa(balance)
	balanceStr = balanceStr[:len(balanceStr) -2] + "." + balanceStr[len(balanceStr) - 2:]

	return balanceStr, nil
}
