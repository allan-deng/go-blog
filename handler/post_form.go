package handler

import (
	"net/http"
	"strconv"
	"strings"
)

type PostForm struct {
	r *http.Request
}

func NewPostForm(r *http.Request) *PostForm {
	form := &PostForm{
		r: r,
	}
	form.r.ParseForm()
	return form
}

func (s *PostForm) GetString(key string) string {
	return s.r.PostFormValue(key)
}

func (s *PostForm) GetInt(key string) int {
	str := s.r.PostFormValue(key)
	num, _ := strconv.Atoi(str)
	return num
}

func (s *PostForm) GetIntSlice(key string) []int {
	str := s.r.PostFormValue(key)
	strs := strings.Split(str, ",")
	var res []int
	for _, str := range strs {
		temp, _ := strconv.Atoi(str)
		res = append(res, temp)
	}
	return res
}

func (s *PostForm) GetBool(key string) bool {
	str := s.r.PostFormValue(key)
	return str == "true"
}

func (s *PostForm) GetBoolContain(key string) bool {
	str := s.r.PostFormValue(key)
	return str != ""
}
