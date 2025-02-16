package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	authorizationHeader     = "Authorization"
	contextEmployeeUsername = "employeeUsername"
)

func (h *Handler) employeeIdentity(c *gin.Context) {
	tokenString, err := c.Cookie(authorizationHeader)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("no authorization token present"))
		return
	}

	secret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get secret jwt key"))
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token not transferred: %w", err))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token invalid"))
		return
	}

	//TODO удалить после отладки
	info := "claims: "
	for key, val := range claims {
		info += " | " + fmt.Sprintf("key: %v value: %v", key, val)
	}
	slog.Info(info)

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("current token expired"))
		return
	}

	employeeUsername := claims["username"]
	if employeeUsername == "" {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid employee username"))
		return
	}

	c.Set(contextEmployeeUsername, employeeUsername)
}

func getEmployeeUsername(c *gin.Context) (string, error) {
	username, ok := c.Get(contextEmployeeUsername)
	if !ok {
		return "", fmt.Errorf("no username present in context")
	}

	usernameStr, ok := username.(string)
	if !ok {
		return "", fmt.Errorf("wrong username format")
	}

	return usernameStr, nil
}

//func loggingMiddleware(w http.ResponseWriter, r *http.Request) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//	})
//}
