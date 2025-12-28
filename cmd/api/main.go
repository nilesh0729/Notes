package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nilesh0729/Notes/internal/api"
	Database "github.com/nilesh0729/Notes/internal/db/Result"
	"github.com/nilesh0729/Notes/internal/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot Load Config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}
	store := Database.ServerConn(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot Start Server : ", err)
	}
}
