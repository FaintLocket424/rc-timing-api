// Package api handles the Gin router, taking incoming requests and using the
// [internal/tracking] and [internal/storage] packages to retrieve the response
// and then sending it back to the user.
package api

import (
	"fmt"
	"net/http"

	"github.com/FaintLocket424/opengrid-bridge/internal/api/responses"
	"github.com/FaintLocket424/opengrid-bridge/internal/storage"
	"github.com/gin-gonic/gin"
)

// SetupRouter creates the Gin router with a rate limiter middleware from
// [internal/api/middleware]. It also initialises the Handler struct which
// contains the methods to process the incoming requests.
func SetupRouter(store storage.Store, tracker Tracker, programVersion string) *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	r.Use(RateLimiter(5, 10))

	handler := NewHandler(store, tracker)

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			})

			v1.GET("/info", func(c *gin.Context) {
				responses.RespondSuccess(
					c, nil,
					fmt.Sprintf("API v1 powered by OpenGrid Timing Bridge %s, licensed under the AGPL-3.0 license. Source available at https://codeberg.org/OpenGrid-RC/bridge", programVersion),
				)
			})

			v1.Use(ExtractTargetURL(handler))

			v1.GET("/live", handler.GetLiveTiming)

			v1.GET("/schedule", handler.GetRaceSchedule)

			results := v1.Group("/results")
			{
				results.GET("/", handler.GetRaceResultsIndex)

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
	}

	return r
}
