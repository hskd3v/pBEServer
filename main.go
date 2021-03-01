package main

import (
	"github.com/harriklein/pBE/pBEServer/app"
	"github.com/harriklein/pBE/pBEServer/db"
	"github.com/harriklein/pBE/pBEServer/ep"
	"github.com/harriklein/pBE/pBEServer/log"
)

func main() {
	log.Init()

	app.Init()

	ep.Init()

	db.Init()

	app.RunServer()

	app.Finish()

}
