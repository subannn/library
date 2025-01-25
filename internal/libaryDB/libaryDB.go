package libaryDB

import (
	"EffectiveMobileTestTask/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
)

type DB struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewDB(logger *slog.Logger, config *config.Config) *DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBConfig.Host, config.DBConfig.Port, config.DBConfig.Username, config.DBConfig.Password, config.DBConfig.Database, config.DBConfig.SSLMode,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	logger.Info("Connected to database")
	return &DB{
		db:     db,
		logger: logger,
	}
}

func (db *DB) Shutdown() {
	if err := db.db.Close(); err != nil {
		db.logger.Debug("Error while closing the database connection", "error", err)
	} else {
		db.logger.Info("Database connection closed successfully")
	}
}
