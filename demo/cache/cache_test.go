package cache_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/miajio/dpsk/pkg/redis"
)

func TestCache(t *testing.T) {
	rc := redis.RedisConfig{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "",
		PoolSize: 20,
	}
	redisClient, err := rc.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer redisClient.Close()
	// 测试连接
	ctx := context.Background()

	// 测试基础操作
	if err := redisClient.Set(ctx, "test_key", "test_value"); err != nil {
		log.Printf("Set failed: %v", err)
	}

	if err := redisClient.SetEx(ctx, "test_key_ex", "test_value", 10*time.Second); err != nil {
		log.Printf("Set failed: %v", err)
	}

	val, err := redisClient.Get(ctx, "test_key")
	if err != nil {
		log.Printf("Get failed: %v", err)
	} else {
		fmt.Printf("Got value: %s\n", val)
	}

	val, err = redisClient.Get(ctx, "test_key_ex")
	if err != nil {
		log.Printf("Get failed: %v", err)
	} else {
		fmt.Printf("Got ex value: %s\n", val)
	}

	// 测试分布式锁
	lockOptions := redis.LockOptions{
		Expiration: 5 * time.Second,
		Timeout:    10 * time.Second,
		RetryDelay: 100 * time.Millisecond,
	}

	// 尝试获取锁
	err = redisClient.TryLock(ctx, "h", "userabc", lockOptions)
	if err != nil {
		log.Printf("TryLock failed: %v", err)
	} else {
		log.Println("TryLock success")
		defer func() {
			if err := redisClient.ReleaseLock(ctx, "h", "userabc"); err != nil {
				log.Printf("ReleaseLock failed: %v", err)
			} else {
				log.Println("Lock released")
			}
		}()
	}

	// 阻塞式获取锁
	err = redisClient.AcquireLock(ctx, "h", "userabc", lockOptions)
	if err != nil {
		log.Fatalf("AcquireLock failed: %v", err)
	}
	log.Println("AcquireLock success")

	// 锁续期
	if err := redisClient.RenewLock(ctx, "h", "userabc", 10*time.Second); err != nil {
		log.Printf("RenewLock failed: %v", err)
	} else {
		log.Println("Lock renewed")
	}

	// 释放锁
	if err := redisClient.ReleaseLock(ctx, "h", "userabc"); err != nil {
		log.Printf("ReleaseLock failed: %v", err)
	} else {
		log.Println("Lock released")
	}

}
