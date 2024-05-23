package service

import (
	"github.com/xissg/open-api-platform/dal/model"
	"github.com/xissg/open-api-platform/models"
)

// 创建接口信息
func (s *Mysql) CreateUser(values ...*model.User) error {
	tx := s.query.Begin()
	err := tx.User.Create(values...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (s *Mysql) GetUserByName(name string) (*model.User, error) {
	cond := s.query.User.UserName.Eq(name)
	alive := s.query.User.IsDelete.Eq(Normal)
	result, err := s.query.User.Where(cond, alive).First()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Mysql) GetUserById(id int64) (*model.User, error) {
	cond := s.query.User.ID.Eq(id)
	alive := s.query.User.IsDelete.Eq(Normal)
	result, err := s.query.User.Where(cond, alive).First()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Mysql) GetUserByAccessKey(accessKey string) (*model.User, error) {
	cond := s.query.User.AccessKey.Eq(accessKey)
	alive := s.query.User.IsDelete.Eq(Normal)
	result, err := s.query.User.Where(cond, alive).First()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 分页查询
func (s *Mysql) GetUserList(ctx models.QueryUserRequest) ([]*model.User, error) {
	offset := (ctx.Page - 1) * ctx.PageSize
	orderBy := s.query.User.ID
	alive := s.query.User.IsDelete.Eq(Normal)
	results, err := s.query.User.Where(alive).Order(orderBy).Limit(ctx.PageSize).Offset(offset).Find()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Mysql) UpdateUser(ctx models.UpdateUserRequest) error {
	cond := s.query.User.ID.Eq(ctx.ID)
	alive := s.query.User.IsDelete.Eq(Normal)
	tx := s.query.Begin()
	_, err := tx.User.Where(cond, alive).Updates(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_ = tx.Commit()
	return nil
}

func (s *Mysql) DeleteUser(id int64) error {
	cond := s.query.User.ID.Eq(id)
	alive := s.query.User.IsDelete.Eq(Normal)
	update := s.query.User.IsDelete
	tx := s.query.Begin()
	_, err := tx.User.Where(cond, alive).Update(update, Delete)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (s *Mysql) IsUserExist(name string) bool {
	user, err := s.GetUserByName(name)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}
