package cron

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/logger"
	"github.com/xissg/open-api-platform/models"
	"github.com/xissg/open-api-platform/service"
	"time"
)

type Fn func()

type Cron struct {
	ticker *time.Ticker //定时器
	runner Fn           //处理函数
}

func (c *Cron) Start() error {
	for {
		select {
		case <-c.ticker.C:
			c.runner()

		}
		defer func() {
			if err := recover(); err != nil {
				logger.SugarLogger.Panic(err)
			}
		}()
	}
}

func (c *Cron) Stop() {
	c.ticker.Stop()
}

const (
	DEFAULT_PAGE      = 1
	DEFAULT_PAGE_SIZE = 100
)

// 默认定时任务
func defaultHandlerFunc() error {
	var interfaceInfo []*model.InterfaceInfo
	var err error

	redisClient := service.NewRedis()
	mysqlClient := service.NewMysql()

	interfaceInfo, err = mysqlClient.GetInterfaceInfoList(models.QueryInfoRequest{
		Page:     DEFAULT_PAGE,
		PageSize: DEFAULT_PAGE_SIZE,
	})

	if err != nil {
		return err
	}

	for i := range interfaceInfo {
		data, _ := json.Marshal(interfaceInfo[i])
		_ = redisClient.ZAdd("interfaceInfo", redis.Z{
			Score:  float64(interfaceInfo[i].ID),
			Member: string(data),
		})
	}

	return nil
}

type Option func(*Cron)

func WithDuration(duration time.Duration) Option {
	return func(c *Cron) {
		c.ticker = time.NewTicker(duration)
	}
}

func WithHandlerFunc(handlerFunc Fn) Option {
	return func(c *Cron) {
		c.runner = handlerFunc
	}
}

func New(options ...Option) *Cron {
	cron := &Cron{}
	for _, option := range options {
		option(cron)
	}
	return cron
}
