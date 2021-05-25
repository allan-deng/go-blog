package main

import (
	"net/http"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"allandeng.cn/allandeng/go-blog/config"
	dao "allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/router"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
)

var log *logging.Logger

func init() {
	log = config.Logger
}

func main() {
	//创建并注册路由
	muxRouter := mux.NewRouter()
	router.Register(muxRouter)
	//连接到数据库并初始化
	dbConf := config.GlobalConfig.Mysql
	db, err := GetDbConnect(dbConf.User, dbConf.Pwd, dbConf.Host, dbConf.Dbname)
	if err != nil {
		log.Panicf("Error : can't connent to database ; %s", err)
	}
	InitDaoAndTable(db)

	srv := "0.0.0.0:" + strconv.Itoa(int(config.GlobalConfig.WebServer.Port))
	log.Infof("Start server on %s", srv)
	err = http.ListenAndServe(srv, muxRouter) // 设置监听的端口
	if err != nil {
		log.Panicf("Panic create server failed : ", err)
	}

}

//连接数据库的函数
func GetDbConnect(user string, password string, host string, dbname string) (*gorm.DB, error) {
	dbString := user + ":" + password + "@tcp(" + host + ")/" + dbname + "?charset=utf8&parseTime=True"
	log.Infof("Connect database string : %s", dbString)
	db, err := gorm.Open("mysql", dbString)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	return db, nil
}

//初始化repository和对应的表
func InitDaoAndTable(db *gorm.DB) {
	dao.InitBlogRepository(db)
	dao.InitCommentRepository(db)
	dao.InitTagRepository(db)
	dao.InitTypeRepository(db)
	dao.InitUserRepository(db)

	dao.GetBlogRepository().InitTable()
	dao.GetCommentRepository().InitTable()
	dao.GetTagRepository().InitTable()
	dao.GetTypeRepository().InitTable()
	dao.GetUserRepository().InitTable()
}
