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

	//首页
	r.HandleFunc("/", BlogHandlerWrapper(handler.IndexHandler))
	//搜索结果
	r.HandleFunc("/search", BlogHandlerWrapper(handler.SearchHandler))
	//博客页面
	r.HandleFunc("/blog/{id:[0-9]+}", BlogHandlerWrapper(handler.BlogHandler))
	//获取验证码
	r.HandleFunc("/captcha", BlogHandlerWrapper(handler.CaptchaHandler))
	//获取footer中newblog
	r.HandleFunc("/footer/newblog", BlogHandlerWrapper(handler.NewBlogHandler))
	//type页面
	r.HandleFunc("/types/{id}", BlogHandlerWrapper(handler.TypeHandler))
	//tag页面
	r.HandleFunc("/tags/{id}", BlogHandlerWrapper(handler.TagHandler))
	//archive页面
	r.HandleFunc("/archives", BlogHandlerWrapper(handler.ArchiveHandler))
	//about
	r.HandleFunc("/about", BlogHandlerWrapper(handler.AboutHandler))
	//创建comments
	r.HandleFunc("/comments", BlogHandlerWrapper(handler.CommentCreateHandler))
	//删除comments
	r.HandleFunc("/comments/delete/{id:[0-9]+}", AuthBlogHandlerWrapper(handler.CommentDeleteHandler))
	//upload
	r.HandleFunc("/uploadfile", AuthBlogHandlerWrapper(handler.UploadHandler))

	//管理端的router
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/", LoginBlogHandlerWrapper(handler.LoginPageHandler))
	adminRouter.HandleFunc("", LoginBlogHandlerWrapper(handler.LoginPageHandler))
	adminRouter.HandleFunc("/login", BlogHandlerWrapper(handler.LoginHandler))
	adminRouter.HandleFunc("/logout", BlogHandlerWrapper(handler.LogoutHandler))
	adminRouter.HandleFunc("/index", AuthBlogHandlerWrapper(handler.AdminIndexHandler))

	adminRouter.HandleFunc("/blogs", AuthBlogHandlerWrapper(handler.AdminBlogListHandler))                      //get:博客列表 post:上传博客
	adminRouter.HandleFunc("/blogs/search", AuthBlogHandlerWrapper(handler.AdminBlogSearchHandler))             //post:搜索
	adminRouter.HandleFunc("/blogs/input", AuthBlogHandlerWrapper(handler.AdminBlogInputHandler))               //get:返回写博客页面
	adminRouter.HandleFunc("/blogs/{id:[0-9]+}/input", AuthBlogHandlerWrapper(handler.AdminBlogUpdateHandler))  //get:更新博客页面
	adminRouter.HandleFunc("/blogs/{id:[0-9]+}/delete", AuthBlogHandlerWrapper(handler.AdminBlogDeleteHandler)) //get:删除博客

	//admin type 操作
	adminRouter.HandleFunc("/types", AuthBlogHandlerWrapper(handler.AdminTypesHandler))
	adminRouter.HandleFunc("/types/{id:[0-9]+}/input", AuthBlogHandlerWrapper(handler.AdminTypesUpdateHandler))
	adminRouter.HandleFunc("/types/{id:[0-9]+}/delete", AuthBlogHandlerWrapper(handler.AdminTypesDeleteHandler))
	adminRouter.HandleFunc("/types/input", AuthBlogHandlerWrapper(handler.AdminTypesInputHandler))
	//admin tag 操作
	adminRouter.HandleFunc("/tags", AuthBlogHandlerWrapper(handler.AdminTagsHandler))
	adminRouter.HandleFunc("/tags/{id:[0-9]+}/input", AuthBlogHandlerWrapper(handler.AdminTagsUpdateHandler))
	adminRouter.HandleFunc("/tags/{id:[0-9]+}/delete", AuthBlogHandlerWrapper(handler.AdminTagsDeleteHandler))
	adminRouter.HandleFunc("/tags/input", AuthBlogHandlerWrapper(handler.AdminTagsInputHandler))

	r.NotFoundHandler = http.HandlerFunc(BlogHandlerWrapper(handler.NotFoundHandler))
}

//记录请求的相关日志
func logMidware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 请求信息写入日志
		log.Infof("Receipt request, URL:%s, Method:%s, Host:%s, Remote:%s", r.RequestURI, r.Method, r.Host, r.RemoteAddr)
		log.Debugf("Receipt request : \n %s", getHttpRequestInfo(r, true))
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

		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(s))
	}
	return buffer.String()
}

func BlogHandlerWrapper(handlerList ...handler.HandlerFunc) http.HandlerFunc {
	return handler.ContextHandler(handler.DefaultHandlerErrorTerminate, handler.BlogContextInit, handlerList...)
}

func AuthBlogHandlerWrapper(handlerList ...handler.HandlerFunc) http.HandlerFunc {
	for i := 0; i < len(handlerList); i++ {
		handlerList[i] = handler.AuthWrapper(handlerList[i], handler.LoginPageHandler)
	}
	return handler.ContextHandler(handler.DefaultHandlerErrorTerminate, handler.BlogContextInit, handlerList...)
}

func LoginBlogHandlerWrapper(handlerList ...handler.HandlerFunc) http.HandlerFunc {
	for i := 0; i < len(handlerList); i++ {
		handlerList[i] = handler.AuthWrapper(handlerList[i], handler.LoginPageHandler)
	}
	return handler.ContextHandler(handler.DefaultHandlerErrorTerminate, handler.BlogContextInit, handlerList...)
}
