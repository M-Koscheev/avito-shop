package service

import "avito-shop/internal/repository"

type Service struct {
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
