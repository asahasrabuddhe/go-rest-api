package responses

import (
	"errors"
	"go-rest-api/types"
	"net/http"
)

type ListExpenseResponse struct {
	*types.Expense
}


func (ListExpenseResponse) Render(w http.ResponseWriter, r *http.Request) error {

	return errors.New("Error at listing Expense")
}
func List1expense(exp *types.Expense) *ListExpenseResponse {
	resp := &ListExpenseResponse{Expense: exp}

	return resp
}
