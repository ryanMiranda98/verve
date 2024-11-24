package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ryanMiranda98/verve/api/counter"
	"github.com/ryanMiranda98/verve/api/db"
	"github.com/ryanMiranda98/verve/api/mq"
)

var (
	REDIS_SRV  = os.Getenv("redis-srv")
	REDIS_PORT = os.Getenv("redis-port")
	MQ_USER    = os.Getenv("mq-user")
	MQ_PASS    = os.Getenv("mq-pass")
	MQ_SRV     = os.Getenv("mq-srv")
)

type ApiServer struct {
	httpServer           *http.Server
	stopChan             chan struct{}
	bgJobs               *BackgroundJobs
	dbClient             db.DB
	uniqueRequestCounter counter.Counter
	writer               io.WriteCloser
}

// Create a new instance of API Server
func NewApiServer(addr string) *ApiServer {
	router := gin.Default()
	server := &ApiServer{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		stopChan: make(chan struct{}),
	}
	server.bgJobs = NewBackgroundJobs(server)
	server.uniqueRequestCounter = counter.NewPrometheusCounter(
		"http_request_get_unique_request_count",
		"Count of unique requests",
		[]string{"requests"},
	)

	// Connect to redis
	redisClient := db.NewRedisDB(fmt.Sprintf("%s:%s", REDIS_SRV, REDIS_PORT))
	err := redisClient.Connect()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to DB")
	server.dbClient = redisClient

	// Connect to RabbitMQ
	rabbitMq := mq.NewRabbitMQ(MQ_USER, MQ_PASS, MQ_SRV)
	err = rabbitMq.SetupConnection()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to MQ")
	// file, _ := os.OpenFile("", os.O_APPEND, 0644)
	server.writer = rabbitMq

	// Routes
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Use(TrackRequestsMiddleware(server))
	router.GET("/api/verve/accept", server.HandleAccept)

	return server
}

// Start starts the API server and background tasks
func (s *ApiServer) Start() error {
	s.bgJobs.Start()
	return s.httpServer.ListenAndServe()
}

// Shutdown terminates the API server and its connections.
func (s *ApiServer) Shutdown(ctx context.Context) error {
	if err := s.dbClient.Close(); err != nil {
		return err
	}
	if err := s.writer.Close(); err != nil {
		return err
	}
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
