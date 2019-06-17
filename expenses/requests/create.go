package requests

import (
	"errors"
	"github.com/asahasrabuddhe/rest-api/expenses"
	"net/http"
)

type CreateExpenseRequest struct {
	*expenses.Expense
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
