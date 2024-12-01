package cache

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrRedisCredentialsNotFound = errors.New("redis credentials not found")
	Client                      *redis.Client
	retry                       int
)

type Config struct {
	address string
}

func GetConfig() (*Config, error) {

	cfg := &Config{}

	redisAddressEnv := os.Getenv("REDIS_ADDRESS")
	if redisAddressEnv != "" {
		cfg.address = redisAddressEnv
	} else {
		return nil, ErrRedisCredentialsNotFound
	}

	return cfg, nil
}

func Start(cfg *Config) error {

	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := Client.Ping(context.Background()).Result()
	times := retry
	for err != nil && times > 0 {
		_, err = Client.Ping(context.Background()).Result()
		times--
		time.After(time.Second)
	}
	if err != nil {
		return err
	}

	slog.Info("redis started...")

	return nil
}

func Close() error {
	slog.Info("redis stopped...")
	return Client.Close()
}
