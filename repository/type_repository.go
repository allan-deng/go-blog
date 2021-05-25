package repository

import (
	"errors"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
)

type ITypeRepository interface {
	//建表
	InitTable() error
	//增
	CreateType(*model.Type) (int64, error)
	//删
	DeleteType(int64) error
	//改
	UpdateType(*model.Type) error
	//按id查找
	FindTypeById(int64) (*model.Type, error)
	//按名称查找
	FindTypeByName(string) (*model.Type, error)
	//查找博客数最多的分类
	FindTop(*Page) ([]model.Type, error)
	//查找所有分类
	FindAll() ([]model.Type, error)
}

type TypeRepository struct {
	mysqlDb *gorm.DB
}

//创建TypeRepository
func newTypeRepository(db *gorm.DB) ITypeRepository {
	return &TypeRepository{mysqlDb: db}
}

/*
这里使用单例模式，但是把dao层的init和get分开，调用之前人为显式的完成初始化
*/
var typeRepositoryIns ITypeRepository

func InitTypeRepository(db *gorm.DB) error {
	typeRepositoryIns = newTypeRepository(db)
	return nil
}

func GetTypeRepository() ITypeRepository {
	return typeRepositoryIns
}

func (s *TypeRepository) InitTable() error {
	return s.mysqlDb.AutoMigrate(&model.Type{}).Error
}

func (s *TypeRepository) CreateType(blogType *model.Type) (int64, error) {
	err := s.mysqlDb.Create(blogType).Error
	return blogType.ID, err
}

func (s *TypeRepository) DeleteType(typeId int64) error {
	p := NewPage(1, 10000)
	blogs, _ := GetBlogRepository().FindBlogByTypeId(typeId, &p)
	if len(blogs) > 0 {
		return errors.New("There are  blogs under this type, cant delete type. ")
	}
	return s.mysqlDb.Delete(&model.Type{}, typeId).Error
}

func (s *TypeRepository) UpdateType(blogType *model.Type) error {
	if blogType.ID <= 0 {
		return errors.New("error: cannot update type without id")
	}

	return s.mysqlDb.Save(blogType).Error
}

func (s *TypeRepository) FindTypeById(typeId int64) (*model.Type, error) {
	blogType := &model.Type{}
	return blogType, s.mysqlDb.Preload("Blogs", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).First(blogType, typeId).Error
}

func (s *TypeRepository) FindTypeByName(typeName string) (*model.Type, error) {
	blogType := &model.Type{}
	return blogType, s.mysqlDb.Where("name = ?", typeName).Preload("Blogs", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).First(blogType).Error
}

func (s *TypeRepository) FindTop(page *Page) ([]model.Type, error) {
	count := 0
	s.mysqlDb.Model(&model.Type{}).Count(&count)
	page.Count = count
	limit := page.Size
	offset := ((page.Index - 1) * page.Size)

	var res []model.Type

	rows, err := s.mysqlDb.Raw("select id , name from (select id , name ,count(*) as nums from (select t.id , t.name from type as t left join  blog as b on t.id = b.type_id) as t group by id order by nums DESC) as c limit ? offset ?", limit, offset).Rows()
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		blogType := model.Type{}
		s.mysqlDb.ScanRows(rows, &blogType)
		res = append(res, blogType)
	}

	//关联blog
	for i := 0; i < len(res); i++ {
		var blogs []model.Blog
		err := s.mysqlDb.Omit("content").Where("type_id = ?", res[i].ID).Find(&blogs).Error
		if err != nil {
			return res, err
		}
		res[i].Blogs = blogs
	}
	return res, err
}

func (s *TypeRepository) FindAll() ([]model.Type, error) {
	var res []model.Type
	return res, s.mysqlDb.Preload("Blogs", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Find(&res).Error
}
