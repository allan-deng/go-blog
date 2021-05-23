package handler

import (
	"html/template"
	"net/http"
	"path"
	"strconv"

	"allandeng.cn/allandeng/go-blog/config"
	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"github.com/Masterminds/sprig"
)

// path:"/admin/blogs"
func AdminBlogListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		AdminBlogCreateHandler(w, r)
	} else {
		AdminBlogListGetHandler(w, r)
	}
}

/*
前端返回的形式
published: false
id: 0
flag: 原创
title: 1231231
content: 3123123
type.id: 1
tagIds: 1,2,3
firstPicture: 123
description: 123123
appreciation: on
commentabled: on
bool类型的为false的直接没有这个字段，为true的返回为on
*/
// path:"/admin/blogs" method:"post"
func AdminBlogCreateHandler(w http.ResponseWriter, r *http.Request) {

}

// path:"/admin/blogs" method:"get"
func AdminBlogListGetHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})
	//应该抽象出来
	session, err := store.Get(r, "cookie-name")
	for key, value := range session.Values {
		log.Debugf("key: %v , value: %v", key, value)
	}
	if err != nil {
		log.Errorf("Error get session failed, host: %s,cookie-name： %s, err: %s", r.Host, r.Header.Get("cookie-name"), err)
		panic(err)
	}
	var user blogmodel.User
	var ok bool
	if user, ok = session.Values["user"].(blogmodel.User); !ok {
		log.Errorf("Error user not login!")
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
		return
	}
	model["user"] = user

	pblog := repository.NewPage(1, 100000)
	blogs, err := repository.GetBlogRepository().FindAll(&pblog)
	if err != nil {
		log.Errorf("Error cant find all blog, err:%s", err)
		panic(err)
	}
	types, err := repository.GetTypeRepository().FindAll()
	if err != nil {
		log.Errorf("Error cant find all blog, err:%s", err)
		panic(err)
	}
	model["blogs"] = blogs
	model["types"] = types
	model["count"] = pblog.Count
	model["pagetitle"] = "博客列表"
	model["active"] = 1
	model["secondactive"] = 0
	model["massage"] = config.GlobalMassage
	model["notice"] = ""
	base := path.Base("views/admin/blogs.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/admin/blogs.html", "views/admin/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}

// path:"/admin/blogs/search" method:"post"
func AdminBlogSearchHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})
	//应该抽象出来
	session, err := store.Get(r, "cookie-name")
	for key, value := range session.Values {
		log.Debugf("key: %v , value: %v", key, value)
	}
	if err != nil {
		log.Errorf("Error get session failed, host: %s,cookie-name： %s, err: %s", r.Host, r.Header.Get("cookie-name"), err)
		panic(err)
	}
	var user blogmodel.User
	var ok bool
	if user, ok = session.Values["user"].(blogmodel.User); !ok {
		log.Errorf("Error user not login!")
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
		return
	}
	model["user"] = user

	if r.Method != "POST" {
		http.Redirect(w, r, "/admin/blogs", http.StatusTemporaryRedirect)
		return
	}
	err = r.ParseForm()
	if err != nil {
		panic(err)
	}

	title := r.PostFormValue("title")
	typeid := r.PostFormValue("typeId")
	recommend := r.PostFormValue("recommend")

	log.Debugf("Received the form: %v, title: %s, typeid: %s, recommend: %s", r.PostForm, title, typeid, recommend)
	tid, _ := strconv.Atoi(typeid)
	rec := false
	if recommend == "true" {
		rec = true
	}
	blogs, err := repository.GetBlogRepository().FindBlogByTitleAndTypeIdAndRecommend(title, int64(tid), rec)
	if err != nil {
		panic(err)
	}

	model["blogs"] = blogs
	base := path.Base("views/admin/_fragments-bloglist.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/admin/_fragments-bloglist.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}

func AdminBlogInputHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})
	//应该抽象出来
	session, err := store.Get(r, "cookie-name")
	for key, value := range session.Values {
		log.Debugf("key: %v , value: %v", key, value)
	}
	if err != nil {
		log.Errorf("Error get session failed, host: %s,cookie-name： %s, err: %s", r.Host, r.Header.Get("cookie-name"), err)
		panic(err)
	}
	var user blogmodel.User
	var ok bool
	if user, ok = session.Values["user"].(blogmodel.User); !ok {
		log.Errorf("Error user not login!")
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
		return
	}
	model["user"] = user

	types, err := repository.GetTypeRepository().FindAll()
	if err != nil {
		log.Errorf("Error cant find all type, err:%s", err)
		panic(err)
	}

	tags, err := repository.GetTagRepository().FindAll()
	if err != nil {
		log.Errorf("Error cant find all tags, err:%s", err)
		panic(err)
	}
	model["blog"] = &blogmodel.Blog{}
	model["types"] = types
	model["tags"] = tags
	model["pagetitle"] = "创建博客"
	model["active"] = 1
	model["secondactive"] = 1
	model["massage"] = config.GlobalMassage
	model["notice"] = ""
	base := path.Base("views/admin/blogs-input.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/admin/blogs-input.html", "views/admin/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}
