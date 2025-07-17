package repository

func (pg *postgres) Deposit(walletID int, amount float64) error {
	return nil
}

func (pg *postgres) Withdraw(walletID int, amount float64) error {
	return nil
}

func (pg *postgres) GetBalance(walledID int) (float64, error) {
	return 1000, nil
}
