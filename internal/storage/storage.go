package storage

import (
	"context"
	"fmt"
	"time"

	"kube/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	config *config.Config
}

type RedisClient struct {
	client *redis.Client
}

func Init(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
	}
}

func (c *Client) UploadFile(filePath string, data []byte) error {
	// TODO: Implement file upload logic
	return nil
}

func (c *Client) DownloadFile(filePath string) ([]byte, error) {
	// TODO: Implement file download logic
	return nil, nil
}

func (c *Client) DeleteFile(filePath string) error {
	// TODO: Implement file deletion logic
	return nil
}

// RedisClient methods
func NewRedisClient(cfg config.RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisClient{
		client: rdb,
	}
}

func (r *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(key string, value string) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisClient) SetWithExpiry(key string, value string, expiry time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiry).Err()
}

func (r *RedisClient) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Ping() error {
	ctx := context.Background()
	return r.client.Ping(ctx).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
