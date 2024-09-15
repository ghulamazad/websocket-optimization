package session

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "redis:6379",
})

// StoreSession stores client session information in redis
func StoreSession(clientID string) string {
	sessionID := "session:" + clientID
	err := rdb.Set(ctx, sessionID, time.Now().String(), 0).Err()
	if err != nil {
		log.Println("Error storing session:", err)
	}
	return sessionID
}

// GetSession retrieves session informations from redis
func GetSession(sessionID string) (string, error) {
	val, err := rdb.Get(ctx, sessionID).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// RemoveSession deletes a session from redis when clients disconnectes.
func RemoveSession(sessionID string) {
	err := rdb.Del(ctx, sessionID).Err()
	if err != nil {
		log.Println("Error removing session:", err)
	}
}
