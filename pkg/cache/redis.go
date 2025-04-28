package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig redis配置
type RedisConfig struct {
	Addr     string `toml:"addr" json:"addr"`         // redis地址
	Password string `toml:"password" json:"password"` // redis密码
	DB       int    `toml:"db" json:"db"`             // redis数据库
	PoolSize int    `toml:"poolSize" json:"poolSize"` // redis连接池大小
}

// RedisClient redis客户端
type RedisClient struct {
	Client *redis.Client
}

// Generator 生成redis客户端
func (redisConfig *RedisConfig) Generator() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &RedisClient{Client: client}, nil
}

// Set 设置key-value
func (r *RedisClient) Set(key string, value interface{}) error {
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, 0).Err()
}

// SetWithExpire 设置key-value并设置过期时间
func (r *RedisClient) SetWithExpire(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取key对应的value
func (r *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	return r.Client.Get(ctx, key).Result()
}

// Exists 判断key是否存在
func (r *RedisClient) Exists(key string) (bool, error) {
	ctx := context.Background()
	result, err := r.Client.Exists(ctx, key).Result()
	return result > 0, err
}

// Delete 删除key
func (r *RedisClient) Delete(key string) error {
	ctx := context.Background()
	return r.Client.Del(ctx, key).Err()
}

// AtomicIncr 原子递增
func (r *RedisClient) AtomicIncr(key string) (int64, error) {
	ctx := context.Background()
	return r.Client.Incr(ctx, key).Result()
}

// AtomicDecr 原子递减
func (r *RedisClient) AtomicDecr(key string) (int64, error) {
	ctx := context.Background()
	return r.Client.Decr(ctx, key).Result()
}
