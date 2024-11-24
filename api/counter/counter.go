package counter

import (
	"context"
)

// Counter defines the common interface for a counter metric.
type Counter interface {
	// Increment increases the counter by a given amount (default is 1).
	Increment(ctx context.Context, value float64) error

	// Set sets the counter to a specific value.
	Set(ctx context.Context, value float64) error

	// Get retrieves the current value of the counter.
	Get(ctx context.Context) (float64, error)

	// Reset resets the current value of the counter.
	Reset(ctx context.Context) error
}
