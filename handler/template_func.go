package handler

import (
	"html/template"

	"allandeng.cn/allandeng/go-blog/config"
	"github.com/op/go-logging"
)

var templateFunc map[string]interface{}
var log *logging.Logger

func init() {
	log = config.Logger
	templateFunc = map[string]interface{}{
		"htmltext": func(x string) interface{} { return template.HTML(x) },
	}
}
