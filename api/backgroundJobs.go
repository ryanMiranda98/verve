package api

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type BackgroundJobs struct {
	server    *ApiServer
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
	ticker    *time.Ticker
}

func NewBackgroundJobs(server *ApiServer) *BackgroundJobs {
	return &BackgroundJobs{
		server:   server,
		stopChan: make(chan struct{}),
		ticker:   time.NewTicker(1 * time.Minute),
	}
}

func (b *BackgroundJobs) Start() {
	b.waitGroup.Add(1)
	go b.runJobs()
}

func (b *BackgroundJobs) Stop() {
	close(b.stopChan)
	b.ticker.Stop()
	b.waitGroup.Wait()
}

func (b *BackgroundJobs) runJobs() {
	defer b.waitGroup.Done()
	for {
		select {
		case <-b.ticker.C:
			b.PublishUniqueRequestCount()
			b.ResetUniqueRequestCount()
		case <-b.stopChan:
			return
		}
	}
}

func (b *BackgroundJobs) PublishUniqueRequestCount() error {
	count, err := getUniqueRequestsCount(b.server, context.Background())
	if err != nil {
		return err
	}
	type Message struct {
		Timestamp string `json:"timestamp"`
		Message   string `json:"message"`
	}
	msg := Message{
		Message:   fmt.Sprintf("Unique Request Count: %d", count),
		Timestamp: time.Now().Format(time.DateTime),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = b.server.writer.Write(msgBytes)
	return err
}

// PushAndResetUniqueRequestCount pushes the unique request count as a log to RabbitMQ
// and resets the count in redis and prometheus every minute.
func (b *BackgroundJobs) ResetUniqueRequestCount() error {
	return resetUniqueRequests(b.server, context.Background())
}
