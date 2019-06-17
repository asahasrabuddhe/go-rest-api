package responses

import (
	"github.com/go-chi/render"
	"go-rest-api/types"
	"net/http"
)

type ExpensesResponse struct {
  Expenses *ListExpenseResponse
}

func NewExpensesResponse(expenses *types.Expenses) []render.Renderer {
	list := []render.Renderer{}
	for _, exp := range *expenses {
		list = append(list, List1expense(&exp))
	}
	return list

}

func (e *ExpensesResponse) Render(w http.ResponseWriter, r *http.Request) error {

	return nil
}