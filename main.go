package main

import (
	"final-project/database"
	"final-project/routes"
)

func main() {
	database.StartDB()
	r := routes.StartApp()
	r.Run(":8080")
}