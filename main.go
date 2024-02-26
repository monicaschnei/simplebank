package main

import (
	"database/sql"
	"log"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"github.com/monicaschnei/simplebank/api"
	db "github.com/monicaschnei/simplebank/db/sqlc"
	"github.com/monicaschnei/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
