package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *ApiServer) HandleAccept(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

	endpoint := ctx.Query("endpoint")
	// If endpoint is provided, send an HTTP request
	if endpoint != "" {
		go func(endpoint string) {
			uniqueCount, _ := getUniqueRequestsCount(s, ctx)
			data := struct {
				Count     int
				Timestamp string
			}{
				Count:     int(uniqueCount),
				Timestamp: time.Now().Format(time.DateTime),
			}
			jsonBytes, err := json.Marshal(data)
			if err != nil {
				log.Printf("Failed to marshal POST data%v\n", err)
				return
			}
			resp, err := http.Post(endpoint+"?count="+strconv.Itoa(int(uniqueCount)), "application/json", bytes.NewBuffer(jsonBytes))
			if err != nil {
				log.Printf("Failed to send request to %s: %v\n", endpoint, err)
				return
			}
			defer resp.Body.Close()
			log.Printf("Request to %s returned status: %s\n", endpoint, resp.Status)
		}(endpoint)
	}
}
