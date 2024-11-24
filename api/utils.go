package api

import (
	"context"
)

var (
	UNIQUE_REQUESTS = "unique_requests"
)

// getUniqueRequestsCount returns the unique request count.
func getUniqueRequestsCount(server *ApiServer, ctx context.Context) (int64, error) {
	count, err := server.dbClient.Get(ctx, UNIQUE_REQUESTS)
	if err != nil {
		return 0, err
	}

	return count.(int64), nil
}

// resetUniqueRequests resets the unique request count in redis and in prometheus
// for the next minute interval.
func resetUniqueRequests(server *ApiServer, ctx context.Context) error {
	err := server.dbClient.Delete(ctx, UNIQUE_REQUESTS)
	if err != nil {
		return err
	}
	server.uniqueRequestCounter.Reset(ctx)
	return nil
}

// updateRequestTracking increments the unique counter by 1 in redis and prometheus.
func updateRequestTracking(server *ApiServer, ctx context.Context, id string) error {
	status, err := server.dbClient.Set(ctx, UNIQUE_REQUESTS, id)
	if err != nil {
		return err
	}

	if status.(int64) == 1 {
		server.uniqueRequestCounter.Increment(ctx, 0)
	}
	return nil
}
