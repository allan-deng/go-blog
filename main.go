package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"allandeng.cn/allandeng/go-blog/metrics"

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
	//创建 promethues metrics
	metrics.Init()

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

	// 监听 admin api 服务 - mTLS
	go func() {
		AdminAPIListenAndServe()
	}()

	// 监听 web 服务
	srv := ":" + strconv.Itoa(int(config.GlobalConfig.WebServer.Port))
	log.Infof("Start server on %s", srv)
	err = http.ListenAndServe(srv, muxRouter) // 设置监听的端口
	if err != nil {
		log.Panicf("Panic create server failed: %v", err)
	}

}

func AdminAPIListenAndServe() {
	// register admin api router
	adminRouter := mux.NewRouter()
	router.AdminApiRegister(adminRouter)

	// load CA certificate file and add it to list of client CAs
	caCertFile, err := ioutil.ReadFile(config.GlobalConfig.MTls.CaCert)
	if err != nil {
		log.Panicf("error reading CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:                caCertPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	tlsConfig.BuildNameToCertificate()

	// serve on port 8443
	server := http.Server{
		Addr:      fmt.Sprintf(":%d", config.GlobalConfig.WebServer.AdminPort),
		Handler:   adminRouter,
		TLSConfig: tlsConfig,
	}

	log.Infof("(HTTPS) Listen on :%d\n", config.GlobalConfig.WebServer.AdminPort)
	if err := server.ListenAndServeTLS(config.GlobalConfig.MTls.ServerCert, config.GlobalConfig.MTls.ServerKey); err != nil {
		log.Fatalf("(HTTPS) error listening to port: %v", err)
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
