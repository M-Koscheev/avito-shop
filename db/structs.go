package db

type Merch string

const (
	TShirt    Merch = "t_shirt"
	Cup       Merch = "cup"
	Book      Merch = "book"
	Pen       Merch = "pen"
	Powerbank Merch = "powerbank"
	Hoody Merch =
)

type Authenticate struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type SendCoins struct {
	ToUser string `json:"toUser" db:"toUser"`
	Amount int    `json:"amount" db:"amount" validator:"min=1"`
}
