package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	mu          sync.Mutex
	requestsMap = make(map[string]struct{})
)

type ApiServer struct {
	Server *http.Server
}

func NewApiServer(addr string) *ApiServer {
	return &ApiServer{
		Server: &http.Server{
			Addr: addr,
		},
	}
}

func (s *ApiServer) SetupRouter() {
	router := gin.Default()
	s.Server.Handler = router

	// router.Use(TrackRequestsMiddleware())
	router.GET("/api/verve/accept", HandleAccept)
}

func LogUniqueRequests() error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	logFile, err := os.Create("./analytics.txt")
	if err != nil {
		log.Fatalf("Error while creating analytics file:%s\n", err)
	}
	log.Printf("Created analytics file\n")
	defer logFile.Close()

	for range ticker.C {
		currentDateTime := time.Now().Format(time.DateTime)
		logFile.Write([]byte(fmt.Sprintf("[%s] Unique requests: %d\n", currentDateTime, getUniqueRequestsCount())))

		resetUniqueRequests()
	}

	return nil
}
