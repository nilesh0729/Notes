package main

import (
	"database/sql"
	"log"

	_"github.com/lib/pq"
	"github.com/nilesh0729/Notes/api"
	Database "github.com/nilesh0729/Notes/db/Result"
)

const (
	DB_DRIVER = "postgres"
	DB_SOURCE = "postgres://root:MummyJi@127.0.0.1:5431/PapaJi?sslmode=disable"
	SERVER_ADDRESS = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(DB_DRIVER, DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}	
	store := Database.ServerConn(conn)
	server := api.NewServer(*store)

	err = server.Start(SERVER_ADDRESS)
	if err != nil{
		log.Fatal("Cannot Start Server : ", err)
	}
}
