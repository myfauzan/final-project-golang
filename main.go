package main

import (
	"final-project/database"
	"final-project/routes"
	"os"
)

func main() {
	database.StartDB()
	
	var PORT = os.Getenv("PORT")

	routes.StartApp().Run(":" + PORT)
}