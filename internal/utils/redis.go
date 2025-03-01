package utils

import (
	"encoding/json"
	"src/internal/config"
	"time"

	"github.com/go-redis/redis"
)

// ConnectRedis 连接到Redis并返回redis.Client对象
func ConnectRedis() *redis.Client {
	config := config.LoadRedisConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       0,
	})
	return rdb
}

// SetCache 数据缓存
func SetCache(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ConnectRedis().Set(key, jsonData, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string, dest interface{}) error {
	val, err := ConnectRedis().Get(key).Result()
	if err == redis.Nil {
		return nil // 缓存未命中
	} else if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
