package repository

import (
	"errors"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
)

type ITagRepository interface {
	//建表
	InitTable() error
	//增
	CreateTag(*model.Tag) (int64, error)
	//删
	DeleteTag(int64) error
	//改
	UpdateTag(*model.Tag) error
	//按id查找
	FindTagById(int64) (*model.Tag, error)
	//按名称查找
	FindTagByName(string) (*model.Tag, error)
	//查找博客数最多的标签
	FindTop(*Page) ([]model.Tag, error)
	//查找所有标签
	FindAll() ([]model.Tag, error)
}

type TagRepository struct {
	mysqlDb *gorm.DB
}

//创建blogRepository
func newTagRepository(db *gorm.DB) ITagRepository {
	return &TagRepository{mysqlDb: db}
}

/*
这里使用单例模式，但是把dao层的init和get分开，调用之前人为显式的完成初始化
*/
var tagRepositoryIns ITagRepository

func InitTagRepository(db *gorm.DB) error {
	tagRepositoryIns = newTagRepository(db)
	return nil
}

func GetTagRepository() ITagRepository {
	return tagRepositoryIns
}

func (s *TagRepository) InitTable() error {
	return s.mysqlDb.AutoMigrate(&model.Tag{}).Error
}

func (s *TagRepository) CreateTag(tag *model.Tag) (int64, error) {
	err := s.mysqlDb.Create(tag).Error
	return tag.ID, err
}

func (s *TagRepository) DeleteTag(tagId int64) error {
	// 先删除连结表
	err := s.mysqlDb.Raw("delete from blog_tag where tag_id = ?", tagId).Error
	if err != nil {
		return err
	}
	return s.mysqlDb.Delete(&model.Tag{}, tagId).Error
}

func (s *TagRepository) UpdateTag(tag *model.Tag) error {
	if tag.ID <= 0 {
		return errors.New("error: cannot update tag without id")
	}

	return s.mysqlDb.Save(tag).Error
}

func (s *TagRepository) FindTagById(tagId int64) (*model.Tag, error) {
	tag := &model.Tag{}
	return tag, s.mysqlDb.First(tag, tagId).Error
}

func (s *TagRepository) FindTagByName(tagName string) (*model.Tag, error) {
	tag := &model.Tag{}
	return tag, s.mysqlDb.Where("name = ?", tagName).Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).First(tag).Error
}

func (s *TagRepository) FindTop(page *Page) ([]model.Tag, error) {
	count := 0
	s.mysqlDb.Model(&model.Tag{}).Count(&count)
	page.Count = count
	limit := page.Size
	offset := ((page.Index - 1) * page.Size)

	var res []model.Tag

	rows, err := s.mysqlDb.Raw("select id , name from (select id , name ,count(*) as nums from (select t.id , t.name from tag as t left join  blog_tag as b on t.id = b.tag_id) as t group by id order by nums DESC) as c limit ? offset ?", limit, offset).Rows()
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		tag := model.Tag{}
		s.mysqlDb.ScanRows(rows, &tag)
		res = append(res, tag)
	}

	//关联blog
	for i := 0; i < len(res); i++ {
		var blogs []model.Blog
		//先获得tag_id 对应的 blog_id
		var blogidlist []int
		rows, err := s.mysqlDb.Raw("select blog_id from blog_tag where tag_id = ?", res[i].ID).Rows()
		if err != nil {
			return res, err
		}
		defer rows.Close()
		for rows.Next() {
			temp := 0
			rows.Scan(&temp)
			blogidlist = append(blogidlist, temp)
		}
		//获取对应id的blog
		err = s.mysqlDb.Omit("content").Find(&blogs, blogidlist).Error
		if err != nil {
			return res, err
		}
		res[i].Blogs = blogs
	}
	return res, err

}

func (s *TagRepository) FindAll() ([]model.Tag, error) {
	var res []model.Tag
	// return res, s.mysqlDb.Preload("Tags", func(db *gorm.DB) *gorm.DB {
	// 	return db.Omit("content")
	// }).Find(&res).Error
	return res, s.mysqlDb.Model(&model.Tag{}).Find(&res).Error
}
