package db

import (
	"context"
)

type DB interface {
	Connect() error
	Close() error
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}) (interface{}, error)
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
}
