package config

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func InitRedisClient() *RedisClient {
	// Support both formats:
	// 1) Env.RedisHost = "redis://:pwd@host:6379/0" or "rediss://..." (recommended)
	// 2) Env.RedisHost = "host:port" + Env.RedisPwd
	var opts *redis.Options

	if strings.HasPrefix(Env.RedisHost, "redis://") || strings.HasPrefix(Env.RedisHost, "rediss://") {
		parsed, err := redis.ParseURL(Env.RedisHost)
		if err != nil {
			log.Printf("Error parsing Redis URL: %v", err)
		} else {
			opts = parsed
		}
	}

	if opts == nil {
		// Backward-compatible: Addr + Password
		opts = &redis.Options{
			Addr:         Env.RedisHost,
			Password:     Env.RedisPwd,
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 2,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		}
	}

	// If password is provided separately, prefer it (useful when Env.RedisHost is host:port).
	if Env.RedisPwd != "" {
		opts.Password = Env.RedisPwd
	}

	// IMPORTANT:
	// Don't force TLS for every Redis endpoint.
	// - If your Redis supports TLS, use a `rediss://` URL in Env.RedisHost so ParseURL enables TLS.
	// - If your Redis is plain TCP (most local/dev setups), use `redis://` or host:port.

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Printf("Error connecting to Redis: %v", err)
	} else {
		log.Println("✅ Connected to Redis")
	}

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func tokenVersionKey(username, siteId string) string {
	return fmt.Sprintf("%s::user::%s", siteId, username)
}

func (s *RedisClient) GetTokenVersion(username, siteId string) int {
	version, err := s.Client.Get(s.Ctx, tokenVersionKey(username, siteId)).Int()
	if err != nil {
		return 0
	}
	return version
}

func (s *RedisClient) IncrementTokenVersion(username, siteId string) {
	s.Client.Incr(s.Ctx, tokenVersionKey(username, siteId))
}

func sessionBlackListKey(sessionId string) string {
	return fmt.Sprintf("session::blacklist::%s", sessionId)
}

func (s *RedisClient) AddSessionIdToBlackList(sessionId string) {
	fmt.Println("black list", sessionBlackListKey(sessionId))
	s.Client.Set(s.Ctx, sessionBlackListKey(sessionId), 1, AccessTokenExpire)
}

func (s *RedisClient) ValidateSessionId(sessionId string) bool {
	existed, _ := s.Client.Get(s.Ctx, sessionBlackListKey(sessionId)).Int()
	return existed == 0
}
