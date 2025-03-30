package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func InitRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: Env.RedisHost,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Error connecting to Redis", err)
	}
	log.Println("Connected to Redis")

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func tokenVersionKey(username, siteId string) string {
	return fmt.Sprintf("%s::user::%s", siteId, username)
}

func (s *RedisClient) GetSessionVersion(username, siteId string) int {
	version, err := s.Client.Get(s.Ctx, tokenVersionKey(username, siteId)).Int()
	if err != nil {
		return 0
	}
	return version
}

func (s *RedisClient) IncrementSessionVersion(username, siteId string) {
	s.Client.Incr(s.Ctx, tokenVersionKey(username, siteId))
}
