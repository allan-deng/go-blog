package repository

import (
	"log"
	"testing"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func getCommentRepository() ICommentRepository {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True")
	if err != nil {
		log.Println(err)
	}
	db.SingularTable(true)
	InitCommentRepository(db)
	return GetCommentRepository()
}

func TestCommentRepository_CreateComment(t *testing.T) {
	repository := getCommentRepository()

	Commenttype := &model.Comment{
		Nickname:        "1",
		Email:           "2",
		Content:         "3",
		Avatar:          "4",
		BlogID:          11,
		ParentCommentID: 0,
		AdminComment:    false,
	}
	id, err := repository.CreateComment(Commenttype)
	if err != nil {
		log.Println(err)
	}
	log.Println(id)

}

func TestCommentRepository_DeleteComment(t *testing.T) {
	repository := getCommentRepository()
	err := repository.DeleteComment(1)
	if err != nil {
		log.Println(err)
	}
}

func TestCommentRepository_UpdateComment_FindCommentById(t *testing.T) {
	repository := getCommentRepository()
	comment, err := repository.FindCommentById(11)
	if err != nil {
		log.Println(err)
	}
	log.Println(comment)
	comment.Nickname = "666"
	log.Println(comment)
	repository.UpdateComment(comment)
}
