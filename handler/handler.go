package handler

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"github.com/gorilla/sessions"
)

//handler函数，需要传递上下文
type HandlerFunc func(ctx *Context, w http.ResponseWriter, r *http.Request)

//handlerError的集中处理函数
type HandlerErrorTerminate func(errs ...HandlerError)

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
func NewHandlerError(err string, r *http.Request) HandlerError {
	handlerError := &HandlerError{}
	//Caller入参为1，即获取NewHandlerError调用时的相关信息
	funcName, file, line, ok := runtime.Caller(1)
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

//HandlerFunc到http.HandlerFunc的包装类
func ContextHandler(errHandler HandlerErrorTerminate, handlerList ...HandlerFunc) http.HandlerFunc {
	ctx := &Context{
		Model:     make(map[string]interface{}),
		Attribute: make(map[string]interface{}),
		Errors:    make([]HandlerError, 0),
	}
	return func(w http.ResponseWriter, r *http.Request) {
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
			handler(ctx, w, r)
		}
		//ctx中错误处理
		errHandler(ctx.Errors...)
	}
}

//默认的HandlerError处理
func DefaultHandlerErrorTerminate(errs ...HandlerError) {
	for _, err := range errs {
		log.Error(err.Error())
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
				log.Debugf("Already logged in user: %s", user.Nickname)
				handler(ctx, w, r)
			}
		}
		ctx.Errors = append(ctx.Errors, NewHandlerError("Not login!", r))
		unauthHandler(ctx, w, r)
	}
}
