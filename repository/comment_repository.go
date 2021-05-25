package repository

import (
	"time"

	"allandeng.cn/allandeng/go-blog/model"
	"github.com/jinzhu/gorm"
)

/*
TODO:
按照博客id查找顶级评论
List<Comment> findByCommentIdAndParentCommentNull(Long CommentId, Sort sort)
*/

type ICommentRepository interface {
	//建表
	InitTable() error
	//增
	CreateComment(*model.Comment) (int64, error)
	//删
	DeleteComment(int64) error
	//改
	UpdateComment(*model.Comment) error
	//用Commentid查找顶级评论
	FindParentCommentByCommentId(CommentId int64) (model.Comment, error)
	//用blogid查找顶级评论
	FindParentCommentByBlogId(blogid int64) ([]model.Comment, error)
	//用blogid查找所有评论
	FindAllCommentByBlogId(blogid int64) ([]model.Comment, error)
	//按id查找
	FindCommentById(int64) (*model.Comment, error)
	//用父评论id找到所有的子评论
	FindChildren(praentId int64) ([]model.Comment, error)
}

type CommentRepository struct {
	mysqlDb *gorm.DB
}

//创建CommentRepository
func newCommentRepository(db *gorm.DB) ICommentRepository {
	return &CommentRepository{mysqlDb: db}
}

/*
这里使用单例模式，但是把dao层的init和get分开，调用之前人为显式的完成初始化
*/
var commentRepositoryIns ICommentRepository

func InitCommentRepository(db *gorm.DB) error {
	commentRepositoryIns = newCommentRepository(db)
	return nil
}

func GetCommentRepository() ICommentRepository {
	return commentRepositoryIns
}

//建表
func (s *CommentRepository) InitTable() error {
	return s.mysqlDb.AutoMigrate(&model.Comment{}).Error
}

//增
func (s *CommentRepository) CreateComment(comment *model.Comment) (int64, error) {
	comment.CreateTime = time.Now()
	return comment.ID, s.mysqlDb.Create(comment).Error
}

//删
func (s *CommentRepository) DeleteComment(commentId int64) error {
	// 先所有子评论
	err := s.mysqlDb.Exec("delete from comment where parent_comment_id = ?", commentId).Error
	if err != nil {
		return err
	}
	return s.mysqlDb.Delete(&model.Comment{}, commentId).Error
}

//改
func (s *CommentRepository) UpdateComment(comment *model.Comment) error {
	return s.mysqlDb.Save(comment).Error
}

//用Commentid查找顶级评论
func (s *CommentRepository) FindParentCommentByCommentId(commentId int64) (model.Comment, error) {
	var res model.Comment
	err := s.mysqlDb.Where("id = ?", commentId).First(&res).Error
	return res, err
}

//用blogid查找顶级评论
func (s *CommentRepository) FindParentCommentByBlogId(blogid int64) ([]model.Comment, error) {
	var res []model.Comment
	err := s.mysqlDb.Order("create_time ASC").Where("parent_comment_id = 0").Where("blog_id = ?", blogid).Find(&res).Error
	return res, err
}

//用blogid查找所有评论
func (s *CommentRepository) FindAllCommentByBlogId(blogid int64) ([]model.Comment, error) {
	var res []model.Comment
	err := s.mysqlDb.Order("create_time ASC").Where("blog_id = ?", blogid).Find(&res).Error
	return res, err
}

//按id查找
func (s *CommentRepository) FindCommentById(commentId int64) (*model.Comment, error) {
	comment := &model.Comment{}
	return comment, s.mysqlDb.First(comment, commentId).Error
}

//用父评论id找到所有的子评论
func (s *CommentRepository) FindChildren(praentId int64) ([]model.Comment, error) {
	var res []model.Comment
	err := s.mysqlDb.Order("create_time ASC").Where("parent_comment_id = ?", praentId).Find(&res).Error
	return res, err
}
