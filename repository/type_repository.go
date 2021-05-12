package repository

import (
	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
)

/*
CURD
TODO:
按名称查找分类
Tag findByName(String name);

查找博客数量最多的分类
@Query("select t from Tag t")
List<Tag> findTop(Pageable pageable);
*/

type ITypeRepository interface {
	InitTable() error
	CreateType(*model.Type) (int64, error)
	DeleteType(int64) error
	UpdateType(*model.Type, ...string) error
	FindTypeById(int64) (*model.Type, error)
	FindTypeByName(string) (*model.Type, error)
	FindTop(Page) ([]model.Type, error)
	FindAll() ([]model.Type, error)
}

type TypeRepository struct {
	mysqlDb *gorm.DB
}

//创建typeRepository
func NewTypeRepository(db *gorm.DB) ITypeRepository {
	return &TypeRepository{mysqlDb: db}
}

func (s *TypeRepository) InitTable() error {
	return s.mysqlDb.AutoMigrate(&model.Type{}).Error
}

func (s *TypeRepository) CreateType(blogType *model.Type) (int64, error) {
	err := s.mysqlDb.Create(blogType).Error
	return blogType.ID, err
}

func (s *TypeRepository) DeleteType(typeId int64) error {
	return s.mysqlDb.Delete(&model.Type{}, typeId).Error
}

func (s *TypeRepository) UpdateType(blogType *model.Type, columns ...string) error {

}

func (s *TypeRepository) FindTypeById(_ int64) (*model.Type, error) {
	panic("not implemented") // TODO: Implement
}

func (s *TypeRepository) FindTypeByName(_ string) (*model.Type, error) {
	panic("not implemented") // TODO: Implement
}

func (s *TypeRepository) FindTop(_ repository.Page) ([]model.Type, error) {
	panic("not implemented") // TODO: Implement
}

func (s *TypeRepository) FindAll() ([]model.Type, error) {
	panic("not implemented") // TODO: Implement
}
