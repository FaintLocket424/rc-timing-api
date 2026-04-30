package api

import (
	"net/http"

	"github.com/FaintLocket424/rc-timing-api/internal/api/middleware"
	"github.com/FaintLocket424/rc-timing-api/internal/manager"
	"github.com/FaintLocket424/rc-timing-api/internal/storage"
	"github.com/gin-gonic/gin"
)

func SetupRouter(store storage.Store, manager *manager.Manager) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
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

			// quali := results.Group("/qualifying")
			// {
			// 	round := quali.Group("/round/:round")
			// 	{
			// 		round.GET("/heat/:heat", handler.getQualiRaceResult)
			// 	}
			// }

			// finals := results.Group("/finals")
			// {
			// 	round := finals.Group("/round/:round")
			// 	{
			// 		round.GET("/final/:final", handler.GetFinalRaceResult)
			// 	}
			// }
		}

		// localhost:8080/results/practice/round/1/heat/1
	}

	return r
}
