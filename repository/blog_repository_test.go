package repository

import (
	"log"
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
	return NewBlogRepository(db)
}

func TestBlogRepository_CreateBlog(t *testing.T) {
	repository := getBlogRepository()
	blogtype := &model.Blog{}
	id, err := repository.CreateBlog(blogtype)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)

}

func TestBlogRepository_DeleteBlog(t *testing.T) {
	repository := getBlogRepository()
	err := repository.DeleteBlog(5)
	if err != nil {
		log.Println(err)
	}
}

// func TestBlogRepository_UpdateBlog_FindBlogById(t *testing.T) {
// 	repository := getBlogRepository()
// 	blogtype, err := repository.FindBlogById(1)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println(blogtype.ID, blogtype.Name)
// 	blogtype.Name = "测试一下"
// 	log.Println(blogtype.ID, blogtype.Name)
// 	repository.UpdateBlog(blogtype)
// }
