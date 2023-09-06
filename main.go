package main

import (
	"os"

	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/servers"
	"github.com/guatom999/Ecommerce-Go/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	// fmt.Println(cfg.Db())

	db := databases.DbConnect(cfg.Db())

	defer db.Close()

	// fmt.Println(db)

	servers.NewServer(cfg, db).Start()

}
