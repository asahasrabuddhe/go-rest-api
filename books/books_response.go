package books

import (
    "net/http"
)

type Response struct {
    Book *Book `json:"data"`
    Success bool `json:"success"`
}

func NewBookResponse(book *Book) *Response {
    return &Response{ Book: book }
}

func (b *Response) Render(w http.ResponseWriter, r *http.Request) error {
    if b.Book != nil {
        b.Success = true
    }

    return nil
}

type ListResponse struct {
    Books *Books `json:"data"`
    Success bool `json:"success"`
}

func NewBooksResponse(books *Books) *ListResponse {
    return &ListResponse{ Books: books }
}

func (b *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
    if b.Books != nil {
        b.Success = true
    }

    return nil
}