package handlers

import (
	_ "avito-shop/docs"
	"avito-shop/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
)

type Handler struct {
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{Services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// TODO узнать нужно ли тут CORSMiddleware

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is up",
		})
	})

	api := router.Group("/api")
	{
		auth := api.POST("/auth")
		buy := api.GET("/buy/:item", h.employeeIdentity)
		sendCoin := api.POST("/sendCoin", h.employeeIdentity)
		info := api.GET("/info", h.employeeIdentity)
	}

	return router
}
