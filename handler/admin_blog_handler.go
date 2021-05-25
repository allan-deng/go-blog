package handler

import (
	"net/http"
	"strconv"

	"allandeng.cn/allandeng/go-blog/model"
	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"github.com/gorilla/mux"
)

// path:"/admin/blogs"
func AdminBlogListHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		AdminBlogCreateHandler(ctx, w, r)
	} else {
		AdminBlogListGetHandler(ctx, w, r)
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
 recommend: on
shareStatement: on
appreciation: on
commentabled: on
bool类型的为false的直接没有这个字段，为true的返回为on
*/
// path:"/admin/blogs" method:"post"
func AdminBlogCreateHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	form := NewPostForm(r)
	id := form.GetInt("id")
	typeid := form.GetInt("type.id")
	tagids := form.GetIntSlice("tagIds")
	title := form.GetString("title")
	content := form.GetString("content")
	flag := form.GetString("flag")
	firstPicture := form.GetString("firstPicture")
	description := form.GetString("description")
	published := form.GetBool("published")
	recommend := form.GetBoolContain("recommend")
	shareStatement := form.GetBoolContain("shareStatement")
	appreciation := form.GetBoolContain("appreciation")
	commentabled := form.GetBoolContain("commentabled")

	blogtype, err := repository.GetTypeRepository().FindTypeById(int64(typeid))
	if err != nil {
		ctx.AddError(r, "cant find type by id, id:%d , err:%s", typeid, err)
	}

	var tags []model.Tag
	for _, tagid := range tagids {
		tag, err := repository.GetTagRepository().FindTagById(int64(tagid))
		if err != nil {
			ctx.AddError(r, "cant find tag by id, id:%d , err:%s", tag.ID, err)
		} else {
			tags = append(tags, *tag)
		}
	}

	blog := &blogmodel.Blog{
		ID:             int64(id),
		Title:          title,
		Content:        content,
		FirstPicture:   firstPicture,
		Flag:           flag,
		Appreciation:   appreciation,
		ShareStatement: shareStatement,
		Commentabled:   commentabled,
		Published:      published,
		Recommend:      recommend,
		Type:           *blogtype,
		Tags:           tags,
		UserID:         0,
		User:           ctx.Model["user"].(model.User),
		Description:    description,
	}

	if blog.ID == 0 {
		_, err = repository.GetBlogRepository().CreateBlog(blog)
	} else {
		err = repository.GetBlogRepository().UpdateBlog(blog)
	}
	if err != nil {
		ctx.AddError(r, "cant update blog, blogid: %d ,err:%s", id, err)
		ctx.Model["notice"] = "操作失败"
	} else {
		ctx.Model["notice"] = "操作成功"
	}
	ctx.Next(AdminBlogListGetHandler)
}

// path:"/admin/blogs" method:"get"
func AdminBlogListGetHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	pblog := repository.NewPage(1, 100000)
	blogs, err := repository.GetBlogRepository().FindAll(&pblog)
	if err != nil {
		ctx.AddError(r, "Error cant find all blog, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}
	types, err := repository.GetTypeRepository().FindAll()
	if err != nil {
		ctx.AddError(r, "Error cant find all blog, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}
	ctx.Model["blogs"] = blogs
	ctx.Model["types"] = types
	ctx.Model["count"] = pblog.Count
	ctx.Model["pagetitle"] = "博客列表"
	ctx.Model["active"] = 1
	ctx.Model["secondactive"] = 0
	RenderTemplate(ctx, w, r, "views/admin/blogs.html", "views/admin/_fragments.html")
}

// path:"/admin/blogs/search" method:"post"
func AdminBlogSearchHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/admin/blogs", http.StatusTemporaryRedirect)
		return
	}
	err := r.ParseForm()
	if err != nil {
		ctx.AddError(r, "parse form error, err:%s", err)
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
		ctx.Next(ErrorHandler)
		return
	}

	ctx.Model["blogs"] = blogs
	RenderTemplate(ctx, w, r, "views/admin/_fragments-bloglist.html")
}

// path:"/admin/blogs/input"
func AdminBlogInputHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	types, err := repository.GetTypeRepository().FindAll()
	if err != nil {
		ctx.AddError(r, "Error cant find all type, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}

	tags, err := repository.GetTagRepository().FindAll()
	if err != nil {
		ctx.AddError(r, "Error cant find all tags, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}
	ctx.Model["blog"] = &blogmodel.Blog{}
	ctx.Model["types"] = types
	ctx.Model["tags"] = tags
	ctx.Model["pagetitle"] = "创建博客"
	ctx.Model["active"] = 1
	ctx.Model["secondactive"] = 1
	RenderTemplate(ctx, w, r, "views/admin/blogs-input.html", "views/admin/_fragments.html")
}

// path:"/blogs/{id:[0-9]+}/delete"
func AdminBlogUpdateHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		ctx.AddError(r, "Error blog format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}
	blog, err := repository.GetBlogRepository().FindBlogById(int64(id))
	if err != nil {
		ctx.AddError(r, "Error cant find blog,id: %d, err:%s", id, err)
		ctx.Next(ErrorHandler)
		return
	}
	tagids := ""
	for i, tag := range blog.Tags {
		if i == len(blog.Tags)-1 {
			tagids = tagids + strconv.Itoa(int(tag.ID))
		} else {
			tagids = tagids + strconv.Itoa(int(tag.ID)) + ","
		}
	}
	blog.TagIds = tagids

	types, err := repository.GetTypeRepository().FindAll()
	if err != nil {
		ctx.AddError(r, "Error cant find all type, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}

	tags, err := repository.GetTagRepository().FindAll()
	if err != nil {
		ctx.AddError(r, "Error cant find all tags, err:%s", err)
		ctx.Next(ErrorHandler)
		return
	}
	ctx.Model["types"] = types
	ctx.Model["tags"] = tags
	ctx.Model["blog"] = blog
	ctx.Model["pagetitle"] = "更新博客"
	ctx.Model["active"] = 1
	ctx.Model["secondactive"] = 1
	RenderTemplate(ctx, w, r, "views/admin/blogs-input.html", "views/admin/_fragments.html")
}

// path:"/blogs/{id:[0-9]+}/delete"
func AdminBlogDeleteHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		ctx.AddError(r, "Error blog format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}

	err = repository.GetBlogRepository().DeleteBlog(int64(id))
	if err != nil {
		ctx.AddError(r, "cant delete blog, blogid: %d ,err:%s", id, err)
		ctx.Model["notice"] = "删除失败"
	} else {
		ctx.Model["notice"] = "删除成功"
	}
	ctx.Model["pagetitle"] = "博客列表"
	http.Redirect(w, r, "/admin/blogs", http.StatusTemporaryRedirect)
}
