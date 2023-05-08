package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

const (
	duration time.Duration = 3600 * time.Second
)

func NewRedisDB(ctx context.Context, cfg Config) (*redis.Client, error) {
	logrus.Printf("%v", cfg)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	status := rdb.Ping(ctx)
	logrus.Print("Connected to Redis : ", status)
	return rdb, nil
}
