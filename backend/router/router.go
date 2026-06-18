package router

import (
	"ocean-ranch/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(
	cageHandler *handlers.CageHandler,
	feedHandler *handlers.FeedHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/cages", cageHandler.GetAll)
		api.GET("/cages/:id", cageHandler.GetByID)
		api.PUT("/cages/:id", cageHandler.UpdateStatus)

		api.POST("/feed", feedHandler.Create)
		api.GET("/feed/records", feedHandler.GetAll)
		api.GET("/feed/records/:cage_id", feedHandler.GetByCageID)
	}

	return r
}
