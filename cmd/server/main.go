package main

import (
	"log"

	"github.com/FaintLocket424/rc-timing-api/internal/routes"
)

func main() {
	version := "0.0.1"
	log.Printf("Starting RC Timing API v%s", version)

	router := routes.SetupRouter()

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
