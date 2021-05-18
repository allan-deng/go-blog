package repository

type Page struct {
	Index int  //当前页码
	Size  int  //每页数量
	Count int  //记录总数
	Nums  int  //总页数
	First bool //是否为首页
	Last  bool //是否为尾页
}

func NewPage(index, size int) Page {
	return Page{
		Index: index,
		Size:  size,
	}
}

//查询之后根据count更新First Last
func (s *Page) Update() {
	s.Nums = s.Count / s.Size
	if s.Count%s.Size > 0 {
		s.Nums++
	}
	s.First = false
	s.Last = false
	if s.Index == 1 {
		s.First = true
		return
	}
	if s.Index == s.Nums {
		s.Last = true
		return
	}
}
