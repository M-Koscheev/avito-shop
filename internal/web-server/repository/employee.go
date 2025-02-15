package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/M-Koscheev/avito-shop/db"
	"github.com/jmoiron/sqlx"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(db *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) GetEmployee(ctx context.Context, username string) (db.EmployeeInfo, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return db.EmployeeInfo{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT username, password_hash, balance FROM employees WHERE username=$1`, username)
	var employee db.EmployeeInfo
	if err = row.Scan(&employee.Username, &employee.PasswordHash, &employee.Balance); err != nil {
		return db.EmployeeInfo{}, fmt.Errorf("failed to select employee data: %w", err)
	}

	return employee, tx.Commit()
}

func (r *EmployeeRepository) RegisterEmployee(ctx context.Context, input db.EmployeeInfo) (db.EmployeeInfo, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return db.EmployeeInfo{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`INSERT INTO employees (username, password_hash) VALUES ($1, $2) RETURNING balance`, input.Username, input.PasswordHash)
	if err = row.Scan(&input.Balance); err != nil {
		return db.EmployeeInfo{}, fmt.Errorf("failed to register employee: %w", err)
	}

	return input, tx.Commit()
}
