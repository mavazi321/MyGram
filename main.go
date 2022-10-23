package main

import (
	"finalproject/database"
	"finalproject/routers"
	"os"
)

func main() {
	var port = os.Getenv("PORT")
	database.ConnectToDatabase()
	routers.StartServer().Run(":" + port)
}
