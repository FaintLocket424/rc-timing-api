package main

import (
	"log"

	"github.com/FaintLocket424/rc-timing-api/internal/api"
	"github.com/FaintLocket424/rc-timing-api/internal/manager"
	"github.com/FaintLocket424/rc-timing-api/internal/storage/cache"
)

func main() {
	version := "0.0.1"
	log.Printf("Starting RC Timing API v%s", version)

	cache := cache.NewCache()
	manager := manager.NewManager(cache)
	router := api.SetupRouter(cache, manager)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
