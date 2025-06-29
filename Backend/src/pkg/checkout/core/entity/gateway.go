package entity

import "time"

type AccountType string

const (
	SOCIALPAY AccountType = "SOCIALPAY"
	BANK      AccountType = "BANK"
	CARD      AccountType = "CARD"
	WALLET    AccountType = "WALLET"
)

type Gateway struct {
	Id      string
	Key     string
	Name    string
	Acronym string

	Icon string
	Type AccountType

	CanProcess bool
	CanSettle  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
