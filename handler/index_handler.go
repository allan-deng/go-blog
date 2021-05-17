package handler

import (
	"html/template"
	"net/http"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/repository"
	"github.com/op/go-logging"
)

var log *logging.Logger

func init() {
	log = config.Logger
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	model := make(map[string]interface{})

	//get types
	ptpye := repository.NewPage(1, 6)
	types, err := repository.GetTypeRepository().FindTop(&ptpye)
	model["types"] = types
	log.Debugf("Find types: %v", types)
	if err != nil {
		log.Errorf("Error repository : ", err)
		ErrorHandler(w, r)
		return
	}
	//get types
	ptag := repository.NewPage(1, 10)
	tags, err := repository.GetTagRepository().FindTop(&ptag)
	model["tags"] = tags
	log.Debugf("Find tags: %v", tags)
	if err != nil {
		log.Errorf("Error repository : ", err)
		ErrorHandler(w, r)
		return
	}
	//get blogs
	pblog := repository.NewPage(1, 10)
	blogs, err := repository.GetBlogRepository().FindAll(&pblog)
	model["blogs"] = blogs
	log.Debugf("Find blogs: %v", blogs)
	if err != nil {
		log.Errorf("Error repository : ", err)
		ErrorHandler(w, r)
		return
	}

	//get recommand blogs
	precommands := repository.NewPage(1, 8)
	recommands, err := repository.GetBlogRepository().FindRecommendTop(&precommands)
	model["recommands"] = recommands
	log.Debugf("Find recommands: %v", recommands)
	if err != nil {
		log.Errorf("Error repository : ", err)
		ErrorHandler(w, r)
		return
	}
	model["page"] = pblog
	model["pagetitle"] = "首页"
	model["active"] = 1
	tmpl, err := template.ParseFiles("views/index.html", "views/_fragments.html")
	if err != nil {
		log.Errorf("Error Create template failed : ", err)
		ErrorHandler(w, r)
		return
	}
	tmpl.Execute(w, model)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {

}
