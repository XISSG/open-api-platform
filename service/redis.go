package service

import (
	"github.com/go-redis/redis"
	"github.com/xissg/open-api-platform/utils"
	"gopkg.in/yaml.v2"
	"os"
)

type Redis struct {
	rdb *redis.Client
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func initRedis() *redis.Client {
	path := "./conf/redis.yaml"
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
