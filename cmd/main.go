package main

import (
	"context"
	"github.com/M-Koscheev/avito-shop/internal/config"
	"github.com/M-Koscheev/avito-shop/internal/servers"
	"github.com/M-Koscheev/avito-shop/internal/web-server/handlers"
	"github.com/M-Koscheev/avito-shop/internal/web-server/repository"
	"github.com/M-Koscheev/avito-shop/internal/web-server/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// @title           Avito-shop API
// @version         1.0
// @description     API for inner Avito merch shop.

// @contact.name   API Support
// @contact.email  mr.kosheef54@gmail.com

// @host      localhost:8080
// @BasePath  /api
// @Schemes http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	slog.Info("starting web-server")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading env: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	//slog.Info("test info log")
	//slog.Error("test error log")
	cfg := config.Load()

	cfg.DB.Password = os.Getenv("POSTGRES_PASSWORD")

	slog.Info("config env variables", "cfg", cfg)

	db, err := servers.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("failed to init db", "error", err)
		os.Exit(1)
	}

	//err = servers.MigrateUp(cfg)
	//if err != nil {
	//	slog.Error("failed to migrate db", "error", err)
	//	os.Exit(1)
	//}

	slog.Info("connected to db and migrated")

	repositoryVar := repository.NewRepository(db)
	servicesVar := services.NewService(repositoryVar)
	handlersVar := handlers.NewHandler(servicesVar)

	srv := servers.Server{}
	exit := make(chan os.Signal, 1)

	go func() {
		if err = srv.Run(handlersVar.InitRoutes(), cfg); err != nil {
			slog.Error("error while running server", "error", err)
		}
	}()

	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	if err := srv.GracefulShutdown(context.Background()); err != nil {
		slog.Error("error occurred on server shutting down")
	}

	if err := db.Close(); err != nil {
		slog.Error("error occurred on db connection close")
	}

	slog.Error("server stopper")
}
