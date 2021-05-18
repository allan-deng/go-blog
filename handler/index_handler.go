package handler

import (
	"html/template"
	"net/http"
	"path"
	"strconv"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/repository"
	"github.com/Masterminds/sprig"
)

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
	pblog := repository.NewPage(pagenum, 10)
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
