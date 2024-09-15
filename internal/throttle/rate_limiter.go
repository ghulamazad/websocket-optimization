package throttle

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisClient wraps the redis connection
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client for rate limiting
func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

// AllowConnection checks if a client is allow to send message based on rate limiting.
func (r *RedisClient) AllowConnection(clientID string, limit int, window time.Duration) bool {
	key := "throttle:" + clientID
	var currentCount int64
	var err error

	// Retry logic for redis operations
	for i:=0; i<3; i++{
	// Increament the count for this client
	currentCount, err = r.client.Incr(ctx, key).Result()
	if err==nil{
		break;
	}
	log.Println("Redis INCR failed, retrying:", err)
	time.Sleep(time.Duration(2^i) * time.Second)
}
	if err != nil {
		log.Println("Redis to increament redis key:", err)
		return false
	}

	// Set the expiration for the first time (time window in seconds)
	if currentCount == 1 {
		for i:=0; i<3; i++{
		err := r.client.Expire(ctx, key, window).Err()
			if err == nil {
				break;
			}
			log.Println("Redis EXPIRE failed, retrying:", err)
			time.Sleep(time.Duration(2 ^ i) * time.Second)
		}
		if err != nil {
			log.Println("Failed to set expiration for Redis key:", err)
		}
	}

	// Check if the current count exceeds the limit
	if currentCount > int64(limit) {
		log.Printf("Client %s exceeded rate limit (%d requests)", clientID, currentCount)
		return false
	}

	return true
}
