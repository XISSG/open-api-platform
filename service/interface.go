package service

import (
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/dal/query"
	"github.com/xissg/open-api-platform/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	query *query.Query
}

func initDBConn() (*gorm.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

func NewService() *Service {
	db, err := initDBConn()
	if err != nil {
		panic(err)
	}
	q := query.Use(db)
	return &Service{
		query: q,
	}
}

// 创建接口信息
func (s *Service) CreateInterfaceInfo(values ...*model.InterfaceInfo) error {
	tx := s.query.Begin()
	err := tx.InterfaceInfo.Create(values...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (s *Service) GetInterfaceInfoById(id int64) (*model.InterfaceInfo, error) {
	cond := s.query.InterfaceInfo.ID.Eq(id)
	alive := s.query.InterfaceInfo.IsDelete.Eq(Normal)
	result, err := s.query.InterfaceInfo.Where(cond, alive).First()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 分页查询
func (s *Service) GetInterfaceInfoList(ctx models.QueryInfoRequest) ([]*model.InterfaceInfo, error) {
	offset := (ctx.Page - 1) * ctx.PageSize
	orderBy := s.query.InterfaceInfo.ID
	alive := s.query.InterfaceInfo.IsDelete.Eq(Normal)
	results, err := s.query.InterfaceInfo.Where(alive).Order(orderBy).Limit(ctx.PageSize).Offset(offset).Find()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Service) UpdateInterfaceInfo(ctx models.UpdateInfoRequest) error {
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

func (s *Service) DeleteInterfaceInfo(id int64) error {
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
