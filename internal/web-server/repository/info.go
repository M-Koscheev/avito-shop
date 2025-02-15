package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/M-Koscheev/avito-shop/db"
	"github.com/jmoiron/sqlx"
	"time"
)

type InfoRepository struct {
	db *sqlx.DB
}

func NewInfoRepository(db *sqlx.DB) *InfoRepository {
	return &InfoRepository{db: db}
}

func (r *InfoRepository) GetInventory(ctx context.Context, username string) ([]db.Inventory, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT products.title, purchases.amount FROM products 
                     INNER JOIN purchases ON purchases.product_id = products.id
                     WHERE purchases.employee=$1
                     ORDER BY purchases.date`, username)
	if err != nil {
		return nil, fmt.Errorf("failed to select employee inventory: %w", err)
	}

	var inventory []db.Inventory
	for rows.Next() {
		var product db.Inventory
		if err = rows.Scan(&product.Type, &product.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan next product: %w", err)
		}

		inventory = append(inventory, product)
	}

	return inventory, tx.Commit()
}

func (r *InfoRepository) PurchaseProduct(ctx context.Context, username string, product db.Merch) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT balance FROM employees WHERE username=$1 FOR UPDATE`, username)
	var balance int
	if err = row.Scan(&balance); err != nil {
		return fmt.Errorf("failed to select employee balance: %w", err)
	}

	row = tx.QueryRow(`SELECT price FROM products WHERE title=$1::productTitle`, product)
	var productPrice int
	if err = row.Scan(&productPrice); err != nil {
		return fmt.Errorf("failed to select product price: %w", err)
	}

	if productPrice > balance {
		return fmt.Errorf("not enough coind - need %v, but hanve only %v", productPrice, balance)
	}

	res, err := tx.Exec(`UPDATE employees SET balance=balance-$1 WHERE username=$2`, balance-productPrice, username)
	if err != nil {
		return fmt.Errorf("failed to update employee balance: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected info: %w", err)
	}

	if affected != 1 {
		return fmt.Errorf("wrong amount of column affected - needed 1, but affected %v", affected)
	}

	return tx.Commit()
}

func (r *InfoRepository) GetTransaction(ctx context.Context, username string) (db.CoinHistory, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return db.CoinHistory{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT from_employee, to_employee, amount FROM coin_transactions 
                     WHERE from_employee=$1 OR to_employee=$1
                     ORDER BY date`, username)
	if err != nil {
		return db.CoinHistory{}, fmt.Errorf("failed to select employee transactions: %w", err)
	}

	var transactions db.CoinHistory
	for rows.Next() {
		var fromEmployee, toEmployee string
		var amount int
		if err = rows.Scan(&fromEmployee, &toEmployee, &amount); err != nil {
			return db.CoinHistory{}, fmt.Errorf("failed to scan next transaction: %w", err)
		}

		if fromEmployee == username {
			transactions.Sent = append(transactions.Sent, db.TransactionTo{
				ToUser: toEmployee, Amount: amount,
			})
		} else {
			transactions.Received = append(transactions.Received, db.TransactionFrom{
				FromUser: fromEmployee, Amount: amount,
			})
		}
	}

	return transactions, tx.Commit()
}

func (r *InfoRepository) SendCoins(ctx context.Context, fromUsername, toUsername string, amount int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	curDate := time.Now()

	row := tx.QueryRow(`SELECT balance FROM employees WHERE username=$1 FOR UPDATE`, fromUsername)
	var balance int
	if err = row.Scan(&balance); err != nil {
		return fmt.Errorf("failed to select from employee balance: %w", err)
	}

	if amount > balance {
		return fmt.Errorf("not enough coind - need %v, but hanve only %v", amount, balance)
	}

	res, err := tx.Exec(`UPDATE employees SET balance=balance-$1 WHERE username=$2`, amount, fromUsername)
	if err != nil {
		return fmt.Errorf("failed to update employee balance: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected info: %w", err)
	}

	if affected != 1 {
		return fmt.Errorf("wrong amount of column affected from employee - needed 1, but affected %v", affected)
	}

	res, err = tx.Exec(`UPDATE employees SET balance=balance+$1 WHERE username=$2`, amount, toUsername)
	if err != nil {
		return fmt.Errorf("failed to update employee balance: %w", err)
	}

	affected, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected info: %w", err)
	}

	if affected != 1 {
		return fmt.Errorf("wrong amount of column affected to employee - needed 1, but affected %v", affected)
	}

	res, err = tx.Exec(`INSERT INTO coin_transactions (from_employee, to_employee, amount, date) VALUES ($1, $2, $3, $4)`,
		fromUsername, toUsername, amount, curDate)

	affected, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected info: %w", err)
	}

	if affected != 1 {
		return fmt.Errorf("wrong amount of column affected coin_transaction - needed 1, but affected %v", affected)
	}

	return tx.Commit()
}
