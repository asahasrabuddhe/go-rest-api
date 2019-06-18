package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go-rest-api/errors"
	"go-rest-api/requests"
	"go-rest-api/responses"
	"go-rest-api/types"
	"log"
	"net/http"
	"strconv"
)


var expenses types.Expenses

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", CreateExpense)
		r.Get("/", ListAllExpense)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(ExpenseContext)
			r.Get("/", ListOneExpense)
			r.Put("/", UpdateExpense)
			r.Delete("/", DeleteExpense)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func CreateExpense(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest

	err := render.Bind(request, &req)
	if err != nil {
		render.Render(writer, request, errors.ErrInvalidRequest(err))
		return
	}

	expenses = append(expenses, *req.Expense)

	j, _ := json.Marshal(req.Expense)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, _ = fmt.Fprintf(writer, `{"success": true, "data": %v}`, string(j))
}



func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(types.Expense)
	if err := render.Render(writer, request, responses.Listexpense(&expense)) ; err != nil{
		render.Render(writer,request,errors.ErrRender(err))
		return
	}
}

func ListAllExpense(writer http.ResponseWriter, request *http.Request) {
	if err := render.Render(writer, request, responses.ExpensesResponse(&expenses)); err != nil{
		render.Render(writer,request,errors.ErrRender(err))
		return
	}
}

func ExpenseContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		expenseID := chi.URLParam(r, "id")
		expid,_:=strconv.Atoi(expenseID)
		for _, expense := range expenses {

			if expense.Id == expid {
				ctx := context.WithValue(r.Context(), "expense", expense )
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}

	})
}


func UpdateExpense(writer http.ResponseWriter, request *http.Request) {

	expense := request.Context().Value("expense").(types.Expense)

	var req requests.UpdateExpenseRequest

	err := render.Bind(request, &req)
	if err != nil {
		render.Render(writer,request,errors.ErrRender(err))
		return
	}
			expenses[expense.Id-1] = *req.Expense


	if err = render.Render(writer, request, responses.Listexpense(&expense)) ; err != nil{
			render.Render(writer,request,errors.ErrRender(err))
			return

	}
}


func DeleteExpense(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(types.Expense)
	expenses = append(expenses[:expense.Id-1], expenses[expense.Id:]...)
	if err := render.Render(writer, request, responses.ExpensesResponse(&expenses)); err != nil{
		render.Render(writer,request,errors.ErrRender(err))
		return
	}
}