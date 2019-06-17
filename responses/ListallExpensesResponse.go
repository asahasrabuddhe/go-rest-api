package responses

import (
	"errors"
	"go-rest-api/types"
	"net/http"
)

type ExpensesResponse struct {
  Expenses *types.Expenses
}

func NewExpensesResponse(expenses *types.Expenses) *ExpensesResponse{
	return &ExpensesResponse{Expenses: expenses}

}

func (e *ExpensesResponse) Render(w http.ResponseWriter, r *http.Request) error {

	return errors.New("Error at listing  Expenses")
}