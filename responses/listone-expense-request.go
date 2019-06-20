package responses

import (
	"go-rest-api/types"
	"net/http"
)

type ListExpenseResponse struct {
	*types.Expense
}


func (ListExpenseResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func List1expense(exp *types.Expense) *ListExpenseResponse {
	resp := &ListExpenseResponse{Expense: exp}

	return resp
}
