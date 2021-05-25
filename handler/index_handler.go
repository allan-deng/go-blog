package handler

import (
	"net/http"
	"strconv"

	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/service"
	"allandeng.cn/allandeng/go-blog/util"
	"github.com/gorilla/mux"
)

// path:"/"
func IndexHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	//get types
	ptpye := repository.NewPage(1, 6)
	types, err := repository.GetTypeRepository().FindTop(&ptpye)
	ctx.Model["types"] = types
	if err != nil {
		ctx.AddError(r, "Error cant find type, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}
	//get types
	ptag := repository.NewPage(1, 10)
	tags, err := repository.GetTagRepository().FindTop(&ptag)
	ctx.Model["tags"] = tags
	if err != nil {
		ctx.AddError(r, "Error cant find tag, err:%s", err)
		ctx.Next(ErrorHandler)
		return
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
	ctx.Model["blogs"] = blogs
	if err != nil {
		ctx.AddError(r, "Error cant find blog, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}

	//get recommand blogs
	precommands := repository.NewPage(1, 8)
	recommands, err := repository.GetBlogRepository().FindRecommendTop(&precommands)
	ctx.Model["recommands"] = recommands
	if err != nil {
		ctx.AddError(r, "Error cant find recommand, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}

	pblog.Update()
	ctx.Model["page"] = pblog
	ctx.Model["pagetitle"] = "首页"
	ctx.Model["active"] = 1
	RenderTemplate(ctx, w, r, "views/index.html", "views/_fragments.html")
}

// path:"/footer/newblog"
func NewBlogHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	//get recommand blogs
	pnewblogs := repository.NewPage(1, 3)
	newblogs, err := repository.GetBlogRepository().FindRecommendTop(&pnewblogs)
	ctx.Model["newblogs"] = newblogs
	if err != nil {
		ctx.AddError(r, "Error cant get footer, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}
	RenderTemplate(ctx, w, r, "views/_fragments-newblog.html")
}

// path:"/search"
func SearchHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	form := NewPostForm(r)
	query := form.GetString("query")
	log.Debugf("Received the form: %v,Search blog with : %s", r.PostForm, query)

	ctx.Model["query"] = query

	//search blogs
	pblog := repository.NewPage(1, 1000)
	blogs, err := repository.GetBlogRepository().FindByQuery(query, &pblog)
	pblog.Update()
	ctx.Model["blogs"] = blogs
	ctx.Model["page"] = pblog
	if err != nil {
		ctx.AddError(r, "Error cant find blog by query, query:%s err:%s", query, err)
		ctx.Next(ErrorHandler)
		return
	}

	ctx.Model["pagetitle"] = "搜索结果:" + query
	ctx.Model["active"] = 0
	RenderTemplate(ctx, w, r, "views/search.html", "views/_fragments.html")
}

// path: /blog/{id}
func BlogHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {

	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		ctx.AddError(r, "Error blog format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}

	blog, err := service.GetAndConvertBlog(int64(id))
	if err != nil {
		ctx.AddError(r, "Error blog not found, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}
	ctx.Model["blog"] = blog

	comments, err := service.ListCommentByBlogId(int64(id))
	if err != nil {
		ctx.AddError(r, "Error can't find blog comments, blogid: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}
	ctx.Model["comments"] = comments

	str, image := util.GetCaptcha(4)
	ctx.Session.Values["captcha"] = str
	ctx.Session.Save(r, w)

	ctx.Model["captcha"] = image

	if _, ok := ctx.Session.Values["user"].(blogmodel.User); ok {
		ctx.Model["admin"] = true
	}

	ctx.Model["pagetitle"] = blog.Title
	ctx.Model["active"] = 0

	RenderTemplate(ctx, w, r, "views/blog.html", "views/_fragments.html")
}

// path:"/captcha"
func CaptchaHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	str, image := util.GetCaptcha(4)
	ctx.Session.Values["captcha"] = str
	ctx.Session.Save(r, w)

	ctx.Model["captcha"] = image

	RenderTemplate(ctx, w, r, "views/_fragments-captcha.html")
}

// path: "/about"
func AboutHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {

	ctx.Model["pagetitle"] = "关于我"
	ctx.Model["active"] = 5
	RenderTemplate(ctx, w, r, "views/about.html", "views/_fragments.html", "views/_fragments-aboutme.html")
}

func ErrorHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	ctx.Model["active"] = 0
	ctx.Model["pagetitle"] = "500-服务器错误"
	RenderTemplate(ctx, w, r, "views/error/500.html", "views/_fragments.html")
}

func NotFoundHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	ctx.Model["active"] = 0
	ctx.Model["pagetitle"] = "404-找不到网页"
	RenderTemplate(ctx, w, r, "views/error/404.html", "views/_fragments.html")
}
