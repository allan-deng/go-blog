package repository

import (
	"log"
	"testing"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getRepository() ITypeRepository {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}
	db.SingularTable(true)
	InitTypeRepository(db)
	return GetTypeRepository()
}

func TestTypeRepository_CreateType(t *testing.T) {
	repository := getRepository()
	blogtype := &model.Type{
		Name: "test",
	}
	id, err := repository.CreateType(blogtype)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)

}

func TestTypeRepository_DeleteType(t *testing.T) {
	repository := getRepository()
	err := repository.DeleteType(5)
	if err != nil {
		log.Println(err)
	}
}

func TestTypeRepository_UpdateType_FindTypeById(t *testing.T) {
	repository := getRepository()
	blogtype, err := repository.FindTypeById(1)
	if err != nil {
		log.Println(err)
	}
	log.Println(blogtype.ID, blogtype.Name)
	blogtype.Name = "测试一下"
	log.Println(blogtype.ID, blogtype.Name, blogtype)
	repository.UpdateType(blogtype)

}

func TestTypeRepository_FindTypeByName(t *testing.T) {
	repository := getRepository()
	blogtype, err := repository.FindTypeByName("测试一下")
	if err != nil {
		log.Println(err)
	}
	log.Println(blogtype.ID, blogtype.Name, blogtype)
}

func TestTypeRepository_FindTop(t *testing.T) {
	repository := getRepository()
	page := &Page{
		Index: 1,
		Size:  5,
	}
	res, err := repository.FindTop(page)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	log.Println(page.Count)
}

func TestTypeRepository_FindAll(t *testing.T) {
	repository := getRepository()
	res, err := repository.FindAll()
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
