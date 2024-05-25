package cron

import (
	"encoding/json"
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/models"
	"github.com/xissg/open-api-platform/service"
	"strconv"
	"time"
)

type Fn func() error

type Cron struct {
	ticker *time.Ticker //定时器
	runner Fn           //处理函数
}

func NewCron(timeDuration time.Duration, fn Fn) *Cron {
	//设置默认处理函数
	if fn == nil {
		fn = defaultHandlerFunc
	}

	return &Cron{
		ticker: time.NewTicker(timeDuration),
		runner: fn,
	}
}

func (c *Cron) Start() error {
	for {
		select {
		case <-c.ticker.C:
			if err := c.runner(); err != nil {
				return err
			}
		}
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
		_ = redisClient.Set(strconv.FormatInt(interfaceInfo[i].ID, 10), data)
	}

	return nil
}

func StartCron() {
	go func() {
		cron := NewCron(time.Hour*24, defaultHandlerFunc)
		err := cron.Start()
		if err != nil {
			return
		}
		defer cron.Stop()
	}()
}
