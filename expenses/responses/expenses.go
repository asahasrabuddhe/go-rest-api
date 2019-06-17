package responses

import (
	"github.com/asahasrabuddhe/rest-api/expenses"
	"net/http"
)

type ExpensesResponse struct {
	Expenses *expenses.Expenses `json:"data"`
	Success  bool               `json:"success"`
}

func NewExpensesResponse(expenses *expenses.Expenses) *ExpensesResponse {
	return &ExpensesResponse{Expenses: expenses}
}

func (e *ExpensesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if e.Expenses != nil {
		e.Success = true
	}
	return nil
}
