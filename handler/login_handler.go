package handler

import (
	"html/template"
	"net/http"
	"path"
	"strings"

	"allandeng.cn/allandeng/go-blog/config"
	blogmodel "allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/util"
	"github.com/Masterminds/sprig"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})

	str, image := util.GetCaptcha(4)
	session, _ := store.Get(r, "cookie-name")
	session.Values["captcha"] = str
	session.Save(r, w)

	model["captcha"] = image
	model["pagetitle"] = "登录"
	base := path.Base("views/admin/login.html")
	err := template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/admin/login.html", "views/_fragments.html")).Execute(w, model)
	if err != nil {
		panic(err)
	}
}

// path:"/login"
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()

	if r.Method != "POST" {
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
		return
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	captcha := r.PostFormValue("captchacode")
	log.Debugf("Received the form: %v, Comment username: %s, password: %s, captcha: %s", r.PostForm, username, password, captcha)

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		log.Errorf("Error get session failed, host: %s,cookie-name： %s, err: %s", r.Host, r.Header.Get("cookie-name"), err)
		panic(err)
	}

	if captchaSession, ok := session.Values["captcha"]; !ok || strings.ToLower(captchaSession.(string)) != strings.ToLower(captcha) {
		log.Errorf("Error login captcha failed, captcha: %s, host: %s", captcha, r.Host)
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
		return
	}

	user, err := repository.GetUserRepository().FindUserByNameAndPassword(username, password)
	if err != nil || user == nil {
		log.Errorf("Error find user by username and password, username: %s,password: %s, err: %s ", username, password, err)
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
		return
	}

	user.Password = ""
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		log.Errorf("Error save session. err：%s, user: %v", err, user)
	}
	for key, value := range session.Values {
		log.Debugf("key: %v , value: %v", key, value)
	}
	http.Redirect(w, r, "/admin/index", http.StatusTemporaryRedirect)
}

// path:"/logout"
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	delete(session.Values, "user")
	session.Save(r, w)
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
}

// path:"/admin/index"
func AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in index_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})
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
	model["pagetitle"] = "后台管理"
	model["active"] = 0
	model["secondactive"] = 0
	model["massage"] = config.GlobalMassage
	base := path.Base("views/admin/index.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/admin/index.html", "views/admin/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}
