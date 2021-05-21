package handler

import (
	"html/template"
	"net/http"
	"path"
	"strconv"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/repository"
	"github.com/Masterminds/sprig"
	"github.com/gorilla/mux"
)

// path:"/types/{id}"
func TypeHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})
	//获取id
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id == 0 || id < -1 {
		log.Errorf("Error type format error, id: %s ,err: %s", idstr, err)
		NotFoundHandler(w, r)
		return
	}
	//获取页码
	query := r.URL.Query()
	pagenum := 1
	if pagestr := query.Get("page"); pagestr != "" {
		pagenum, err := strconv.Atoi(pagestr)
		if err != nil || pagenum <= 0 {
			pagenum = 1
		}
	}
	//获取所有type
	typepage := repository.NewPage(1, 10000)
	types, err := repository.GetTypeRepository().FindTop(&typepage)
	if err != nil {
		log.Errorf("Error cant find all types. err: %s", err)
		panic(err)
	}
	if id == -1 {
		id = int(types[0].ID)
	}
	model["types"] = types
	//获取对应type的所有博客
	pblog := repository.NewPage(pagenum, 8)
	blogs, err := repository.GetBlogRepository().FindBlogByTypeId(int64(id), &pblog)
	pblog.Update()
	if pblog.Nums < pblog.Index {
		pblog.Index = 1
		blogs, err = repository.GetBlogRepository().FindBlogByTypeId(int64(id), &pblog)
	}
	model["blogs"] = blogs
	if err != nil {
		panic(err)
	}

	pblog.Update()
	model["page"] = pblog
	model["pagetitle"] = "分类"
	model["active"] = 2
	model["activetypeid"] = id
	model["massage"] = config.GlobalMassage
	base := path.Base("views/types.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/types.html", "views/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}

// path:"/tags/{id}"
func TagHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})
	//获取id
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id == 0 || id < -1 {
		log.Errorf("Error type format error, id: %s ,err: %s", idstr, err)
		NotFoundHandler(w, r)
		return
	}
	//获取页码
	query := r.URL.Query()
	pagenum := 1
	if pagestr := query.Get("page"); pagestr != "" {
		pagenum, err := strconv.Atoi(pagestr)
		if err != nil || pagenum <= 0 {
			pagenum = 1
		}
	}
	//获取所有tag
	tagpage := repository.NewPage(1, 10000)
	tags, err := repository.GetTagRepository().FindTop(&tagpage)
	if err != nil {
		log.Errorf("Error cant find all tags. err: %s", err)
		panic(err)
	}
	if id == -1 {
		id = int(tags[0].ID)
	}
	model["tags"] = tags
	//获取对应tag的所有博客
	pblog := repository.NewPage(pagenum, 8)
	blogs, err := repository.GetBlogRepository().FindBlogByTagId(int64(id), &pblog)
	pblog.Update()
	if pblog.Nums < pblog.Index {
		pblog.Index = 1
		blogs, err = repository.GetBlogRepository().FindBlogByTagId(int64(id), &pblog)
	}
	model["blogs"] = blogs
	if err != nil {
		panic(err)
	}

	pblog.Update()
	model["page"] = pblog
	model["pagetitle"] = "标签"
	model["active"] = 3
	model["activetypeid"] = id
	model["massage"] = config.GlobalMassage
	base := path.Base("views/tags.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/tags.html", "views/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}
