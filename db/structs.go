package db

import (
	"fmt"
	"time"
)

type Merch string

const (
	TShirt    Merch = "t-shirt"
	Cup       Merch = "cup"
	Book      Merch = "book"
	Pen       Merch = "pen"
	Powerbank Merch = "powerbank"
	Hoody     Merch = "hoody"
	Umbrella  Merch = "umbrella"
	Socks     Merch = "socks"
	Wallet    Merch = "wallet"
	PinkHoody Merch = "pink-hoody"
)

var merchMap = map[Merch]struct{}{
	TShirt:    {},
	Cup:       {},
	Book:      {},
	Pen:       {},
	Powerbank: {},
	Hoody:     {},
	Umbrella:  {},
	Socks:     {},
	Wallet:    {},
	PinkHoody: {},
}

func ParseMerch(given string) (Merch, error) {
	merch := Merch(given)
	_, ok := merchMap[merch]
	if !ok {
		return "", fmt.Errorf("invalid merch title given")
	}

	return merch, nil
}

const (
	TokenTTL = 12 * time.Hour
)

type AuthRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type ErrorResponse struct {
	Error string `json:"errors"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type Inventory struct {
	Type     Merch `json:"type"`
	Quantity int   `json:"quantity"`
}

type TransactionFrom struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type TransactionTo struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []TransactionFrom `json:"received"`
	Sent     []TransactionTo   `json:"sent"`
}

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type EmployeeInfo struct {
	Username     string `db:"username"`
	PasswordHash []byte `db:"password_hash"`
	Balance      int    ` db:"balance"`
}

type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

type InvalidRequestError struct {
	Message string
}

func (e InvalidRequestError) Error() string {
	return e.Message
}
