// Package api handles the Gin router, taking incoming requests and using the
// [internal/manager] and [internal/storage] packages to retrieve the response
// and then sending it back to the user.
package api

import (
	"net/http"

	"github.com/FaintLocket424/opengrid-bridge/internal/api/middleware"
	"github.com/FaintLocket424/opengrid-bridge/internal/manager"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage"
	"github.com/gin-gonic/gin"
)

// SetupRouter creates the Gin router with a rate limiter middleware from
// [internal/api/middleware]. It also initialises the Handler struct which
// contains the methods to process the incoming requests.
func SetupRouter(store storage.Store, manager *manager.Manager) *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	r.Use(middleware.RateLimiter(5, 10))

	handler := NewHandler(store, manager)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		v1.GET("/live", handler.GetLiveTiming)

		results := v1.Group("/results")
		{
			practice := results.Group("/practice")
			{
				round := practice.Group("/round/:round")
				{
					round.GET("/heat/:heat", handler.GetPracticeRaceResult)
				}
			}

			quali := results.Group("/qualifying")
			{
				round := quali.Group("/round/:round")
				{
					round.GET("/heat/:heat", handler.GetQualiRaceResult)
				}
			}

			finals := results.Group("/finals")
			{
				round := finals.Group("/round/:round")
				{
					round.GET("/final/:final", handler.GetFinalRaceResult)
				}
			}
		}
	}

	return r
}
