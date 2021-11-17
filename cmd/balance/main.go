package main

import (
	"avito-intership/server"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	app := server.NewApp()

	if err := app.Run(":" + port); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
