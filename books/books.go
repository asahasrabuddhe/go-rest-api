package books

type Book struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Author string  `json:"author"`
	ISBN   string  `json:"isbn"`
	Price  float64 `json:"price"`
}

type Books []Book
