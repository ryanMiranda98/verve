package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TrackRequestsMiddleware checks for valid id in the request and updates the unique request
// count in redis.
func TrackRequestsMiddleware(server *ApiServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Query("id")
		if id == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "failed",
				"errors":  []string{"query parameter `id` missing"},
			})
			return
		}

		updateRequestTracking(server, ctx, id)
		ctx.Next()
	}
}
