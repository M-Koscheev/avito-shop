package repository

import (
	"context"
	"github.com/M-Koscheev/avito-shop/db"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Employee
	Info
}

type Employee interface {
	GetEmployee(ctx context.Context, username string) (db.EmployeeInfo, error)
	RegisterEmployee(ctx context.Context, input db.EmployeeInfo) (db.EmployeeInfo, error)
}

type Info interface {
	GetInventory(ctx context.Context, username string) ([]db.Inventory, error)
	PurchaseProduct(ctx context.Context, username string, product db.Merch) error
	GetTransaction(ctx context.Context, username string) (db.CoinHistory, error)
	SendCoins(ctx context.Context, fromUsername, toUsername string, amount int) error
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Employee: NewEmployeeRepository(db),
		Info:     NewInfoRepository(db),
	}
}
