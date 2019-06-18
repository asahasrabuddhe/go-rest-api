package response

import (
	"go-rest-api/types"
	"net/http"
)

type ListExpensesResponse struct {
	Expenses *types.Expenses
}

func ExpensesResponse(expenses *types.Expenses) *ListExpensesResponse {
	resp := &ListExpensesResponse{Expenses: expenses}
	return resp
}

func (ListExpensesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
