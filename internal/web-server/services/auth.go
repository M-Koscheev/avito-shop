package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/M-Koscheev/avito-shop/db"
	"github.com/M-Koscheev/avito-shop/internal/web-server/repository"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	EmployeeUsername string `json:"username"`
}

type AuthService struct {
	repo repository.Employee
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) AuthorizeEmployee(ctx context.Context, input db.AuthRequest) (string, error) {
	user, err := s.repo.GetEmployee(ctx, input.Username)
	h := sha256.New()
	h.Write([]byte(input.Password))
	hs := h.Sum(nil)

	if errors.Is(err, sql.ErrNoRows) {
		// генерируем нового пользователя
		_, err = s.repo.RegisterEmployee(ctx, db.EmployeeInfo{
			Username:     input.Username,
			PasswordHash: hs,
		})
		if err != nil {
			return "", fmt.Errorf("failed to register new employee: %w", err)
		}

	} else if err != nil {
		return "", fmt.Errorf("username checkup failed: %w", err)
	} else {
		// проверяем совпадают ли хэши паролей
		if !bytes.Equal(hs, user.PasswordHash) {
			return "", db.UnauthorizedError{Message: "wrong password"}
		}
	}

	return s.generateToken(input.Username)
}

func (s *AuthService) generateToken(username string) (string, error) {
	// Генерация токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(db.TokenTTL)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		username,
	})

	secret, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		return "", fmt.Errorf("no jwt secret present in env file")
	}

	return token.SignedString([]byte(secret))
}
