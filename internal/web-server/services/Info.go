package services

import (
	"context"
	"fmt"
	"github.com/M-Koscheev/avito-shop/db"
	"github.com/M-Koscheev/avito-shop/internal/web-server/repository"
	"golang.org/x/sync/errgroup"
)

type InfoService struct {
	infoRepo     repository.Info
	employeeRepo repository.Employee
}

func NewInfoService(infoRepo repository.Repository, employeeRepo repository.Repository) *InfoService {
	return &InfoService{infoRepo: infoRepo, employeeRepo: employeeRepo}
}

func (s *InfoService) BuyMerch(ctx context.Context, username string, merch db.Merch) error {
	_, err := s.employeeRepo.GetEmployee(ctx, username)
	if err != nil {
		return db.InvalidRequestError{Message: fmt.Sprintf("failed to get employee with given username: %v", err)}
	}

	if err = s.infoRepo.PurchaseProduct(ctx, username, merch); err != nil {
		return fmt.Errorf("failed to buy %v for employee %v: %w", merch, username, err)
	}

	return nil
}

func (s *InfoService) SendCoin(ctx context.Context, fromUsername, toUsername string, amount int) error {
	var wg errgroup.Group
	for _, username := range []string{fromUsername, toUsername} {
		wg.Go(func() error {
			_, err := s.employeeRepo.GetEmployee(ctx, username)
			if err != nil {
				return db.InvalidRequestError{Message: fmt.Sprintf("failed to get employee with given username: %v", err)}
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	if err := s.infoRepo.SendCoins(ctx, fromUsername, toUsername, amount); err != nil {
		return fmt.Errorf("failed to send %v coins from employee %v to employee %v: %w", amount, fromUsername, toUsername, err)
	}

	return nil
}

func (s *InfoService) EmployeeInfo(ctx context.Context, username string) (db.InfoResponse, error) {
	info := db.InfoResponse{}
	employee, err := s.employeeRepo.GetEmployee(ctx, username)
	if err != nil {
		return db.InfoResponse{}, fmt.Errorf("failed to get employee %v info: %w", username, err)
	}

	info.Coins = employee.Balance

	inventory, err := s.infoRepo.GetInventory(ctx, username)
	if err != nil {
		return db.InfoResponse{}, fmt.Errorf("failed to get employee %v inventory: %w", username, err)
	}

	info.Inventory = inventory

	transactions, err := s.infoRepo.GetTransaction(ctx, username)
	if err != nil {
		return db.InfoResponse{}, fmt.Errorf("failed to get employee %v transactions history: %w", username, err)
	}

	info.CoinHistory = transactions

	return info, nil
}
