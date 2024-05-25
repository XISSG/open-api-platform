package service

import (
	"fmt"
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/dal/query"
	"github.com/xissg/open-api-platform/models"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type Mysql struct {
	query *query.Query
}

type DbConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

func initDBConn() *gorm.DB {
	var dbConfig DbConfig

	path := "./conf/mysql.yaml"
	data, err := os.ReadFile(path)
	err = yaml.Unmarshal(data, &dbConfig)
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func NewMysql() *Mysql {
	db := initDBConn()
	q := query.Use(db)
	return &Mysql{
		query: q,
	}
}

// 创建接口信息
func (s *Mysql) CreateInterfaceInfo(values ...*model.InterfaceInfo) error {
	tx := s.query.Begin()
	err := tx.InterfaceInfo.Create(values...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (s *Mysql) GetInterfaceInfoById(id int64) (*model.InterfaceInfo, error) {
	cond := s.query.InterfaceInfo.ID.Eq(id)
	alive := s.query.InterfaceInfo.IsDelete.Eq(Normal)
	result, err := s.query.InterfaceInfo.Where(cond, alive).First()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 分页查询
func (s *Mysql) GetInterfaceInfoList(ctx models.QueryInfoRequest) ([]*model.InterfaceInfo, error) {
	offset := (ctx.Page - 1) * ctx.PageSize
	orderBy := s.query.InterfaceInfo.ID.Desc()
	alive := s.query.InterfaceInfo.IsDelete.Eq(Normal)
	results, err := s.query.InterfaceInfo.Where(alive).Order(orderBy).Limit(ctx.PageSize).Offset(offset).Find()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Mysql) UpdateInterfaceInfo(ctx models.UpdateInfoRequest) error {
	cond := s.query.InterfaceInfo.ID.Eq(ctx.ID)
	alive := s.query.InterfaceInfo.IsDelete.Eq(Normal)
	tx := s.query.Begin()
	_, err := tx.InterfaceInfo.Where(cond, alive).Updates(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_ = tx.Commit()
	return nil
}

const (
	Delete = int32(1)
	Normal = int32(0)
)

func (s *Mysql) DeleteInterfaceInfo(id int64) error {
	cond := s.query.InterfaceInfo.ID.Eq(id)
	alive := s.query.InterfaceInfo.IsDelete.Eq(Normal)
	update := s.query.InterfaceInfo.IsDelete
	tx := s.query.Begin()
	_, err := tx.InterfaceInfo.Where(cond, alive).Update(update, Delete)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}
