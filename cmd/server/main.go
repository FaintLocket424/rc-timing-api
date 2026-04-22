package main

import (
	"log"

	"github.com/FaintLocket424/rc-timing-api/internal/api"
	"github.com/FaintLocket424/rc-timing-api/internal/storage/cache"
)

func main() {
	version := "0.0.1"
	log.Printf("Starting RC Timing API v%s", version)

	cache := cache.NewCache()
	router := api.SetupRouter(cache)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
