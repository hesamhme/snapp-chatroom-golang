package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClientInterface interface {
    AddUserToChatroom(chatroom, username string) error
    GetUsersInChatroom(chatroom string) ([]string, error)
    RemoveUserFromChatroom(chatroom, username string) error
    AddChatroom(chatroom string) error
    GetChatrooms() ([]string, error)
    RemoveChatroom(chatroom string) error
}

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

func (r *RedisClient) AddChatroom(chatroom string) error {
    return r.client.SAdd(ctx, "active_chatrooms", chatroom).Err()
}

func (r *RedisClient) GetChatrooms() ([]string, error) {
    return r.client.SMembers(ctx, "active_chatrooms").Result()
}

func (r *RedisClient) RemoveChatroom(chatroom string) error {
    return r.client.SRem(ctx, "active_chatrooms", chatroom).Err()
}

func (r *RedisClient) ClearChatrooms() error {
    keys, err := r.client.Keys(ctx, "chatroom:*").Result()
    if err != nil {
        return err
    }
    for _, key := range keys {
        r.client.Del(ctx, key)
    }
    return nil
}
