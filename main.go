package main

import (
	"finalproject/database"
	"finalproject/routers"
)

func main() {
	var port = ":8080"
	database.ConnectToDatabase()
	routers.StartServer().Run(port)
}
