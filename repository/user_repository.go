package repository

import (
	"errors"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
)

type IUserRepository interface {
	//建表
	InitTable() error
	//增
	CreateUser(*model.User) (int64, error)
	//删
	DeleteUser(int64) error
	//改
	UpdateUser(*model.User) error
	//按用户名和密码查找
	FindUserByNameAndPassword(string, string) (*model.User, error)
	//按id查找
	FindUserById(int64) (*model.User, error)
}

type UserRepository struct {
	mysqlDb *gorm.DB
}

//创建userRepository
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{mysqlDb: db}
}

func (s *UserRepository) InitTable() error {
	return s.mysqlDb.AutoMigrate(&model.User{}).Error
}

func (s *UserRepository) CreateUser(user *model.User) (int64, error) {
	err := s.mysqlDb.Create(user).Error
	return user.ID, err
}

func (s *UserRepository) DeleteUser(userId int64) error {
	return s.mysqlDb.Delete(&model.User{}, userId).Error
}

func (s *UserRepository) UpdateUser(user *model.User) error {
	if user.ID <= 0 {
		return errors.New("error: cannot update user without id")
	}
	return s.mysqlDb.Save(user).Error
}

func (s *UserRepository) FindUserById(userId int64) (*model.User, error) {
	user := &model.User{}
	return user, s.mysqlDb.First(user, userId).Error
}

func (s *UserRepository) FindUserByNameAndPassword(username, password string) (*model.User, error) {
	user := &model.User{}
	return user, s.mysqlDb.Where("username = ?", username).Where("password = ?", password).First(user).Error
}
