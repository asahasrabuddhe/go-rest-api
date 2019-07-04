package expenses

import (
	"errors"
	"net/http"
)

type Create struct {
	*Expense
}

func (c *Create) Bind(r *http.Request) error {
	if c.Description == "" {
		return errors.New("description is either empty or invalid")
	}

	if c.Amount == 0 {
		return errors.New("amount is either empty or invalid")
	}

	if c.Type == "" {
		return errors.New("description is either empty or invalid")
	}

	return nil
}

type Update struct {
	*Create
}

func (u *Update) Bind(r *http.Request) error {
	if u.Id == 0 {
		return errors.New("id is empty or invalid")
	}

	return u.Create.Bind(r)
}