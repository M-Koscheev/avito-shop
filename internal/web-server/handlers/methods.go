package handlers

import (
	"errors"
	"fmt"
	"github.com/M-Koscheev/avito-shop/db"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// @Summary		Аутентификация и получение JWT-токена.
// @Tags			auth
// @Description	Аутентификация пользователя для дальнейшего использования сервиса
// @ID	auth
// @Accept			json
// @Produce		json
// @Param			input	body		db.AuthRequest	true	"Данные для аутентификации."
// @Success		200		{object}	db.AuthResponse "Аутентификация и получение JWT-токена."
// @Failure		400	{object}	db.ErrorResponse "Неверный запрос."
// @Failure		401	{object}	db.ErrorResponse "Неавторизован."
// @Failure		500		{object}	db.ErrorResponse "Внутренняя ошибка сервера."
// @Router			/auth [post]
func (h *Handler) authenticateEmployee(c *gin.Context) {
	ctx := c.Request.Context()
	slog.Info("authentication handler", "accepted", true)
	var input db.AuthRequest
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(400, fmt.Errorf("неверный запрос"))
		return
	}

	token, err := h.Services.Authentication.AuthorizeEmployee(ctx, input)
	slog.Info("authentication handler", "token", token, "error", err)
	var unauthorizeErr db.UnauthorizedError
	if errors.As(err, &unauthorizeErr) {
		c.AbortWithError(401, err)
		return
	} else if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(authorizationHeader, token, int(db.TokenTTL.Seconds()), "", "", true, false)

	c.JSON(200, db.AuthResponse{
		Token: token,
	})
}

// @Summary		Получить информацию о монетах, инвентаре и истории транзакций.
// @Tags	info
// @Description Получить информацию о монетах, инвентаре и истории транзакций авторизированного сотрудника.
// @ID	info
// @Produce		json
// @Success		200		{object}	db.InfoResponse "Успешный ответ."
// @Failure		400	{object}	db.ErrorResponse "Неверный запрос."
// @Failure		401	{object}	db.ErrorResponse "Неавторизован."
// @Failure		500		{object}	db.ErrorResponse "Внутренняя ошибка сервера."
// @Router			/info [get]
func (h *Handler) getInfo(c *gin.Context) {
	ctx := c.Request.Context()
	slog.Info("get info handler", "accepted", true)
	employeeUsername, err := getEmployeeUsername(c)
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	info, err := h.Services.Info.EmployeeInfo(ctx, employeeUsername)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, info)
}

// @Summary		Отправить монеты другому пользователю.
// @Tags			coins
// @Description	Отправить монеты указанному пользователю.
// @ID	send-coins
// @Accept			json
// @Produce		json
// @Param			input	body		db.SendCoinRequest	true	"Данные о пользователе и количестве монет."
// @Success		200		{string}	string "Успешный ответ."
// @Failure		400	{object}	db.ErrorResponse "Неверный запрос."
// @Failure		401	{object}	db.ErrorResponse "Неавторизован."
// @Failure		500		{object}	db.ErrorResponse "Внутренняя ошибка сервера."
// @Router			/sendCoin [post]
func (h *Handler) sendCoin(c *gin.Context) {
	ctx := c.Request.Context()
	slog.Info("send info handler", "accepted", true)
	employeeUsername, err := getEmployeeUsername(c)
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	var req db.SendCoinRequest
	if err = c.BindJSON(&req); err != nil {
		c.AbortWithError(400, fmt.Errorf("error getting input body from request"))
		return
	}

	err = h.Services.Info.SendCoin(ctx, employeeUsername, req.ToUser, req.Amount)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, "Монеты отправлены.")
}

// @Summary		Купить предмет за монеты.
// @Tags			items
// @Description	Купить предмет за монеты.
// @ID	buy-item
// @Accept			json
// @Produce		json
// @Param			item	path		db.Merch	true	"Название предмета."
// @Success		200		{string}	string "Успешный ответ."
// @Failure		400	{object}	db.ErrorResponse "Неверный запрос."
// @Failure		401	{object}	db.ErrorResponse "Неавторизован."
// @Failure		500		{object}	db.ErrorResponse "Внутренняя ошибка сервера."
// @Router			/buy/{item} [get]
func (h *Handler) buyMerch(c *gin.Context) {
	ctx := c.Request.Context()
	slog.Info("buy merch handler", "accepted", true)
	employeeUsername, err := getEmployeeUsername(c)
	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	var merch db.Merch
	if err = c.BindJSON(&merch); err != nil {
		c.AbortWithError(400, fmt.Errorf("error getting input body from request"))
		return
	}

	err = h.Services.Info.BuyMerch(ctx, employeeUsername, merch)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, "Предмет куплен.")
}
