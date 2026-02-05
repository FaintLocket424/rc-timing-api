package api

import (
	"github.com/FaintLocket424/rc-timing-api/internal/api/handlers"
	"github.com/FaintLocket424/rc-timing-api/internal/api/middleware"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initialises the Gin router with middleware and the handlers
func SetupRouter(store *service.DataStore) *gin.Engine {
	r := gin.Default()

	eventHandler := handlers.NewEventHandler(store)
	practiceHandler := handlers.NewPracticeHandler(store)
	qualifyingHandler := handlers.NewQualifyingHandler(store)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.Use(gzip.Gzip(gzip.DefaultCompression))
		//v1.Use(middleware.NotModifiedMiddleware())
		v1.Use(middleware.ScrapeParamsMiddleware())

		event := v1.Group("/event")
		{
			// Endpoint for reporting broken website scraping
			event.POST("/report", nil)

			// Retrieve metadata about an event:
			// - Number of drivers
			// - Event structure (practice? qualifying? finals?)
			// - Classes present (2wd, 4wd, truck etc.)
			event.GET("/meta", eventHandler.GetMeta)

			// Event schedule
			event.GET("/schedule", nil)

			// Event stats
			// - Fastest lap
			event.GET("/stats", nil)

			// Retrieve drivers list
			event.GET("/drivers/list", nil)

			// Retrieve data about a driver (url encoded name)
			event.GET("/driver", nil)

			practice := event.Group("/practice")
			{
				// Metadata
				// - Number of rounds
				// - Number of heats
				// - Number to count
				practice.GET("/meta", nil)

				// Practice Heats list
				practice.GET("/list", nil)

				// Overall finishing positions for all practice heats
				practice.GET("/overalls", nil)

				// Overall finishing positions for a specific round
				practice.GET("/overalls/round/:id", nil)

				// Routes for a specific practice heat (e.g. Practice 2)
				heat := practice.Group("/heat/:id")
				{
					// Info about a specific round (e.g. Round 2)
					heat.GET("/round/:round_id", practiceHandler.GetHeatRoundResult)
				}
			}

			qualifying := event.Group("/qualifying")
			{
				// Metadata
				// - Number of rounds
				// - Number of heats
				// - Number to count
				qualifying.GET("/meta", nil)

				// Qualifying Heats list
				qualifying.GET("/list", nil)

				// Overall finishing positions for all Qualifying heats
				qualifying.GET("/overalls", nil)

				// Overall finishing positions for a specific round
				qualifying.GET("/overalls/round/:id", nil)

				// Routes for a specific Qualifying heat (e.g. Heat 2)
				heat := qualifying.Group("/heat/:id")
				{
					// Info about a specific round (e.g. Round 2)
					heat.GET("/round/:round_id", qualifyingHandler.GetHeatRoundResult)
				}
			}

			finals := event.Group("/finals")
			{
				// Metadata
				// - Number of legs
				// - Number of finals
				// - Number to count
				finals.GET("/meta", nil)

				// Finals list
				finals.GET("/list", nil)

				// Overall finishing positions for all finals
				finals.GET("/overalls", nil)

				// Routes for a specific final (e.g. A Final)
				final := finals.Group("/final/:id")
				{
					// Overall finishing positions for a final
					final.GET("/overalls", nil)

					// Info about a specific leg (e.g. Leg 2)
					final.GET("/leg/:leg_id", nil)
				}
			}
		}
	}

	return r
}
