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

func (r *Redis) Set(key string, value interface{}) error {
	return r.rdb.Set(key, value, utils.RandomExpireTime()).Err()
}

func (r *Redis) Get(key string) (string, error) {
	return r.rdb.Get(key).Result()
}

func (r *Redis) Delete(key string) error {
	return r.rdb.Del(key).Err()
}

func (r *Redis) ZAdd(key string, z ...redis.Z) error {
	return r.rdb.ZAdd(key, z...).Err()
}

func (r *Redis) ZRange(key string, page, pageSize int64) ([]string, error) {
	start := (page - 1) * pageSize
	stop := page * pageSize
	res, err := r.rdb.ZRange(key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Redis) ZRem(key string, z ...interface{}) error {
	return r.rdb.ZRem(key, z...).Err()
}

func (r *Redis) HIncrBy(key string, field string, incr int64) {
	r.rdb.HIncrBy(key, field, incr)
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	res, err := r.rdb.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}
