package handler

import (
	"net/http"
	"strconv"

	"allandeng.cn/allandeng/go-blog/repository"
	"github.com/gorilla/mux"
)

// path:"/types/{id}"
func TypeHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	//获取id
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id == 0 || id < -1 {
		ctx.AddError(r, "Error type format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
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
		ctx.AddError(r, "Error cant find all types. err: %s", err)
		ctx.Next(ErrorHandler)
		return
	}
	if id == -1 {
		id = int(types[0].ID)
	}
	ctx.Model["types"] = types
	//获取对应type的所有博客
	pblog := repository.NewPage(pagenum, 8)
	blogs, err := repository.GetBlogRepository().FindBlogByTypeId(int64(id), &pblog)
	pblog.Update()
	if pblog.Nums < pblog.Index {
		pblog.Index = 1
		blogs, err = repository.GetBlogRepository().FindBlogByTypeId(int64(id), &pblog)
	}
	ctx.Model["blogs"] = blogs
	if err != nil {
		ctx.AddError(r, "err: %s", err)
		ctx.Next(ErrorHandler)
		return
	}
	pblog.Update()
	ctx.Model["page"] = pblog
	ctx.Model["pagetitle"] = "分类"
	ctx.Model["active"] = 2
	ctx.Model["activetypeid"] = id
	RenderTemplate(ctx, w, r, "views/types.html", "views/_fragments.html")
}

// path:"/tags/{id}"
func TagHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	//获取id
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id == 0 || id < -1 {
		ctx.AddError(r, "Error type format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
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
		ctx.AddError(r, "Error cant find all tags. err: %s", err)
		ctx.Next(ErrorHandler)
		return
	}
	if id == -1 {
		id = int(tags[0].ID)
	}
	ctx.Model["tags"] = tags
	//获取对应tag的所有博客
	pblog := repository.NewPage(pagenum, 8)
	blogs, err := repository.GetBlogRepository().FindBlogByTagId(int64(id), &pblog)
	pblog.Update()
	if pblog.Nums < pblog.Index {
		pblog.Index = 1
		blogs, err = repository.GetBlogRepository().FindBlogByTagId(int64(id), &pblog)
	}
	ctx.Model["blogs"] = blogs
	if err != nil {
		ctx.AddError(r, "err: %s", err)
		ctx.Next(ErrorHandler)
		return
	}

	pblog.Update()
	ctx.Model["page"] = pblog
	ctx.Model["pagetitle"] = "标签"
	ctx.Model["active"] = 3
	ctx.Model["activetypeid"] = id
	RenderTemplate(ctx, w, r, "views/tags.html", "views/_fragments.html")
}
