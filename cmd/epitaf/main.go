package main

import (
	"github.com/aureleoules/epitaf/api"
	"github.com/aureleoules/epitaf/db"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	// db.Init()

	go db.Connect()
	api.Serve()
}
