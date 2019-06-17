package requests

import (
	"errors"
	"go-rest-api/types"
	"net/http"
)

type CreateExpenseRequest struct {
	*types.Expense
}

func (c *CreateExpenseRequest) Bind(r *http.Request) error {
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
