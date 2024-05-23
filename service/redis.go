package service

import (
	"github.com/go-redis/redis"
	"github.com/xissg/open-api-platform/utils"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"time"
)

// redis仅读取和删除，不更新，
type Redis struct {
	rdb *redis.Client
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func initRedis() *redis.Client {
	path, _ := os.Getwd()
	path = filepath.Dir(path)
	path = filepath.Join(path, "conf", "redis.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg RedisConfig
	err = yaml.Unmarshal(data, &cfg)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // 没有密码，默认值
		DB:       cfg.DB,       // 默认DB 0
	})
	return rdb
}
func NewRedis() *Redis {
	rdb := initRedis()
	return &Redis{
		rdb: rdb,
	}
}

func (r *Redis) Get(key string) (interface{}, error) {
	res := r.rdb.Get(key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return res, nil
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.rdb.Set(key, value, utils.RandomExpireTime()).Err()
}

func (r *Redis) Delete(key string) error {
	return r.rdb.Del(key).Err()
}

func (r *Redis) Exists(key string) bool {
	exists, _ := r.rdb.Exists(key).Result()
	if exists != 1 {
		return false
	}
	return true
}

func (r *Redis) Expire(key string, expireTime time.Duration) error {
	return r.rdb.Expire(key, expireTime).Err()
}

func (r *Redis) Incr(key string) {
	r.rdb.Incr(key)
}
