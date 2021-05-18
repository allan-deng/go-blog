package router

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/handler"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

var log *logging.Logger

func init() {
	log = config.Logger
}

func Register(r *mux.Router) {
	//添加日志记录器
	r.Use(logMidware)
	//静态文件路由
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/", handler.IndexHandler)
	//TODO:
	// r.HandleFunc("/footer/newblog", handler.NewBlogHandler)

	//管理端的router
	// adminRouter := r.PathPrefix("/admin").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)
}

//记录请求的相关日志
func logMidware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 请求信息写入日志
		log.Debugf("Receipt request : \n %s", getHttpRequestInfo(r, true))
		log.Infof("Receipt request : \n %s", getHttpRequestInfo(r, false))
		next.ServeHTTP(w, r)
	})
}

//获取请求的详细信息并转换为string
func getHttpRequestInfo(r *http.Request, detail bool) string {
	var buffer bytes.Buffer
	buffer.WriteString("Request Host: ")
	buffer.WriteString(r.Host)
	buffer.WriteString("\n")
	buffer.WriteString("Request URL: ")
	buffer.WriteString(r.RequestURI)
	buffer.WriteString("\n")
	buffer.WriteString("Remote Addr: ")
	buffer.WriteString(r.RemoteAddr)
	buffer.WriteString("\n")
	buffer.WriteString("Header: ")
	headerJson, _ := json.Marshal(r.Header)
	buffer.WriteString(string(headerJson))
	buffer.WriteString("\n")
	buffer.WriteString("ContentLength: ")
	buffer.WriteString(strconv.Itoa(int(r.ContentLength)))
	buffer.WriteString(" bytes\n")
	if detail {
		buffer.WriteString("Content: ")
		s, _ := ioutil.ReadAll(r.Body)
		buffer.WriteString(string(s))
		buffer.WriteString("\n")
	}
	return buffer.String()
}
