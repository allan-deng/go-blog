package repository

import (
	"log"
	"strconv"
	"testing"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getBlogRepository() IBlogRepository {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}
	db.SingularTable(true)
	InitBlogRepository(db)
	return GetBlogRepository()
}

func TestBlogRepository_CreateBlog(t *testing.T) {
	repository := getBlogRepository()

	blogtype := &model.Blog{
		Title:        "文章1",
		Content:      "内容内容内容内容内容内容内容",
		FirstPicture: "11",
		Flag:         "22",
		Views:        13,
		Recommend:    true,
		TypeID:       1,
		Tags:         []model.Tag{{ID: 23}},
		UserID:       0,
		Description:  "描述描述描述描述描述描述",
	}
	id, err := repository.CreateBlog(blogtype)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)

}

func TestBlogRepository_DeleteBlog(t *testing.T) {
	repository := getBlogRepository()
	err := repository.DeleteBlog(12)
	if err != nil {
		log.Println(err)
	}
}

func TestBlogRepository_UpdateBlog_FindBlogById(t *testing.T) {
	repository := getBlogRepository()
	blog, err := repository.FindBlogById(11)
	if err != nil {
		log.Println(err)
	}
	log.Println(blog)
	blog.Title = "测试一下"
	blog.TypeID = 1
	newTag := model.Tag{
		Name: "newTag",
	}
	blog.Tags = append(blog.Tags, newTag)
	log.Println(blog)
	repository.UpdateBlog(blog)
}

func TestBlogRepository_FindAll(t *testing.T) {
	repository := getBlogRepository()
	res, err := repository.FindAll(&Page{Index: 2, Size: 2})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}

func TestBlogRepository_FindRecommendTop(t *testing.T) {
	repository := getBlogRepository()
	res, err := repository.FindRecommendTop(&Page{Index: 1, Size: 2})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}

func TestBlogRepository_FindByYear(t *testing.T) {
	repository := getBlogRepository()
	res, err := repository.FindByYear(strconv.Itoa(2021))
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}

func TestBlogRepository_FindGroupYear(t *testing.T) {
	repository := getBlogRepository()
	res, err := repository.FindGroupYear()
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}

func TestBlogRepository_FindByQuery(t *testing.T) {
	repository := getBlogRepository()
	res, err := repository.FindByQuery("内容", &Page{Index: 1, Size: 10})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
