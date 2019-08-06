package responses

import (
	"go-rest-api/types"
	"net/http"
)

type ExpensesResponse struct {
 *types.Expenses
}

func NewExpensesResponse(expenses *types.Expenses) *ExpensesResponse{
	return &ExpensesResponse{ expenses}

}

func (e *ExpensesResponse) Render(w http.ResponseWriter, r *http.Request) error {

	return nil
}