package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go-rest-api/errrs"
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
		var err error
		flag:=1
		expenseID := chi.URLParam(r, "id")
		b,_:=strconv.Atoi(expenseID)
		for _, a := range expenses {
			if a.Id == b {
				flag=0
				ctx := context.WithValue(r.Context(), "expense", a )
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
		if flag ==1{
			err=errors.New("ID not Found")
			render.Render(w, r, errrs.ErrRender(err))
		}
	})
}
func CreateExpense(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest
	var err error

	err = render.Bind(request, &req)
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
	expenses = append(expenses, *req.Expense)
	render.Render(writer, request, responses.List1expense(req.Expense))
}
func UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(types.Expense)
	var req requests.UpdateExpenseRequest

	err:= render.Bind(request,&req)
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
		for index, e := range expenses {
			if e.Id == expense.Id {
				expenses[index] = *req.Expense
			}
	}
		errs:=render.Render(writer, request, responses.List1expense(&expense))
	if errs != nil {
		render.Render(writer,request,errrs.ErrRender(errs))
		return
	}
}
func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
	exp := request.Context().Value("expense").(types.Expense)
	err:=render.Render(writer, request, responses.List1expense(&exp))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
}
func ListAllExpense(writer http.ResponseWriter, request *http.Request) {
	err := render.Render(writer, request, responses.NewExpensesResponse(&expenses))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
}
func DeleteExpense(writer http.ResponseWriter, request *http.Request) {

	exp := request.Context().Value("expense").(types.Expense)
	expenses =append(expenses[:exp.Id], expenses[exp.Id+1:]...)
	err:=render.Render(writer,request,responses.NewExpensesResponse(&expenses))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
}
