package service

import (
	"wallet-app/pkg/repository"
)

type Database interface {
	Close()
	Deposit(walletID int, amount float64) error
	Withdraw(walletID int, amount float64) error
	GetBalance(walledID int) (float64, error)
}

type Service struct {
	Database
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Database: repo.Database,
	}
}
