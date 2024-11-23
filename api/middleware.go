package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TrackRequestsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "failed",
				"errors":  []string{"query parameter `id` missing"},
			})
			return
		}

		updateRequestTracking(id)
		c.Next()
	}
}
