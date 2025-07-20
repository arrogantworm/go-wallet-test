package models

import "github.com/google/uuid"

type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
}
