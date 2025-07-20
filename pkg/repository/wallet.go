package repository

import (
	"context"
	"fmt"
	custom_errors "wallet-app/pkg/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (pg *postgresDB) NewWallet(ctx context.Context, walletID uuid.UUID, amount float64) error {
	query := `INSERT INTO wallets (id, balance) VALUES (@walletID, @amount)`
	args := pgx.NamedArgs{"walletID": walletID, "amount": amount}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("create wallet: %w", err)
	}
	return nil
}

func (pg *postgresDB) Deposit(ctx context.Context, walletID uuid.UUID, amount float64) error {
	querySelect := `SELECT balance FROM wallets WHERE id=@walletID FOR UPDATE`
	argsSelect := pgx.NamedArgs{
		"walletID": walletID,
	}

	queryUpdate := `UPDATE wallets SET balance = balance + @amount WHERE id=@walletID`
	argsUpdate := pgx.NamedArgs{
		"walletID": walletID,
		"amount":   amount,
	}

	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var balance float64

	if err := tx.QueryRow(ctx, querySelect, argsSelect).Scan(&balance); err != nil {
		return fmt.Errorf("select for update: %w", err)
	}

	_, err = tx.Exec(ctx, queryUpdate, argsUpdate)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	return tx.Commit(ctx)
}

func (pg *postgresDB) Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) error {
	querySelect := `SELECT balance FROM wallets WHERE id=@walletID FOR UPDATE`
	argsSelect := pgx.NamedArgs{
		"walletID": walletID,
	}

	queryUpdate := `UPDATE wallets SET balance = balance - @amount WHERE id=@walletID`
	argsUpdate := pgx.NamedArgs{
		"walletID": walletID,
		"amount":   amount,
	}

	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var balance float64

	if err := tx.QueryRow(ctx, querySelect, argsSelect).Scan(&balance); err != nil {
		return fmt.Errorf("select for update: %w", err)
	}

	if balance < amount {
		return custom_errors.ErrNotEnoughFunds
	}

	_, err = tx.Exec(ctx, queryUpdate, argsUpdate)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	return tx.Commit(ctx)
}

func (pg *postgresDB) GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	var balance float64
	err := pg.db.QueryRow(ctx, `SELECT balance FROM wallets WHERE id = $1`, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}
	return balance, nil
}
