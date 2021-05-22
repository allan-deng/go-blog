package handler

import (
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"

	"allandeng.cn/allandeng/go-blog/config"
	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/service"
	"github.com/Masterminds/sprig"
	"github.com/gorilla/mux"
)

// path: "/comments" post
func CommentCreateHandler(w http.ResponseWriter, r *http.Request) {
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

	parentid := r.PostFormValue("parentComment.id")
	blogid := r.PostFormValue("blog.id")
	nickname := r.PostFormValue("nickname")
	email := r.PostFormValue("email")
	content := r.PostFormValue("content")
	captcha := r.PostFormValue("captchacode")
	log.Debugf("Received the form: %v, Comment blogid: %s, parentid: %s, content: %s,nickname: %s,email: %s ", r.PostForm, blogid, parentid, content, nickname, email)
	bid, _ := strconv.Atoi(blogid)
	pid, _ := strconv.Atoi(parentid)

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Errorf("Error get session failed, host: %s,cookie-name： %s, err: %s", r.Host, r.Header.Get("cookie-name"), err)
		panic(err)
	}

	model := make(map[string]interface{})
	model["commentmassage"] = ""
	if captchaSession, ok := session.Values["captcha"]; !ok || strings.ToLower(captchaSession.(string)) != strings.ToLower(captcha) {
		log.Errorf("Error captcha failed, captcha: %s, host: %s", captcha, r.Host)
		model["commentmassage"] = "验证码错误，请点击验证码刷新！"
		comments, err := service.ListCommentByBlogId(int64(bid))
		if err != nil {
			log.Errorf("Error can't find blog comments, blogid: %s ,err: %s", bid, err)
		}
		model["comments"] = comments

		base := path.Base("views/_fragments-comment.html")
		err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/_fragments-comment.html")).Execute(w, model)

		if err != nil {
			panic(err)
		}
		return
	}

	admin := false
	avater := config.GlobalMassage.CommentAvatar
	if u, ok := session.Values["user"].(blogmodel.User); ok {
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
	_, err = repository.GetCommentRepository().CreateComment(comment)
	if err != nil {
		log.Errorf("Error add comment failed,comment:%v , err： %s", comment, err)
		panic(err)
	}

	comments, err := service.ListCommentByBlogId(int64(bid))
	if err != nil {
		log.Errorf("Error can't find blog comments, blogid: %s ,err: %s", bid, err)
	}
	model["comments"] = comments

	base := path.Base("views/_fragments-comment.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/_fragments-comment.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}

// path: "/comments/delete/{id}"
func CommentDeleteHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()

	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		log.Errorf("Error commentid format error, id: %s ,err: %s", idstr, err)
		panic(err)
	}

	comment, err := repository.GetCommentRepository().FindCommentById(int64(id))
	if err != nil {
		log.Errorf("Error Comments do not have this ID, id: %s,err： %s", id, err)
		NotFoundHandler(w, r)
		return
	}
	blogid := strconv.Itoa(int(comment.BlogID))

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Errorf("Error get session failed, host: %s,cookie-name： %s, err: %s", r.Host, r.Header.Get("cookie-name"), err)
		panic(err)
	}
	if _, ok := session.Values["user"].(blogmodel.User); !ok {
		log.Errorf("Error no permission to delete comments, comment id: %s , host: %s", id, r.Host)
		http.Redirect(w, r, "/blog/"+blogid, http.StatusTemporaryRedirect)
		return
	}
	repository.GetCommentRepository().DeleteComment(comment.ID)
	http.Redirect(w, r, "/blog/"+blogid, http.StatusTemporaryRedirect)
}
