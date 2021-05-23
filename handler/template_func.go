package handler

import (
	"encoding/gob"
	"html/template"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/model"
	"github.com/op/go-logging"
)

var templateFunc map[string]interface{}
var log *logging.Logger

func init() {
	log = config.Logger
	templateFunc = map[string]interface{}{
		"htmltext": func(x string) interface{} { return template.HTML(x) },
	}
	gob.Register(model.User{})
}
