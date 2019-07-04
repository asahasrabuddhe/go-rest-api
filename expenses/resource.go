package expenses

import (
	"context"
	"errors"
	"github.com/asahasrabuddhe/rest-api/logger"
	"github.com/asahasrabuddhe/rest-api/renderers"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type Resource struct{}

func (e Resource) Routes() chi.Router {
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

func (e Resource) Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := chi.URLParam(r, "id"); id != "" {
			if idInt, err := strconv.Atoi(id); err != nil {
				_ = render.Render(w, r, &renderers.ErrorResponse{
					Err:            err,
					HTTPStatusCode: 422,
					StatusText:     "unable to parse id",
				})
			} else {
				logger.LogEntrySetField(r, "id", idInt)
				for _, expense := range Exp {
					if expense.Id == idInt {
						ctx := context.WithValue(r.Context(), "expense", expense)
						next.ServeHTTP(w, r.WithContext(ctx))

						return
					}
				}

				_ = render.Render(w, r, &renderers.ErrorResponse{
					Err:            errors.New("resource not found"),
					HTTPStatusCode: 404,
					StatusText:     "resource not found",
				})
			}
		}
	})
}

func (e Resource) Create(writer http.ResponseWriter, request *http.Request) {
	var req Create

	err := render.Bind(request, &req)
	if err != nil {
		logger.LogEntrySetField(request, "error", err)
		return
	}

	req.Expense.Id = len(Exp) + 1
	Exp = append(Exp, *req.Expense)

	_ = render.Render(writer, request, NewExpenseResponse(req.Expense))
}

func (e Resource) GetOne(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(Expense)
	_ = render.Render(writer, request, NewExpenseResponse(&expense))
}

func (e Resource) GetAll(writer http.ResponseWriter, request *http.Request) {
	_ = render.Render(writer, request, NewExpensesResponse(&Exp))
}

func (e Resource) Update(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(Expense)

	var req Update

	err := render.Bind(request, &req)
	if err != nil {
		logger.LogEntrySetField(request, "error", err)
		return
	}

	Exp[expense.Id-1] = *req.Expense

	_ = render.Render(writer, request, NewExpensesResponse(&Exp))
}

func (e Resource) Delete(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(Expense)

	Exp = append(Exp[:expense.Id], Exp[expense.Id+1:]...)
}
