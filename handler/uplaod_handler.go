package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"allandeng.cn/allandeng/go-blog/config"
)

func UploadHandler(ctx *Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var filepath string
	failMsg, _ := json.Marshal(UploadResult{Success: 0, Message: "上传失败"})
	r.ParseMultipartForm(32 << 20) //设置临时文件大小

	mulitpartfile, fileHeader, err := r.FormFile("editormd-image-file") //mulitpart中的文件句柄
	if err != nil {
		ctx.AddError(r, "uploaf file failed! cant get mulitpartfile. err: %s", err)
		w.Write(failMsg)
		return
	}
	defer mulitpartfile.Close()

	exist, err := PathExists(config.GlobalConfig.WebServer.UploadPath)
	if err != nil {
		ctx.AddError(r, "uploaf file failed! cant get file path.path:%s , err: %s", config.GlobalConfig.WebServer.UploadPath, err)
		w.Write(failMsg)
		return
	}
	if !exist {
		err := os.Mkdir(config.GlobalConfig.WebServer.UploadPath, os.ModePerm)
		if err != nil {
			ctx.AddError(r, "uploaf file failed! create path failed.path:%s , err: %s", config.GlobalConfig.WebServer.UploadPath, err)
			w.Write(failMsg)
			return
		}
	}

	filepath = config.GlobalConfig.WebServer.UploadPath + fileHeader.Filename

	osfile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		ctx.AddError(r, "uploaf file failed! cant get osfile. err: %s", err)
		w.Write(failMsg)
		return
	}
	defer osfile.Close()

	_, err = io.Copy(osfile, mulitpartfile)
	if err != nil {
		ctx.AddError(r, "uploaf file failed! cant write osfile. err: %s", err)
		w.Write(failMsg)
		return
	}

	successMsg, _ := json.Marshal(UploadResult{Success: 1, Message: "上传成功", Url: strings.TrimPrefix(filepath, ".")})
	w.Write(successMsg)
}

type UploadResult struct {
	Success int    `json:"success"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
