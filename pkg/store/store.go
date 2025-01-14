package store

import "context"

type DataStore interface {
	Get(ctx context.Context, key string) (interface{}, error)
	GetString(ctx context.Context, key string) string
	GetBool(ctx context.Context, key string) bool
	DataMap(ctx context.Context) map[string]interface{}

	Set(ctx context.Context, key string, value interface{}) error
}
