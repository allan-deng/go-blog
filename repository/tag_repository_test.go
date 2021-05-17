package repository

import (
	"log"
	"testing"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getTagRepository() ITagRepository {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}
	db.SingularTable(true)
	InitTagRepository(db)
	return GetTagRepository()
}

func TestTagRepository_CreateTag(t *testing.T) {
	repository := getTagRepository()
	blogtype := &model.Tag{
		Name: "test",
	}
	id, err := repository.CreateTag(blogtype)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)

}

func TestTagRepository_DeleteTag(t *testing.T) {
	repository := getTagRepository()
	err := repository.DeleteTag(5)
	if err != nil {
		log.Println(err)
	}
}

func TestTagRepository_UpdateTag_FindTagById(t *testing.T) {
	repository := getTagRepository()
	blogtype, err := repository.FindTagById(1)
	if err != nil {
		log.Println(err)
	}
	log.Println(blogtype.ID, blogtype.Name)
	blogtype.Name = "测试一下"
	log.Println(blogtype.ID, blogtype.Name, blogtype)
	repository.UpdateTag(blogtype)

}

func TestTagRepository_FindTagByName(t *testing.T) {
	repository := getTagRepository()
	blogtype, err := repository.FindTagByName("测试一下")
	if err != nil {
		log.Println(err)
	}
	log.Println(blogtype.ID, blogtype.Name)
}

func TestTagRepository_FindTop(t *testing.T) {
	repository := getTagRepository()
	page := &Page{
		Index: 1,
		Size:  100,
	}
	res, err := repository.FindTop(page)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	log.Println(page.Count)
}

func TestTagRepository_FindAll(t *testing.T) {
	repository := getTagRepository()
	res, err := repository.FindAll()
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
