package custom_errors

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotEnoughFunds = errors.New("not enough funds")
	ErrWalletNotFound = pgx.ErrNoRows
)
