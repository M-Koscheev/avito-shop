package services

import (
	"context"
	"github.com/M-Koscheev/avito-shop/db"
	"github.com/M-Koscheev/avito-shop/internal/web-server/repository"
)

type Service struct {
	Authentication
	Info
}

type Authentication interface {
	AuthorizeEmployee(ctx context.Context, input db.AuthRequest) (string, error)
}

type Info interface {
	BuyMerch(ctx context.Context, username string, merch db.Merch) error
	SendCoin(ctx context.Context, fromUsername, toUsername string, amount int) error
	EmployeeInfo(ctx context.Context, username string) (db.InfoResponse, error)
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authentication: NewAuthService(*repo),
		Info:           NewInfoService(*repo, *repo),
	}
}
