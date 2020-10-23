package main

import (
	"github.com/aureleoules/epitaf/api"
	"github.com/aureleoules/epitaf/db"
)

func main() {
	go db.Connect()
	api.Serve()
}
