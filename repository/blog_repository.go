package repository

import (
	"errors"
	"time"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
)

/*
TODO:

查找年份
查找某年的文章
*/

type IBlogRepository interface {
	//建表
	InitTable() error
	//增
	CreateBlog(*model.Blog) (int64, error)
	//删
	DeleteBlog(int64) error
	//改
	UpdateBlog(*model.Blog) error
	//按id查找
	FindBlogById(int64) (*model.Blog, error)
	//按typeid查找
	FindBlogByTypeId(int64, *Page) ([]model.Blog, error)
	//按tagid查找
	FindBlogByTagId(int64, *Page) ([]model.Blog, error)
	//分页查找全部，按createtime排序
	FindAll(*Page) ([]model.Blog, error)
	//模糊查找，按createtime排序
	FindByQuery(string, *Page) ([]model.Blog, error)
	//已推荐的文章，按updatetime排序
	FindRecommendTop(*Page) ([]model.Blog, error)
	//查找所有发布年份
	FindGroupYear() ([]string, error)
	//安年份查找，按createtime排序
	FindByYear(string) ([]model.Blog, error)
	//博客计数
	Count() (int64, error)
	//按标题、type、recommend查找
	FindBlogByTitleAndTypeIdAndRecommend(title string, typeid int64, recommend bool) ([]model.Blog, error)
}

type BlogRepository struct {
	mysqlDb *gorm.DB
}

//创建blogRepository
func newBlogRepository(db *gorm.DB) IBlogRepository {
	return &BlogRepository{mysqlDb: db}
}

/*
这里使用单例模式，但是把dao层的init和get分开，调用之前人为显式的完成初始化
*/
var blogRepositoryIns IBlogRepository

func InitBlogRepository(db *gorm.DB) error {
	blogRepositoryIns = newBlogRepository(db)
	return nil
}

func GetBlogRepository() IBlogRepository {
	return blogRepositoryIns
}

func (s *BlogRepository) InitTable() error {
	return s.mysqlDb.AutoMigrate(&model.Blog{}).Error
}

func (s *BlogRepository) CreateBlog(blog *model.Blog) (int64, error) {
	blog.CreateTime = time.Now()
	blog.UpdateTime = time.Now()
	err := s.mysqlDb.Create(blog).Error
	//多对多的连接由gorm完成不需要自行添加连接
	// tags := blog.Tags
	// for _, tag := range tags {
	// 	if len(tag.Name) != 0 {
	// 		//tag的名称不为空时，创建连结
	// 		tag := &model.Tag{
	// 			Name: tag.Name,
	// 		}
	// 		if s.mysqlDb.Debug().Where("name = ?", tag.Name).First(tag).RecordNotFound() {
	// 			//如果不存在则插入
	// 			err := s.mysqlDb.Debug().Create(tag).Error
	// 			return 0, err
	// 		}
	// 		s.mysqlDb.Debug().Raw("INSERT INTO blog_tag ( blog_id , tag_id ) VALUES ( ?, ? )", blog.ID, tag.ID)
	// 	}
	// }
	return blog.ID, err
}

func (s *BlogRepository) DeleteBlog(blogId int64) error {
	// 先删除连结表
	err := s.mysqlDb.Debug().Exec("delete from blog_tag where blog_id = ?", blogId).Error
	if err != nil {
		return err
	}
	return s.mysqlDb.Delete(&model.Blog{}, blogId).Error
}

func (s *BlogRepository) UpdateBlog(blog *model.Blog) error {
	blog.UpdateTime = time.Now()
	if blog.ID <= 0 {
		return errors.New("error: cannot update type without id")
	}
	// tags := blog.Tags
	// for _, tag := range tags {
	// 	if len(tag.Name) != 0 {
	// 		//tag的名称不为空时，创建连结
	// 		tag := &model.Tag{}
	// 		if s.mysqlDb.Where("name = ?", tag.Name).First(tag).RecordNotFound() {
	// 			//如果不存在则插入
	// 			err := s.mysqlDb.Create(tag).Error
	// 			return err
	// 		}
	// 		if s.mysqlDb.Raw("select * from blog_tag where blog_id = ? and tag_id = ?", blog.ID, tag.ID).RecordNotFound() {
	// 			//连结不存在时创建连结
	// 			s.mysqlDb.Raw("INSERT INTO blog_tag ( blog_id , tag_id ) VALUES ( ?, ? )", blog.ID, tag.ID)
	// 		}
	// 	}
	// }
	//save操作会自动添加连接
	return s.mysqlDb.Omit("CreateTime").Save(blog).Error
}

func (s *BlogRepository) FindBlogById(blogId int64) (*model.Blog, error) {
	blog := &model.Blog{}
	return blog, s.mysqlDb.Preload("User").Preload("Type").Preload("Tags").Preload("Comments").First(blog, blogId).Error
}

