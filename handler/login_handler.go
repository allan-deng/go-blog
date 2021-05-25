package handler

import (
	"net/http"
	"strings"

	"allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/util"
)

// path:"/admin"
func LoginPageHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if _, ok := ctx.Model["user"].(model.User); ok {
		http.Redirect(w, r, "/admin/index", http.StatusTemporaryRedirect)
		return
	}
	str, image := util.GetCaptcha(4)
	ctx.Session.Values["captcha"] = str
	ctx.Session.Save(r, w)

	ctx.Model["captcha"] = image
	ctx.Model["pagetitle"] = "登录"
	RenderTemplate(ctx, w, r, "views/admin/login.html", "views/_fragments.html")
}

// path:"/admin/login"
func LoginHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ctx.Next(LoginPageHandler)
		return
	}
	err := r.ParseForm()
	if err != nil {
		ctx.AddError(r, "parse form error, err:%s", err)
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	captcha := r.PostFormValue("captchacode")
	log.Debugf("Received the form: %v, Comment username: %s, password: %s, captcha: %s", r.PostForm, username, password, captcha)

	if captchaSession, ok := ctx.Session.Values["captcha"]; !ok || strings.ToLower(captchaSession.(string)) != strings.ToLower(captcha) {
		log.Infof("Login captcha failed, captcha: %s, host: %s", captcha, r.Host)
		ctx.Model["notice"] = "验证码错误"
		ctx.Next(LoginPageHandler)
		return
	}

	user, err := repository.GetUserRepository().FindUserByNameAndPassword(username, password)
	if err != nil || user == nil {
		log.Debugf("Cant find user by username and password, username: %s,password: %s, err: %s ", username, password, err)
		ctx.Model["notice"] = "用户名或密码错误"
		ctx.Next(LoginPageHandler)
		return
	}
	user.Password = ""
	ctx.Session.Values["user"] = user
	ctx.Session.Save(r, w)
	http.Redirect(w, r, "/admin/index", http.StatusTemporaryRedirect)
}

// path:"/logout"
func LogoutHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	delete(ctx.Session.Values, "user")
	ctx.Session.Save(r, w)
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
}

// path:"/admin/index"
func AdminIndexHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	ctx.Model["pagetitle"] = "后台管理"
	ctx.Model["active"] = 0
	ctx.Model["secondactive"] = 0
	RenderTemplate(ctx, w, r, "views/admin/index.html", "views/admin/_fragments.html")
}
