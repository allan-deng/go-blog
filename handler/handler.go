package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/metrics"
	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/util"
	"github.com/Masterminds/sprig"
	"github.com/gorilla/sessions"
)

//handler函数，需要传递上下文
type HandlerFunc func(ctx *Context, w http.ResponseWriter, r *http.Request)

//handlerError的集中处理函数
type HandlerErrorTerminate func(ctx *Context, w http.ResponseWriter, r *http.Request)

//context初始化函数
type ContextInitFunc func(ctx *Context)

//handler上下文
type Context struct {
	//传递给前端的模型
	Model map[string]interface{}
	//传递给后端的属性
	Attribute map[string]interface{}
	//当前请求的session
	Session *sessions.Session
	//请求处理过程的所有错误
	Errors []HandlerError
	//下一个执行的handler
	NextHandler HandlerFunc
}

//添加错误
func (s *Context) AddError(r *http.Request, format string, a ...interface{}) {
	s.Errors = append(s.Errors, NewHandlerError(fmt.Sprintf(format, a...), r, 2))
}

//判断是否存在错误
func (s *Context) HasError() bool {
	if len(s.Errors) > 0 {
		return true
	}
	return false
}

//添加下一个执行的handler
func (s *Context) Next(nextHandler HandlerFunc) {
	s.NextHandler = nextHandler
}

//判断是否有下一个handler
func (s *Context) HasNextHandler() bool {
	if s.NextHandler == nil {
		return false
	}
	return true
}

//请求处理过程中的错误
type HandlerError struct {
	funcName string
	fileName string
	line     int
	url      string
	host     string
	method   string
	err      string
}

//获取一个HandlerError
func NewHandlerError(err string, r *http.Request, skip int) HandlerError {
	handlerError := &HandlerError{}
	//Caller入参为调用栈的层级
	funcName, file, line, ok := runtime.Caller(skip)
	if ok {
		handlerError.funcName = runtime.FuncForPC(funcName).Name()
		handlerError.fileName = file
		handlerError.line = line
	}
	handlerError.err = err

	handlerError.url = r.URL.Path
	handlerError.host = r.Host
	handlerError.method = r.Method
	return *handlerError
}

func (s *HandlerError) Error() string {
	line := strconv.Itoa(s.line)
	return fmt.Sprintf("Hadnler error occurred. func: %s - %s:%s; Url: %s|%s; Host: %s; Err: %s ", s.funcName, s.fileName, line, s.url, s.method, s.host, s.err)
}

//HandlerFunc到http.HandlerFunc的包装
func ContextHandler(errHandler HandlerErrorTerminate, contextinit ContextInitFunc, handlerList ...HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := util.NewResponseWriter(w)

		ctx := &Context{
			Model:     make(map[string]interface{}),
			Attribute: make(map[string]interface{}),
			Errors:    make([]HandlerError, 0),
		}
		contextinit(ctx)

		//获取session
		session, err := store.Get(r, "cookie-name")
		for key, value := range session.Values {
			log.Debugf("session: key: %v , value: %v", key, value)
		}
		if err != nil {
			log.Errorf("Error get session failed, host: %s,cookie-name： %s, err: %s", r.Host, r.Header.Get("cookie-name"), err)
		}
		ctx.Session = session

		//处理所有handler
		for _, handler := range handlerList {
			if ctx.HasError() {
				break
			}
			handler(ctx, rw, r)
			if ctx.HasNextHandler() {
				ctx.NextHandler(ctx, rw, r)
			}
		}

		//ctx中错误处理
		errHandler(ctx, rw, r)

		//指标上报
		elapsed := time.Since(start)
		url := r.RequestURI
		method := r.Method
		code := rw.Status()

		metrics.RequestDurations.WithLabelValues(url, method, strconv.Itoa(code)).Observe(float64(elapsed.Milliseconds()))
		metrics.TotalRequests.WithLabelValues(url, method, strconv.Itoa(code)).Inc()
	}
}

//默认的HandlerError处理
func DefaultHandlerErrorTerminate(ctx *Context, w http.ResponseWriter, r *http.Request) {
	for _, err := range ctx.Errors {
		log.Error(err.Error())
	}
	if p := recover(); p != nil {
		log.Errorf("message push panic: %v . stack: \n %s", p, string(debug.Stack()))
		ErrorHandler(ctx, w, r)
	}
}

//鉴权包装类，handler被包装的函数，unauthHandler鉴权失败时调用的函数
func AuthWrapper(handler HandlerFunc, unauthHandler HandlerFunc) HandlerFunc {
	return func(ctx *Context, w http.ResponseWriter, r *http.Request) {
		//应该抽象出来
		session := ctx.Session
		if session != nil {
			for key, value := range session.Values {
				log.Debugf("key: %v , value: %v", key, value)
			}
			var user blogmodel.User
			var ok bool
			if user, ok = session.Values["user"].(blogmodel.User); ok {
				ctx.Model["user"] = user
				log.Infof("Already logged in user: %s", user.Nickname)
				handler(ctx, w, r)
				return
			}
		}
		ctx.AddError(r, "Not login!")
		unauthHandler(ctx, w, r)
	}
}

//blog常用model的初始化
func BlogContextInit(ctx *Context) {
	ctx.Model["pagetitle"] = "Allan的个人博客"
	ctx.Model["active"] = 1
	ctx.Model["secondactive"] = 1
	ctx.Model["massage"] = config.GlobalMassage
	ctx.Model["notice"] = ""
}

//渲染模板
func RenderTemplate(ctx *Context, w http.ResponseWriter, r *http.Request, templateFiles ...string) {
	if len(templateFiles) < 1 {
		ctx.AddError(r, "No input template file!")
		return
	}
	base := path.Base(templateFiles[0])
	err := template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles(templateFiles...)).Execute(w, ctx.Model)
	if err != nil {
		ctx.AddError(r, "Render template error! err: %s", err)
	}
}
