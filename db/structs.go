package db

import "time"

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
