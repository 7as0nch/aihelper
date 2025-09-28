package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"time"
)

var RedisClient *redis.Client
var ctx = context.Background()

// 初始化Redis连接
func InitRedis() {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	db := viper.GetInt("redis.db")
	password := viper.GetString("redis.password")

	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established successfully")
}

// Get 获取Redis中的值
func Get(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Set 设置Redis中的值
func Set(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Del 删除Redis中的键
func Del(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// HashGet 获取哈希表中的字段值
func HashGet(hashKey, field string) (string, error) {
	return RedisClient.HGet(ctx, hashKey, field).Result()
}

// HashSet 设置哈希表中的字段值
func HashSet(hashKey, field string, value interface{}) error {
	return RedisClient.HSet(ctx, hashKey, field, value).Err()
}

// ListPush 将元素添加到列表的头部或尾部
func ListPush(key string, values ...interface{}) error {
	return RedisClient.RPush(ctx, key, values...).Err()
}

// ListPop 从列表的头部或尾部弹出元素
func ListPop(key string) (string, error) {
	return RedisClient.LPop(ctx, key).Result()
}

// ListRange 获取列表指定范围的元素
func ListRange(key string, start, stop int64) ([]string, error) {
	return RedisClient.LRange(ctx, key, start, stop).Result()
}