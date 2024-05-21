package service

import (
	"github.com/go-redis/redis"
	"github.com/xissg/open-api-platform/utils"
)

// redis仅读取和删除，不更新，
type Redis struct {
	rdb *redis.Client
}

func initRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
	return rdb
}
func NewRedis() *Redis {
	rdb := initRedis()
	return &Redis{
		rdb: rdb,
	}
}

func (r *Redis) Get(key string) interface{} {
	return r.rdb.Get(key)
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.rdb.Set(key, value, utils.RandomExpireTime()).Err()
}
