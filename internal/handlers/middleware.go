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
	authorizationHeader = "Authorization"
	contextEmployeeId   = "employeeId"
)

func (h *Handler) employeeIdentity(c *gin.Context) {
	tokenString, err := c.Cookie(authorizationHeader)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("no authorization token present"))
		return
	}

	secret, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get secret jwt key"))
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token not transferred"))
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

	employeeId := int(claims["id"].(float64))
	if employeeId <= 0 {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid employee ID"))
		return
	}

	c.Set(contextEmployeeId, employeeId)
}

func getEmployeeId(c *gin.Context) (int, error) {
	idString, ok := c.Get(contextEmployeeId)
	if !ok {
		return 0, fmt.Errorf("no id present in context")
	}

	id, ok := idString.(int)
	if !ok {
		return 0, fmt.Errorf("wrong id format")
	}

	return id, nil
}