func (s *BlogRepository) FindBlogByTypeId(typeId int64, page *Page) ([]model.Blog, error) {
	count := 0
	s.mysqlDb.Model(&model.Blog{}).Where("type_id = ?", typeId).Count(&count)
	page.Count = count
	limit := page.Size
	offset := ((page.Index - 1) * page.Size)

	var res []model.Blog
	err := s.mysqlDb.Offset(offset).Limit(limit).Preload("User").Preload("Type").Preload("Comments").Preload("Tags").Order("create_time DESC").Where("type_id = ?", typeId).Find(&res).Error
	return res, err
	// return blog, s.mysqlDb.Preload("User").Preload("Type").Preload("Tags").Preload("Comments").First(blog, blogId).Error
}

func (s *BlogRepository) FindBlogByTitleAndTypeIdAndRecommend(title string, typeid int64, recommend bool) ([]model.Blog, error) {
	var res []model.Blog
	rec := 0
	if recommend {
		rec = 1
	}
	title = "%" + title + "%"
	err := s.mysqlDb.Preload("User").Preload("Type").Preload("Comments").Preload("Tags").Order("create_time DESC").Where("type_id = ?", typeid).Where("recommend = ?", rec).Where("title like ?", title).Find(&res).Error
	return res, err
}

func (s *BlogRepository) FindBlogByTagId(tagId int64, page *Page) ([]model.Blog, error) {
	count := 0
	row := s.mysqlDb.Raw("select count(*) from blog_tag where tag_id = ?", tagId).Row()
	row.Scan(&count)
	page.Count = count
	limit := page.Size
	offset := ((page.Index - 1) * page.Size)

	var res []model.Blog
	rows, err := s.mysqlDb.Raw("select * from blog where id in (select blog_id from blog_tag where tag_id = ?) ORDER BY create_time DESC limit ? offset ?", tagId, limit, offset).Rows()
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp model.Blog
		s.mysqlDb.ScanRows(rows, &temp)
		user, err := GetUserRepository().FindUserById(temp.UserID)

		blogtype, err := GetTypeRepository().FindTypeById(temp.TypeID)

		comments, err := GetCommentRepository().FindAllCommentByBlogId(temp.ID)

		temp.User = *user
		temp.Type = *blogtype
		temp.Comments = comments
		res = append(res, temp)
		if err != nil {

		}

	}
	return res, err
	// return blog, s.mysqlDb.Preload("User").Preload("Type").Preload("Tags").Preload("Comments").First(blog, blogId).Error
}

func (s *BlogRepository) FindAll(page *Page) ([]model.Blog, error) {
	count := 0
	s.mysqlDb.Model(&model.Blog{}).Count(&count)
	page.Count = count
	limit := page.Size
	offset := ((page.Index - 1) * page.Size)

	var res []model.Blog

	err := s.mysqlDb.Offset(offset).Limit(limit).Preload("User").Preload("Type").Preload("Comments").Preload("Tags").Order("create_time DESC").Find(&res).Error

	return res, err
}

func (s *BlogRepository) FindRecommendTop(page *Page) ([]model.Blog, error) {
	count := 0
	s.mysqlDb.Model(&model.Blog{}).Count(&count)
	page.Count = count
	limit := page.Size
	offset := ((page.Index - 1) * page.Size)

	var res []model.Blog

	err := s.mysqlDb.Offset(offset).Limit(limit).Preload("User").Preload("Type").Preload("Comments").Preload("Tags").Order("update_time DESC").Where("recommend = 1").Find(&res).Error

	return res, err
}

func (s *BlogRepository) FindByQuery(query string, page *Page) ([]model.Blog, error) {
	count := 0
	s.mysqlDb.Model(&model.Blog{}).Count(&count)
	page.Count = count
	limit := page.Size
	offset := ((page.Index - 1) * page.Size)
	query = "%" + query + "%"
	var res []model.Blog

	err := s.mysqlDb.Offset(offset).Limit(limit).Preload("User").Preload("Type").Preload("Comments").Preload("Tags").Order("create_time DESC").Where("title like ? or content like ?", query, query).Find(&res).Error

	return res, err
}

func (s *BlogRepository) FindGroupYear() ([]string, error) {
	var res []string
	rows, err := s.mysqlDb.Raw("SELECT DATE_FORMAT(b.create_time,'%Y') as year FROM blog as b  GROUP BY DATE_FORMAT(b.create_time,'%Y') order by year desc").Rows()
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp string
		rows.Scan(&temp)
		res = append(res, temp)
	}
	return res, err
}

func (s *BlogRepository) FindByYear(year string) ([]model.Blog, error) {
	var res []model.Blog
	rows, err := s.mysqlDb.Raw("select * from blog as b where DATE_FORMAT(b.create_time,'%Y') = ? order by create_time desc", year).Rows()
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var blog model.Blog
		s.mysqlDb.ScanRows(rows, &blog)
		res = append(res, blog)
	}
	return res, err
}

func (s *BlogRepository) Count() (int64, error) {
	var count int64
	return count, s.mysqlDb.Model(&model.Blog{}).Count(&count).Error
}
