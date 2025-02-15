package handlers

import (
	_ "github.com/M-Koscheev/avito-shop/docs"
	"github.com/M-Koscheev/avito-shop/internal/web-server/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
)

type Handler struct {
	Services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{Services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// TODO узнать нужно ли тут CORSMiddleware

	// TODO graceful shutdown

	// TODO logger

	// TODO make repository private

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is up",
		})
	})

	api := router.Group("/api")
	{
		api.POST("/auth", h.authenticateEmployee)
		api.GET("/buy/:item", h.employeeIdentity, h.buyMerch)
		api.POST("/sendCoin", h.employeeIdentity, h.sendCoin)
		api.GET("/info", h.employeeIdentity, h.getInfo)
	}

	return router
}
