package main

import (
	"log"

	_ "github.com/FaintLocket424/rc-timing-api/docs"
	"github.com/FaintLocket424/rc-timing-api/internal/api"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
)

// @title RC Timing API
// @version 1.0
// @description A web API for live RC timing data.

// @contact.name FaintLocket424
// @contact.url https://github.com/FaintLocket424/rc-timing-api

// @host localhost:8080
// @BasePath /api/v1
func main() {
	store := service.NewStore()
	r := api.SetupRouter(store)

	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		println(err.Error())
		return
	}
}
