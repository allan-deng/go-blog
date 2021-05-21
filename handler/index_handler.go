package handler

import (
	"html/template"
	"net/http"
	"path"
	"strconv"

	"allandeng.cn/allandeng/go-blog/config"
	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/service"
	"allandeng.cn/allandeng/go-blog/util"
	"github.com/Masterminds/sprig"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// path:"/"
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})

	//get types
	ptpye := repository.NewPage(1, 6)
	types, err := repository.GetTypeRepository().FindTop(&ptpye)
	model["types"] = types
	if err != nil {
		panic(err)
	}
	//get types
	ptag := repository.NewPage(1, 10)
	tags, err := repository.GetTagRepository().FindTop(&ptag)
	model["tags"] = tags
	if err != nil {
		panic(err)
	}

	query := r.URL.Query()
	pagenum := 1
	if pagestr := query.Get("page"); pagestr != "" {
		pagenum, err = strconv.Atoi(pagestr)
		if err != nil || pagenum <= 0 {
			pagenum = 1
		}
	}
	//get blogs
	pblog := repository.NewPage(pagenum, 8)
	blogs, err := repository.GetBlogRepository().FindAll(&pblog)
	pblog.Update()
	if pblog.Nums < pblog.Index {
		pblog.Index = 1
		blogs, err = repository.GetBlogRepository().FindAll(&pblog)
	}
	model["blogs"] = blogs
	if err != nil {
		panic(err)
	}

	//get recommand blogs
	precommands := repository.NewPage(1, 8)
	recommands, err := repository.GetBlogRepository().FindRecommendTop(&precommands)
	model["recommands"] = recommands
	if err != nil {
		panic(err)
	}

	pblog.Update()
	model["page"] = pblog
	model["pagetitle"] = "首页"
	model["active"] = 1
	model["massage"] = config.GlobalMassage
	base := path.Base("views/index.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/index.html", "views/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}

// path:"/footer/newblog"
func NewBlogHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})
	//get recommand blogs
	pnewblogs := repository.NewPage(1, 3)
	newblogs, err := repository.GetBlogRepository().FindRecommendTop(&pnewblogs)
	model["newblogs"] = newblogs
	if err != nil {
		panic(err)
	}
	base := path.Base("views/_fragments-newblog.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/_fragments-newblog.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}

// path:"/search"
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	query := r.PostFormValue("query")
	log.Debugf("Received the form: %v,Search blog with : %s", r.PostForm, query)
	model := make(map[string]interface{})
	model["query"] = query

	//search blogs
	pblog := repository.NewPage(1, 1000)
	blogs, err := repository.GetBlogRepository().FindByQuery(query, &pblog)
	pblog.Update()
	model["blogs"] = blogs
	model["page"] = pblog
	if err != nil {
		panic(err)
	}

	model["pagetitle"] = "搜索结果:" + query
	model["active"] = 0
	model["massage"] = config.GlobalMassage
	base := path.Base("views/search.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/search.html", "views/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}

}

// path: /blog/{id}
func BlogHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})

	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		log.Errorf("Error blog format error, id: %s ,err: %s", idstr, err)
		panic(err)
	}

	blog, err := service.GetAndConvertBlog(int64(id))
	if err != nil {
		log.Errorf("Error blog not found, id: %s ,err: %s", idstr, err)
		NotFoundHandler(w, r)
		return
	}
	model["blog"] = blog

	comments, err := service.ListCommentByBlogId(int64(id))
	if err != nil {
		log.Errorf("Error can't find blog comments, blogid: %s ,err: %s", idstr, err)
	}
	model["comments"] = comments

	str, image := util.GetCaptcha(4)

	session, _ := store.Get(r, "cookie-name")
	session.Values["captcha"] = str
	session.Save(r, w)

	model["captcha"] = image

	if _, ok := session.Values["user"].(blogmodel.User); ok {
		model["admin"] = true
	}

	model["pagetitle"] = blog.Title
	model["active"] = 0
	model["massage"] = config.GlobalMassage
	model["commentmassage"] = ""
	base := path.Base("views/blog.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/blog.html", "views/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}

// path:"/captcha"
func CaptchaHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})

	str, image := util.GetCaptcha(4)
	session, _ := store.Get(r, "cookie-name")
	session.Values["captcha"] = str
	session.Save(r, w)

	model["captcha"] = image
	base := path.Base("views/_fragments-captcha.html")
	err := template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/_fragments-captcha.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	model := make(map[string]interface{})
	model["massage"] = config.GlobalMassage
	model["active"] = 0
	model["pagetitle"] = "500-服务器错误"
	base := path.Base("views/error/500.html")
	err := template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/error/500.html", "views/_fragments.html")).Execute(w, model)
	if err != nil {
		log.Errorf("Error in error_handler : %s", err)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	model := make(map[string]interface{})
	model["massage"] = config.GlobalMassage
	model["active"] = 0
	model["pagetitle"] = "404-找不到网页"
	base := path.Base("views/error/404.html")
	err := template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/error/404.html", "views/_fragments.html")).Execute(w, model)
	if err != nil {
		log.Errorf("Error in error_handler : %s", err)
	}
}
