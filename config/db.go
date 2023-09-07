package config

import "fmt"

type IDbConfig interface {
	Url() string
	MaxOpenConns() int
}

type db struct {
	host          string
	port          int
	protocal      string
	username      string
	password      string
	database      string
	sslMode       string
	maxConnection int
}

func (c *config) Db() IDbConfig {
	return c.db
}

func (d *db) Url() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.host,
		d.port,
		d.username,
		d.password,
		d.database,
		d.sslMode,
	)
}

func (d *db) MaxOpenConns() int { return d.maxConnection }
