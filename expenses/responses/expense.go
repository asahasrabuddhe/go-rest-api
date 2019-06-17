package responses

import (
	"github.com/asahasrabuddhe/rest-api/expenses"
	"net/http"
)

type ExpenseResponse struct {
	Expense *expenses.Expense `json:"data"`
	Success bool              `json:"success"`
}

func NewExpenseResponse(expense *expenses.Expense) *ExpenseResponse {
	return &ExpenseResponse{Expense: expense}
}

func (e *ExpenseResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if e.Expense != nil {
		e.Success = true
	}
	return nil
}
