package handler

import (
	"net/http"

	"allandeng.cn/allandeng/go-blog/repository"
	"allandeng.cn/allandeng/go-blog/service"
)

func ArchiveHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	count, err := repository.GetBlogRepository().Count()
	if err != nil {
		ctx.AddError(r, "err: %s", err)
		ctx.Next(ErrorHandler)
		return
	}
	ctx.Model["count"] = count

	archive, err := service.ArchiveBlog()
	if err != nil {
		ctx.AddError(r, "err: %s", err)
		ctx.Next(ErrorHandler)
		return
	}
	ctx.Model["archive"] = archive

	ctx.Model["pagetitle"] = "归档"
	ctx.Model["active"] = 4
	RenderTemplate(ctx, w, r, "views/archives.html", "views/_fragments.html")
}
