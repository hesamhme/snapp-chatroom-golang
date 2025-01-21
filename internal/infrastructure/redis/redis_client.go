package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return &RedisClient{client: client}
}

func (r *RedisClient) AddUserToChatroom(chatroom, username string) error {
	key := fmt.Sprintf("chatroom:%s:users", chatroom)
	return r.client.SAdd(ctx, key, username).Err()
}

func (r *RedisClient) GetUsersInChatroom(chatroom string) ([]string, error) {
	key := fmt.Sprintf("chatroom:%s:users", chatroom)
	return r.client.SMembers(ctx, key).Result()
}

func (r *RedisClient) RemoveUserFromChatroom(chatroom, username string) error {
	key := fmt.Sprintf("chatroom:%s:users", chatroom)
	return r.client.SRem(ctx, key, username).Err()
}
