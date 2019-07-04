package expenses

import "net/http"

type Response struct {
	Expense *Expense `json:"data,omitempty"`
	Success bool     `json:"success"`
}

func NewExpenseResponse(expense *Expense) *Response {
	return &Response{Expense: expense}
}

func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if e.Expense != nil {
		e.Success = true
	}
	return nil
}

type ListResponse struct {
	Expenses *Expenses `json:"data,omitempty"`
	Success  bool      `json:"success"`
}

func NewExpensesResponse(expenses *Expenses) *ListResponse {
	return &ListResponse{Expenses: expenses}
}

func (e *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if e.Expenses != nil {
		e.Success = true
	}
	return nil
}
