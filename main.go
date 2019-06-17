package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
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
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", CreateExpense)
		r.Get("/", ListAllExpense)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(ArticleCtx)
			r.Get("/", ListOneExpense)
			r.Put("/", UpdateExpense)
			r.Delete("/", DeleteExpense)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	//	var e error
		expenseID := chi.URLParam(r, "id")
		b,_:=strconv.Atoi(expenseID)
		for _, a := range expenses {

			if a.Id == b {
				ctx := context.WithValue(r.Context(), "expense", a )
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}

	})
}


func CreateExpense(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest

	err := render.Bind(request, &req)
	if err != nil {
		log.Println(err)
		return
	}

	expenses = append(expenses, *req.Expense)

	//j, _ := json.Marshal(req.Expense)
	//writer.Header().Set("Content-Type", "application/json")
	//writer.WriteHeader(http.StatusCreated)
	//
	//_, _ = fmt.Fprintf(writer, `{"success": true, "data": %v}`, string(j))
	//render.Status(request, http.StatusCreated)
	render.Render(writer, request, responses.List1expense(req.Expense))
}
func UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(types.Expense)
	var req requests.UpdateExpenseRequest

	err:= render.Bind(request,&req)
	if err != nil {
		log.Println(err)
		return
	}
// for loop -> get index of expense
		for index, e := range expenses {
			if e.Id == expense.Id {
				expenses[index] = *req.Expense
			}
	}
	//j, _ := json.Marshal(req.Expense)
	//writer.Header().Set("Content-Type", "application/json")
	//writer.WriteHeader(http.StatusCreated)
	//
	//_, _ = fmt.Fprintf(writer, `{"success": true, "data": %v}`, string(j))

	render.Render(writer, request, responses.List1expense(&expense))
}


func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
	exp := request.Context().Value("expense").(types.Expense)
	_ = render.Render(writer, request, responses.List1expense(&exp))
}


func ListAllExpense(writer http.ResponseWriter, request *http.Request) {
	_ = render.Render(writer, request, responses.NewExpensesResponse(&expenses))

}


func DeleteExpense(writer http.ResponseWriter, request *http.Request) {

	exp := request.Context().Value("expense").(types.Expense)
	expenses =append(expenses[:exp.Id], expenses[exp.Id+1:]...)

}
