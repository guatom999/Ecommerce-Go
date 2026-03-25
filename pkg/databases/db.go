package databases

import (
	"log"
	"time"

	"github.com/guatom999/Ecommerce-Go/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// Connect DB
	db, err := sqlx.Connect("pgx", cfg.Url())

	if err != nil {
		log.Printf("connect to db failed cause of %s", err)
	}

	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	db.DB.SetMaxIdleConns(cfg.MaxOpenConns() / 2)
	db.DB.SetConnMaxLifetime(1 * time.Minute)
	db.DB.SetConnMaxIdleTime(1 * time.Minute)

	return db
}
