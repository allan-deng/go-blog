package model

import (
	"time"
)

type User struct {
	ID int64
	//用户昵称
	Nickname string
	//登陆用户名
	Username string
	//登陆密码，MD5
	Password string
	//邮箱
	Email string
	//头像
	Avatar string
	//该用户的博客列表,一对多
	Blogs []Blog
}

type Type struct {
	//id
	ID int64
	//分类名称,非空
	Name string
	//属于该分类的博客列表,一对多
	Blogs []Blog //`gorm:"FOREIGNKEY:BlogID;ASSOCIATION_FOREIGNKEY:ID"` //`gorm:"foreignKey:BlogID"` //`gorm:"many2many:blog_type;"`
}

type Comment struct {
	//id
	ID int64
	//添加评论的用户名
	Nickname string `gorm:"comment:'添加评论的用户名';type:varchar(255)"`
	//评论者的邮箱
	Email string `gorm:"comment:'评论者的邮箱';type:varchar(255)"`
	//评论内容
	Content string `gorm:"comment:'评论内容';type:varchar(255)"`
	//评论头像
	Avatar string `gorm:"comment:'头像地址';type:varchar(255)"`
	//评论创建时间
	CreateTime time.Time `gorm:"comment:'创建时间';default:CURRENT_TIMESTAMP"`
	//评论关联的博客
	BlogID int64 `gorm:"comment:'博客id';NOT NULL;INDEX"`
	//提交评论的验证码：不存入数据库
	Captchacode string `gorm:"-"`
	//子评论 一对多
	ReplyComments []Comment
	//父评论
	ParentCommentID int64 `gorm:"comment:'父评论id';NOT NULL;INDEX"`
	ParentComment   *Comment
	//是否为管理员评论
	AdminComment bool `gorm:"comment:'是否为管理员评论';type:tinyint;default:0"`
}

type Tag struct {
	//id
	ID int64
	//标签名称,非空
	Name string
	//属于该标签的博客列表,多对多
	Blogs []Blog `gorm:"many2many:blog_tag;"`
}

type Blog struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`
	//标题
	Title string `gorm:"comment:'标题';type:varchar(255)"`
	//内容 大字段，懒加载
	Content string `gorm:"comment:'内容';type:longtext"`
	//首图
	FirstPicture string `gorm:"comment:'首图';type:varchar(255)"`
	//标记：原创、翻译、转载
	Flag string `gorm:"comment:'标记';type:varchar(255)"`
	//浏览次数
	Views int32 `gorm:"comment:'浏览次数';type:int"`
	//是否开启赞赏
	Appreciation bool `gorm:"comment:'是否开启赞赏';type:tinyint;default:0"`
	//是否显示分享版权
	ShareStatement bool `gorm:"comment:'是否显示分享版权';type:tinyint;default:0"`
	//是否可评论
	Commentabled bool `gorm:"comment:'是否可评论';type:tinyint;default:0"`
	//是否发布
	Published bool `gorm:"comment:'是否发布';type:tinyint;default:0"`
	//是否首页推荐
	Recommend bool `gorm:"comment:'是否首页推荐';type:tinyint;default:0"`
	//创建时间
	CreateTime time.Time `gorm:"comment:'创建时间';default:CURRENT_TIMESTAMP"` //在创建时，如果该字段值为零值，则使用当前时间填充
	//更新时间
	UpdateTime time.Time `gorm:"comment:'修改时间';default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` //在创建时该字段值为零值或者在更新时，使用当前时间戳秒数填充
	//文章类型
	TypeID int64 `gorm:"comment:'文章类型id';NOT NULL;INDEX"`
	Type   Type
	//该博客的标签,多对多
	Tags []Tag `gorm:"comment:'博客标签';many2many:blog_tag;"`
	//发布用户
	UserID int64 `gorm:"comment:'发布用户ID';NOT NULL;INDEX"`
	User   User
	//该博客的评论,一对多
	Comments []Comment
	//tag的id,用于返回给前端不存入数据库
	TagIds string `gorm:"-"`
	//博客描述
	Description string `gorm:"comment:'博客描述';type:text"`
}
