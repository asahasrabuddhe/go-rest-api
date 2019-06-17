package resource

import (
	"context"
	"errors"
	"github.com/asahasrabuddhe/rest-api/expenses"
	"github.com/asahasrabuddhe/rest-api/expenses/requests"
	"github.com/asahasrabuddhe/rest-api/expenses/responses"
	"github.com/asahasrabuddhe/rest-api/logger"
	"github.com/asahasrabuddhe/rest-api/renderers"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type ExpenseResource struct{}

func (e ExpenseResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", e.Create)
	r.Get("/", e.GetAll)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(e.Context)

		r.Get("/", e.GetOne)
		r.Put("/", e.Update)
		r.Delete("/", e.Delete)
	})

	return r
}

func (e ExpenseResource) Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := chi.URLParam(r, "id"); id != "" {
			if idInt, err := strconv.Atoi(id); err != nil {
				// error
			} else {
				for index, expense := range expenses.Exp {
					if index == idInt {
						ctx := context.WithValue(r.Context(), "expense", expense)
						next.ServeHTTP(w, r.WithContext(ctx))
					}
				}

				_ = render.Render(w, r, &renderers.ErrorResponse{
					Err: errors.New("resource not found"),
					HTTPStatusCode: 404,
					StatusText: "resource not found",
				})
			}
		}
	})
}

func (e ExpenseResource) Create(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest

	err := render.Bind(request, &req)
	if err != nil {
		logger.LogEntrySetField(request, "error", err)
		return
	}

	req.Expense.Id = len(expenses.Exp) + 1
	expenses.Exp = append(expenses.Exp, *req.Expense)

	_ = render.Render(writer, request, responses.NewExpenseResponse(req.Expense))
}

func (e ExpenseResource) GetOne(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(expenses.Expense)
	_ = render.Render(writer, request, responses.NewExpenseResponse(&expense))
}

func (e ExpenseResource) GetAll(writer http.ResponseWriter, request *http.Request) {
	_ = render.Render(writer, request, responses.NewExpensesResponse(&expenses.Exp))
}

func (e ExpenseResource) Update(writer http.ResponseWriter, request *http.Request) {

}

func (e ExpenseResource) Delete(writer http.ResponseWriter, request *http.Request) {

}
