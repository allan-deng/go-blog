package handler

import (
	"net/http"
	"strconv"

	"allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"github.com/gorilla/mux"
)

// path:"/admin/types"
func AdminTypesHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		AdminTypesCreateHandler(ctx, w, r)
	} else {
		AdminTypesListHandler(ctx, w, r)
	}
}

// path:"/admin/types" method:post
func AdminTypesCreateHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	form := NewPostForm(r)
	id := form.GetInt("id")
	name := form.GetString("name")
	blogtype := &model.Type{
		ID:   int64(id),
		Name: name,
	}
	var err error
	if id == 0 {
		_, err = repository.GetTypeRepository().CreateType(blogtype)
	} else {
		err = repository.GetTypeRepository().UpdateType(blogtype)
	}
	if err != nil {
		ctx.AddError(r, "create/update type failed, id:%d", id)
		ctx.Model["notice"] = "操作失败"

	} else {
		ctx.Model["notice"] = "操作成功"
	}
	ctx.Next(AdminTypesListHandler)
}

func AdminTypesListHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	types, err := repository.GetTypeRepository().FindAll()
	if err != nil {
		ctx.AddError(r, "cant find all types, err: %s", err)
	}
	ctx.Model["types"] = types
	ctx.Model["pagetitle"] = "分类列表"
	ctx.Model["active"] = 2
	RenderTemplate(ctx, w, r, "views/admin/types.html", "views/admin/_fragments.html")
}

// path:"/admin/types/input" method:get
func AdminTypesInputHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	ctx.Model["type"] = &model.Type{}
	ctx.Model["pagetitle"] = "创建分类"
	ctx.Model["active"] = 2
	RenderTemplate(ctx, w, r, "views/admin/type-input.html", "views/admin/_fragments.html")
}

// path:"/admin/types/{id:[0-9]+}/delete"
func AdminTypesDeleteHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		log.Errorf("Error type format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}
	err = repository.GetTypeRepository().DeleteType(int64(id))
	if err != nil {
		ctx.AddError(r, "cant delete type, typeid: %d ,err:%s", id, err)
		ctx.Model["notice"] = "删除失败"
	} else {
		ctx.Model["notice"] = "删除成功"
	}
	ctx.Model["pagetitle"] = "分类列表"
	ctx.Next(AdminTypesListHandler)
}

// path:"/admin/types/{id:[0-9]+}/input"
func AdminTypesUpdateHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		log.Errorf("Error type format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}
	blogtype, err := repository.GetTypeRepository().FindTypeById(int64(id))
	if err != nil {
		ctx.AddError(r, "cant find type by id: %d", id)
	}
	ctx.Model["type"] = blogtype
	ctx.Model["pagetitle"] = "更新分类"
	ctx.Model["active"] = 2
	RenderTemplate(ctx, w, r, "views/admin/type-input.html", "views/admin/_fragments.html")
}

// path:"/admin/types"
func AdminTagsHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		AdminTagsCreateHandler(ctx, w, r)
	} else {
		AdminTagsListHandler(ctx, w, r)
	}
}

// path:"/admin/tags" method:post
func AdminTagsCreateHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	form := NewPostForm(r)
	id := form.GetInt("id")
	name := form.GetString("name")
	blogtag := &model.Tag{
		ID:   int64(id),
		Name: name,
	}
	var err error
	if id == 0 {
		_, err = repository.GetTagRepository().CreateTag(blogtag)
	} else {
		err = repository.GetTagRepository().UpdateTag(blogtag)
	}
	if err != nil {
		ctx.AddError(r, "create/update tag failed, id:%d", id)
		ctx.Model["notice"] = "操作失败"

	} else {
		ctx.Model["notice"] = "操作成功"
	}
	ctx.Next(AdminTagsListHandler)
}

func AdminTagsListHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	tags, err := repository.GetTagRepository().FindAll()
	if err != nil {
		ctx.AddError(r, "cant find all tags, err: %s", err)
	}
	ctx.Model["tags"] = tags
	ctx.Model["pagetitle"] = "分类列表"
	ctx.Model["active"] = 3
	RenderTemplate(ctx, w, r, "views/admin/tags.html", "views/admin/_fragments.html")
}

// path:"/admin/tags/input" method:get
func AdminTagsInputHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	ctx.Model["tag"] = &model.Tag{}
	ctx.Model["pagetitle"] = "创建分类"
	ctx.Model["active"] = 3
	RenderTemplate(ctx, w, r, "views/admin/tag-input.html", "views/admin/_fragments.html")
}

// path:"/admin/tags/{id:[0-9]+}/delete"
func AdminTagsDeleteHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		log.Errorf("Error tag format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}
	err = repository.GetTagRepository().DeleteTag(int64(id))
	if err != nil {
		ctx.AddError(r, "cant delete tag, tagid: %d ,err:%s", id, err)
		ctx.Model["notice"] = "删除失败"
	} else {
		ctx.Model["notice"] = "删除成功"
	}
	ctx.Model["pagetitle"] = "分类列表"
	ctx.Next(AdminTagsListHandler)
}

// path:"/admin/tags/{id:[0-9]+}/input"
func AdminTagsUpdateHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		log.Errorf("Error tag format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}
	blogtag, err := repository.GetTagRepository().FindTagById(int64(id))
	if err != nil {
		ctx.AddError(r, "cant find tag by id: %d", id)
	}
	ctx.Model["tag"] = blogtag
	ctx.Model["pagetitle"] = "更新分类"
	ctx.Model["active"] = 3
	RenderTemplate(ctx, w, r, "views/admin/tag-input.html", "views/admin/_fragments.html")
}
