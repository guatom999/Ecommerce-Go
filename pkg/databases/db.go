package databases

import (
	"log"

	"github.com/guatom999/Ecommerce-Go/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// Connect DB
	db, err := sqlx.Connect("pgx", cfg.Url())

	if err != nil {
		log.Fatalf("connect to db failed cause of %s", err)
	}

	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())

	return db
}
