package books

import (
    "errors"
    "net/http"
)

type Create struct {
    *Book
}

func (c *Create) Bind(r *http.Request) error {
    // TODO - Add Validations

    return nil
}

type Update struct {
    *Create
}

func (u *Update) Bind(r *http.Request) error {
    if u.Id == 0 {
        return errors.New("id is empty or is invalid")
    }

    return u.Create.Bind(r)
}