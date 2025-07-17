package repository

type Database interface {
	Close()
	Deposit(walletID int, amount float64) error
	Withdraw(walletID int, amount float64) error
	GetBalance(walledID int) (float64, error)
}

type Repository struct {
	Database
}

func NewRepository(pg *postgres) *Repository {
	return &Repository{
		Database: pg,
	}
}
