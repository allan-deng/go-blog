package handler

import (
	"html/template"
	"net/http"
	"path"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/service"
	"github.com/Masterminds/sprig"
)

func ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Error in archive_handler : %s", err)
			ErrorHandler(w, r)
		}
	}()
	model := make(map[string]interface{})

	count, err := repository.GetBlogRepository().Count()
	if err != nil {
		panic(err)
	}
	model["count"] = count

	archive, err := service.ArchiveBlog()
	if err != nil {
		panic(err)
	}
	model["archive"] = archive

	model["pagetitle"] = "归档"
	model["active"] = 4
	model["massage"] = config.GlobalMassage
	base := path.Base("views/archives.html")
	err = template.Must(template.New(base).Funcs(sprig.FuncMap()).Funcs(templateFunc).ParseFiles("views/archives.html", "views/_fragments.html")).Execute(w, model)

	if err != nil {
		panic(err)
	}
}
