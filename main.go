package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

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

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
	// 注意:如果没有调用 ParseForm 方法，下面无法获取表单的数据
	fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // 这个写入到 w 的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./views/login.html")
		log.Debug(t.Execute(w, nil))
	} else {
		err := r.ParseForm() // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
		if err != nil {
			// handle error http.Error() for example
			log.Fatal("ParseForm: ", err)
		}
		// 请求的是登录数据，那么执行登录的逻辑判断
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // 输出到服务器端
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")

	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("./views/upload.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)                 //设置临时文件大小
		file, handler, err := r.FormFile("uploadfile") //获取文件句柄
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)

		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

type MyMux struct {
	// handlerMap map[string]func(http.ResponseWriter, http.Request)
}

// func (s *MyMux) InitMux() {
// 	s.handlerMap = make(map[string]func(http.ResponseWriter, http.Request))
// }

// func (s *MyMux) AddRouter(path string, handler func(http.ResponseWriter, http.Request)) {
// 	s.handlerMap[path] = handler
// }

func (s *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		sayhelloName(w, r)
	case "/login/":
		login(w, r)
	case "/upload":
		upload(w, r)
	// case "/blog/": //如果需要在path中传入参数在路径后添加一个/，自己在handler中解析地址
	// 	blog(w, r)
	case "/admim": //使用wrapper函数实验拦截器
		HandleIterceptor(login)
	default:
		http.NotFound(w, r)
	}
	// if r.URL.Path == "/" {
	// 	sayhelloName(w, r)
	// 	return
	// }
	// http.NotFound(w, r)
	// return
}

//如果需要使用拦截器，使用wrapper实现
func HandleIterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handleIterceptor")
		h(w, r)
	}
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

	//路由
	// http.HandleFunc("/", sayhelloName)              // 设置访问的路由
	// mux := &MyMux{}
	// err := http.ListenAndServe("0.0.0.0:9090", mux) // 设置监听的端口
	err = http.ListenAndServe("0.0.0.0:9090", muxRouter) // 设置监听的端口
	if err != nil {
		log.Panicf("Panic create server failed : ", err)
	}

}

//连接数据库的函数
func GetDbConnect(user string, password string, host string, dbname string) (*gorm.DB, error) {
	dbString := user + ":" + password + "@tcp(" + host + ")/" + dbname + "?charset=utf8&parseTime=True"
	log.Debugf("Connect database string : %s", dbString)
	db, err := gorm.Open("mysql", dbString)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	return db, nil
}

//初始化repository和对应的表
func InitDaoAndTable(db *gorm.DB) {
	//TODO: 此处写法不优雅，后续再细究
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
