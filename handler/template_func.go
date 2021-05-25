package handler

import (
	"encoding/gob"
	"html/template"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/model"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	templateFunc map[string]interface{}
	log          = config.Logger
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = NewCookieStore(key)
)

func init() {
	templateFunc = map[string]interface{}{
		"htmltext": func(x string) interface{} { return template.HTML(x) },
	}
	gob.Register(model.User{})
}

func NewCookieStore(keyPairs ...[]byte) *sessions.CookieStore {
	cs := &sessions.CookieStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 0,
		},
	}
	cs.MaxAge(cs.Options.MaxAge)
	return cs
}
