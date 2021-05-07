package model_test

import (
	"fmt"
	"log"
	"testing"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func TestCreateTable(t *testing.T) {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	db.SingularTable(true)
	err = db.AutoMigrate(&model.Blog{}, &model.Type{}, &model.Comment{}, &model.Tag{}, &model.User{}).Error
	if err != nil {
		log.Println(err)
	}
}

func TestBlogModel(t *testing.T) {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	//禁止复表
	db.SingularTable(true)
	err = db.AutoMigrate(&model.Blog{}, &model.Type{}).Error
	if err != nil {
		log.Println(err)
	}
	db.Model(&model.Blog{}).AddForeignKey("type_id", "types(id)", "RESTRICT", "RESTRICT")
	blogtype := &model.Type{
		Name: "TestType",
	}
	blog := &model.Blog{
		Title:          "test",
		Content:        "asdasdasdaxc牛逼dasd1231",
		FirstPicture:   "/pic/666.jpg",
		Flag:           "原创",
		Views:          0,
		Appreciation:   false,
		ShareStatement: true,
		Commentabled:   false,
		Published:      true,
		Recommend:      false,
		TypeID:         1,
	}
	insertBlog(blogtype, db)
	insertBlog(blog, db)
}

func TestFindTypeAndBlog(t *testing.T) {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	db.SingularTable(true)
	var at model.Type
	db.Debug().First(&at)
	db.Debug().Model(&at).Related(&at.Blogs).Find(&at.Blogs)
	var at2 model.Type
	var at3 []model.Type
	db.Debug().Preload("Blogs").First(&at2)
	db.Debug().Preload("Blogs").Find(&at3)
	// db.Debug().Model(&at).Association("Blogs").Find(&at.Blogs)
	fmt.Println(at)
	fmt.Println(at2)
	fmt.Println(at3)
	// v, ok := res.Value.([]model.Blog)
	// if ok {
	// 	for b, _ := range v {
	// 		fmt.Println(b)
	// 	}
	// } else {
	// 	fmt.Println("err")
	// }
}

func insertBlog(data interface{}, db *gorm.DB) error {
	temp := data
	res := db.Create(temp)
	if res.Error != nil {
		return db.Error
	}
	return nil
}
