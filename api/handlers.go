package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleAccept(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed",
			"errors":  []string{"query parameter `id` missing"},
		})
		return
	}

	updateRequestTracking(id)

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

	// If endpoint is provided, send an HTTP request
	endpoint := c.Query("endpoint")
	if endpoint != "" {
		go func(endpoint string) {
			uniqueCount := getUniqueRequestsCount()
			resp, err := http.Get(endpoint + "?count=" + strconv.Itoa(uniqueCount))
			if err != nil {
				log.Printf("Failed to send request to %s: %v", endpoint, err)
				return
			}
			defer resp.Body.Close()
			log.Printf("Request to %s returned status: %s", endpoint, resp.Status)
		}(endpoint)
	}

	return
}
