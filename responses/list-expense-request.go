
package responses

import (
	"go-rest-api/expenseDB"
	"go-rest-api/types"
	"net/http"
)

type ListExpenseResponse struct {
	*types.Expense
}

var mh expenseDB.MongoHandler

func Listexpense(expense *types.Expense) *ListExpenseResponse {
	resp := &ListExpenseResponse{Expense: expense}

	return resp

}

func (ListExpenseResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}


