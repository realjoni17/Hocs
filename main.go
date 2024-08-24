package main

import (
	"github.com/realjoni17/Hdocs/database"
	"github.com/realjoni17/Hdocs/server"
)

func main() {
	database.StartDatabase()
	server := server.NewServer()
	server.Run()
}

//
