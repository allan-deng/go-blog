package handler

import (
	"net/http"
	"strconv"
	"strings"

	"allandeng.cn/allandeng/go-blog/config"
	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/service"
	"github.com/gorilla/mux"
)

// path: "/comments" post
func CommentCreateHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	form := NewPostForm(r)

	parentid := form.GetInt("parentComment.id")
	blogid := form.GetInt("blog.id")
	nickname := form.GetString("nickname")
	email := form.GetString("email")
	content := form.GetString("content")
	captcha := form.GetString("captchacode")
	log.Debugf("Received the form: %v, Comment blogid: %s, parentid: %s, content: %s,nickname: %s,email: %s ", r.PostForm, blogid, parentid, content, nickname, email)
	bid := blogid
	pid := parentid

	ctx.Model["commentmassage"] = ""
	if captchaSession, ok := ctx.Session.Values["captcha"]; !ok || strings.ToLower(captchaSession.(string)) != strings.ToLower(captcha) {
		log.Errorf("Error captcha failed, captcha: %s, host: %s", captcha, r.Host)
		ctx.Model["commentmassage"] = "验证码错误，请点击验证码刷新！"
		comments, err := service.ListCommentByBlogId(int64(bid))
		if err != nil {
			log.Errorf("Error can't find blog comments, blogid: %s ,err: %s", bid, err)
		}
		ctx.Model["comments"] = comments
		RenderTemplate(ctx, w, r, "views/_fragments-comment.html")
		return
	}

	admin := false
	avater := config.GlobalMassage.CommentAvatar
	if u, ok := ctx.Session.Values["user"].(blogmodel.User); ok {
		admin = true
		avater = u.Avatar
	}

	comment := &blogmodel.Comment{
		Nickname:        nickname,
		Email:           email,
		Content:         content,
		Avatar:          avater,
		BlogID:          int64(bid),
		ParentCommentID: int64(pid),
		AdminComment:    admin,
	}
	_, err := repository.GetCommentRepository().CreateComment(comment)
	if err != nil {
		ctx.AddError(r, "Error add comment failed,comment:%v , err： %s", comment, err)
		ctx.Next(ErrorHandler)
		return
	}

	comments, err := service.ListCommentByBlogId(int64(bid))
	if err != nil {
		ctx.AddError(r, "Error can't find blog comments, blogid: %s ,err: %s", bid, err)
	}
	ctx.Model["comments"] = comments

	RenderTemplate(ctx, w, r, "views/_fragments-comment.html")
}

// path: "/comments/delete/{id}"
func CommentDeleteHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		ctx.AddError(r, "Error commentid format error, id: %s ,err: %s", idstr, err)
		ctx.Next(ErrorHandler)
		return
	}

	comment, err := repository.GetCommentRepository().FindCommentById(int64(id))
	if err != nil {
		ctx.AddError(r, "Error Comments do not have this ID, id: %s,err： %s", id, err)
		ctx.Next(NotFoundHandler)
		return
	}
	blogid := strconv.Itoa(int(comment.BlogID))

	if _, ok := ctx.Session.Values["user"].(blogmodel.User); !ok {
		ctx.AddError(r, "Error no permission to delete comments, comment id: %s , host: %s", id, r.Host)
		http.Redirect(w, r, "/blog/"+blogid, http.StatusTemporaryRedirect)
		return
	}
	repository.GetCommentRepository().DeleteComment(comment.ID)
	http.Redirect(w, r, "/blog/"+blogid, http.StatusTemporaryRedirect)
}
