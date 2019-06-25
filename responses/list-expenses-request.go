package responses

import (
	"go-rest-api/types"
	"net/http"
)

type ListExpensesResponse struct {
	Expenses []*types.Expense
}

func ExpensesResponse(expenses []*types.Expense) *ListExpensesResponse {
	resp := &ListExpensesResponse{Expenses: expenses}
	return resp
}

func (ListExpensesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

