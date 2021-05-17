package repository

import (
	"log"
	"testing"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getUserRepository() IUserRepository {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}
	db.SingularTable(true)
	InitUserRepository(db)
	return GetUserRepository()
}

func TestUserRepository_CreateUser(t *testing.T) {
	repository := getUserRepository()
	blogtype := &model.User{
		Nickname: "nick",
		Username: "user",
		Password: "p1",
		Email:    "111",
		Avatar:   "222",
	}
	id, err := repository.CreateUser(blogtype)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)

}

func TestUserRepository_DeleteUser(t *testing.T) {
	repository := getUserRepository()
	err := repository.DeleteUser(1)
	if err != nil {
		log.Println(err)
	}
}

func TestUserRepository_UpdateUser_FindUserById(t *testing.T) {
	repository := getUserRepository()
	blogtype, err := repository.FindUserById(2)
	if err != nil {
		log.Println(err)
	}
	log.Println(blogtype.ID, blogtype.Nickname)
	blogtype.Nickname = "测试一下"
	log.Println(blogtype.ID, blogtype.Nickname)
	repository.UpdateUser(blogtype)

}

func TestUserRepository_FindUserByNameAndPassword(t *testing.T) {
	repository := getUserRepository()
	user, err := repository.FindUserByNameAndPassword("user", "p1")
	if err != nil {
		log.Println(err)
	}
	log.Println(user)
}
