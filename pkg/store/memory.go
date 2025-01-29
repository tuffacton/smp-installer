package store

import (
	"context"
	"fmt"
	"strings"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/util"
	"github.com/rs/zerolog/log"
)

const (
	OutputDirectoryKey = "output_directory"
	ProviderKey        = "provider"
	ProfileKey         = "profile"
)

type memoryStore struct {
	data map[string]interface{}
}

// GetBool implements DataStore.
func (c *memoryStore) GetBool(ctx context.Context, key string) bool {
	val, err := c.Get(ctx, key)
	if err != nil {
		log.Err(err).Msgf("failed to get key %s", key)
		return false
	}
	return val.(bool)
}

// GetString implements DataStore.
func (c *memoryStore) GetString(ctx context.Context, key string) string {
	val, err := c.Get(ctx, key)
	if err != nil {
		log.Err(err).Msgf("failed to get key %s", key)
		return ""
	}
	return val.(string)
}

// DataMap implements DataStore.
func (c *memoryStore) DataMap(ctx context.Context) map[string]interface{} {
	return c.data
}

// Get implements DataStore.
func (c *memoryStore) Get(ctx context.Context, path string) (interface{}, error) {
	segments := strings.Split(path, ".")
	value := c.data
	for idx, seg := range segments {
		_, ok := value[seg]
		if !ok {
			return nil, fmt.Errorf("no path found %s at %s", path, seg)
		}
		if idx == len(segments)-1 {
			return value[seg], nil
		}
		value = value[seg].(map[string]interface{})
	}
	return nil, fmt.Errorf("no path found %s", path)
}

// Set implements DataStore.
func (c *memoryStore) Set(ctx context.Context, path string, value interface{}) error {
	segments := strings.Split(path, ".")
	currentData := c.data
	for idx, seg := range segments {
		// log.Info().Msgf("checking for segment: %s", seg)
		if idx == len(segments)-1 {
			currentData[seg] = value
			return nil
		}
		_, ok := currentData[seg]
		if !ok {
			currentData[seg] = make(map[string]interface{})
		}
		currentData = currentData[seg].(map[string]interface{})
	}
	return nil
}

// Set implements DataStore.
func (c *memoryStore) AddAll(ctx context.Context, data map[string]interface{}) error {
	c.data = util.MergeMaps(data, c.data)
	return nil
}

func NewMemoryStore() DataStore {
	return &memoryStore{
		data: make(map[string]interface{}),
	}
}

func NewMemoryStoreWithData(data map[string]interface{}) DataStore {
	return &memoryStore{
		data: data,
	}
}
