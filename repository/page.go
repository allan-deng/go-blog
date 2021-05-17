package repository

type Page struct {
	Index int
	Size  int
	Count int
}

func NewPage(index, size int) Page {
	return Page{
		Index: index,
		Size:  size,
	}
}
