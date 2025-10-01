package repositories

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheManager struct {
	rdb *redis.Client
}

func NewCacheManager(rdb *redis.Client) *CacheManager {
	return &CacheManager{rdb: rdb}
}

// GetFromCache attempts to retrieve and unmarshal data from Redis cache
// Returns true if cache hit, false if cache miss
func (c *CacheManager) GetFromCache(ctx context.Context, key string, dest interface{}) bool {
	cmd := c.rdb.Get(ctx, key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			log.Printf("Key %s does not exist\n", key)
		} else {
			log.Println("Redis Error. \nCause: ", cmd.Err().Error())
		}
		return false
	}

	// Cache hit
	cmdByte, err := cmd.Bytes()
	if err != nil {
		log.Println("internal server error.\nCause: ", err.Error())
		return false
	}

	if err := json.Unmarshal(cmdByte, dest); err != nil {
		log.Println("internal server error.\nCause: ", err.Error())
		return false
	}

	log.Println("✅ cache-hit!")
	return true
}

// SetCache marshals and stores data in Redis cache with TTL
func (c *CacheManager) SetCache(ctx context.Context, key string, data interface{}, ttl time.Duration) {
	bt, err := json.Marshal(data)
	if err != nil {
		log.Println("internal server error.\nCause: ", err.Error())
		return
	}

	if err := c.rdb.Set(ctx, key, string(bt), ttl).Err(); err != nil {
		log.Println("redis error.\nCause: ", err.Error())
		return
	}

	log.Printf("Cache set successfully for key: %s", key)
}

// CacheOrFetch is a generic function that implements cache-aside pattern
func (c *CacheManager) CacheOrFetch(
	ctx context.Context,
	key string,
	ttl time.Duration,
	dest interface{},
	fetchFunc func() (interface{}, error),
) error {
	// Try to get from cache first
	if c.GetFromCache(ctx, key, dest) {
		return nil
	}

	// Cache miss - fetch from source
	log.Println("❎ cache-missed!")

	data, err := fetchFunc()
	if err != nil {
		return err
	}

	// Cache the result
	c.SetCache(ctx, key, data, ttl)

	bt, err := json.Marshal(data)
	if err != nil {
		log.Println("internal server error.\nCause: ", err.Error())
		// Still return the data even if caching failed
		return json.Unmarshal(bt, dest)
	}

	return json.Unmarshal(bt, dest)
}
