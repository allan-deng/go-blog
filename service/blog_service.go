package service

import (
	"bytes"
	"sort"

	"allandeng.cn/allandeng/go-blog/config"
	"allandeng.cn/allandeng/go-blog/model"
	"allandeng.cn/allandeng/go-blog/repository"
	"github.com/PuerkitoBio/goquery"
	"github.com/op/go-logging"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var log *logging.Logger

func init() {
	log = config.Logger
}

func GetAndConvertBlog(blogid int64) (*model.Blog, error) {
	blog, err := GetBlogIncreaseViews(blogid)
	if err != nil {
		log.Errorf("Error blog is nil, id: ", blogid)
		return blog, err
	}

	htmlstr := MarkdownToHtml(blog.Content)
	blog.Content = htmlstr
	return blog, err
}

func GetBlogIncreaseViews(blogid int64) (*model.Blog, error) {
	rep := repository.GetBlogRepository()
	b, err := rep.FindBlogById(blogid)
	if err != nil {
		log.Errorf("Error Can't find blog, id: ", blogid)
		return b, err
	}
	b.Views = b.Views + 1
	err = rep.UpdateBlog(b)
	if err != nil {
		log.Errorf("Error Can't update blog views, id: ", blogid)
		return b, err
	}
	return b, err
}

func MarkdownToHtml(markdown string) string {
	//convert markdown to html
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		panic(err)
	}

	//html add attribute
	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		panic(err)
	}
	docAddAttr(doc, "h1", "class", "ui header")
	docAddAttr(doc, "h2", "class", "ui header")
	docAddAttr(doc, "h3", "class", "ui header")
	docAddAttr(doc, "h4", "class", "ui header")
	docAddAttr(doc, "h5", "class", "ui header")
	docAddAttr(doc, "h6", "class", "ui header")
	docAddAttr(doc, "a", "target", "_blank")
	docAddAttr(doc, "img", "class", "ui rounded image m-shadow-small")
	docAddAttr(doc, "table", "class", "ui celled table")
	docAddAttr(doc, "code", "class", "token operator language- ")

	res, err := doc.Html()
	if err != nil {
		log.Errorf("convert doc to html error! err:%s", err)
	}
	return res
}

func docAddAttr(doc *goquery.Document, tag string, attr string, value string) {
	doc.Find(tag).Each(func(index int, ele *goquery.Selection) {
		newValue := ele.AttrOr(attr, "")
		newValue = newValue + " " + value
		ele.SetAttr(attr, newValue)
	})
}

func ListCommentByBlogId(blogid int64) ([]model.Comment, error) {
	parentComments, err := repository.GetCommentRepository().FindParentCommentByBlogId(blogid)
	if err != nil {
		log.Errorf("Error can't find parent comment, blogid: %s", blogid)
		return parentComments, err
	}
	for i := 0; i < len(parentComments); i++ {
		// comments := make([]model.Comment, 0)
		// var comments []model.Comment
		comments, err := AddAllChildren(parentComments[i])
		if err != nil {
			log.Errorf("Error can't find parents' chilren  comment, parentCommentId: %s", parentComments[i].ID)
			return parentComments, err
		}
		parentComments[i].ReplyComments = append(parentComments[i].ReplyComments, comments...)
	}
	return parentComments, err
}
func AddAllChildren(parentComment model.Comment) ([]model.Comment, error) {
	childrenComments, err := repository.GetCommentRepository().FindChildren(parentComment.ID)
	if err != nil {
		return nil, err
	}
	comments := make([]model.Comment, 0)
	for i := 0; i < len(childrenComments); i++ {
		p, err := repository.GetCommentRepository().FindParentCommentByCommentId(childrenComments[i].ParentCommentID)
		childrenComments[i].ParentComment = &p
		comment, err := AddAllChildren(childrenComments[i])
		comments = append(comments, childrenComments[i])
		comments = append(comments, comment...)
		if err != nil {
			log.Errorf("Error can't find parents' chilren  comment, parentCommentId: %s", childrenComments[i].ID)
			return comments, err
		}
	}
	return comments, nil
}

type Archive struct {
	Year  string
	Blogs []model.Blog
}

func ArchiveBlog() ([]Archive, error) {
	var res []Archive
	rep := repository.GetBlogRepository()
	years, err := rep.FindGroupYear()
	if err != nil {
		log.Errorf("Error can't find group year, err: %s", err)
		return nil, err
	}
	sort.Sort(sort.Reverse(sort.StringSlice((years))))
	for i := 0; i < len(years); i++ {
		blogs, err := rep.FindByYear(years[i])
		if err != nil {
			log.Errorf("Error can't find blog by year, year: %s, err: %s", years[i], err)
			return nil, err
		}
		temp := &Archive{
			Year:  years[i],
			Blogs: blogs,
		}
		res = append(res, *temp)
	}
	return res, nil
}
