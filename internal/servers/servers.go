package servers

import (
	"context"
	"fmt"
	"github.com/M-Koscheev/avito-shop/internal/config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.Username, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db does not respond: %w", err)
	}

	return db, nil
}

//func MigrateUp(cfg *config.Config) error {
//	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s",
//		cfg.Host, cfg.Port, cfg.DBName, cfg.Username, cfg.Password, cfg.SSLMode))
//	if err != nil {
//		return fmt.Errorf("failed to open db connection: %w", err)
//	}
//
//	driver, err := postgres.WithInstance(db, &postgres.Config{})
//	if err != nil {
//		return fmt.Errorf("failed to create postgres driver: %w", err)
//	}
//	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
//	if err != nil {
//		return fmt.Errorf("failed to read migration or connect to db: %w", err)
//	}
//	m.Up()
//
//	return nil
//}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler, config *config.Config) error {
	s.httpServer = &http.Server{
		Addr:           config.Server.Address,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    config.Server.Timeout,
		WriteTimeout:   config.Server.Timeout,
		IdleTimeout:    config.Server.IdleTimeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) GracefulShutdown(ctx context.Context) error {
	// TODO graceful
	return s.httpServer.Shutdown(ctx)
}

/*
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.Username, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateUp(cfg Config) error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBName, cfg.Username, cfg.Password, cfg.SSLMode))
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migration",
		"postgres", driver)
	m.Up()
	return nil
}
*/
