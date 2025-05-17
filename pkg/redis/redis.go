package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient redis客户端
type RedisClient struct {
	client *redis.Client
}

// Redis 配置
type RedisConfig struct {
	Addr     string `toml:"addr" json:"addr"`           // redis地址
	Password string `toml:"password" json:"password"`   // redis密码
	DB       int    `toml:"db" json:"db"`               // redis数据库
	PoolSize int    `toml:"pool_size" json:"pool_size"` // redis连接池大小
}

func (c *RedisConfig) Init() (*RedisClient, error) {
	client := &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Password: c.Password,
			DB:       c.DB,
			PoolSize: c.PoolSize,
		}),
	}
	// 设置500毫秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

// Ping 测试Redis连接
func (rc *RedisClient) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
}

// Close 关闭Redis连接
func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

// Set 设置缓存
func (rc *RedisClient) Set(ctx context.Context, key string, value any) error {
	return rc.client.Set(ctx, key, value, 0).Err()
}

// SetEx 设置带过期时间的缓存
func (rc *RedisClient) SetEx(ctx context.Context, key string, value any, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}

// GetWithTTL 获取缓存并返回过期时间
func (rc *RedisClient) GetWithTTL(ctx context.Context, key string) (string, time.Duration, error) {
	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		return "", 0, err
	}
	ttl, err := rc.client.TTL(ctx, key).Result()
	if err != nil {
		return "", 0, err
	}
	return result, ttl, nil
}

// Del 删除缓存
func (rc *RedisClient) Del(ctx context.Context, keys ...string) error {
	return rc.client.Del(ctx, keys...).Err()
}

// LockOptions 锁选项
type LockOptions struct {
	Expiration time.Duration // 锁过期时间
	Timeout    time.Duration // 锁等待时间 (仅用于AcquireLock)
	RetryDelay time.Duration // 重试间隔时间 (仅用于AcquireLock)
}

// TryLock 尝试获取锁
// 如果获取锁失败，则返回错误
func (rc *RedisClient) TryLock(ctx context.Context, key, token string, options LockOptions) error {
	success, err := rc.client.SetNX(ctx, key, token, options.Expiration).Result()
	if err != nil {
		return err
	}
	if !success {
		return errors.New("lock already held by another client")
	}
	return nil
}

// AcquireLock 获取锁(阻塞式, 直到获取成功或超时)
// 如果获取锁失败，则返回错误
func (rc *RedisClient) AcquireLock(ctx context.Context, key, token string, options LockOptions) error {
	if options.RetryDelay == 0 {
		options.RetryDelay = 100 * time.Millisecond
	}
	start := time.Now()
	for {
		// 尝试获取锁
		success, err := rc.client.SetNX(ctx, key, token, options.Expiration).Result()
		if err != nil {
			return err
		}
		if success {
			return nil
		}
		// 检查是否超时
		if time.Since(start) > options.Timeout {
			return errors.New("acquire lock timeout")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(options.RetryDelay):
		}
	}
}

// ReleaseLock 释放分布式锁
// token必须与获取锁时使用的token相同
func (rc *RedisClient) ReleaseLock(ctx context.Context, key, token string) error {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	result, err := rc.client.Eval(ctx, script, []string{key}, token).Result()
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return errors.New("lock token mismatch or lock already released")
	}
	return nil
}

// RenewLock 续期分布式锁
// 延长锁的过期时间
// token必须与获取锁时使用的token相同
func (rc *RedisClient) RenewLock(ctx context.Context, key, token string, expiration time.Duration) error {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("pexpire", KEYS[1], ARGV[2])
	else
		return 0
	end
	`

	expirationMs := expiration.Milliseconds()
	result, err := rc.client.Eval(ctx, script, []string{key}, token, expirationMs).Result()
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return errors.New("lock token mismatch or lock does not exist")
	}

	return nil
}
