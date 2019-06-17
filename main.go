package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		var exp *types.Expense
		var e error
		articleID := chi.URLParam(r, "id")
		exp,e = GetExpense(articleID)
		if e != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), "expense", exp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetExpense(id string) (*types.Expense, error) {
	b,_:=strconv.Atoi(id)
	for _, a := range expenses {
		if a.Id == b {
			return &a,nil
		}
	}
	return nil, errors.New("article not found.")
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
	expense := req.Expense
	render.Status(request, http.StatusCreated)
	render.Render(writer, request, responses.List1expense(expense))
}
func UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest

	err:= render.Bind(request,&req)
	if err != nil {
		log.Println(err)
		return
	}
	vars := req.Id
	var temp1 types.Expenses
	for index, articles := range expenses {
		if (articles.Id == vars) {
			temp1 = expenses[index+1:]
			expenses=append(expenses[:index],*req.Expense)
			expenses=append(expenses[:index],temp1...)
		}
	}
	//j, _ := json.Marshal(req.Expense)
	//writer.Header().Set("Content-Type", "application/json")
	//writer.WriteHeader(http.StatusCreated)
	//
	//_, _ = fmt.Fprintf(writer, `{"success": true, "data": %v}`, string(j))

	render.Render(writer, request, responses.List1expense(req.Expense))
}


func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
	exp := request.Context().Value("expense").(*types.Expense)
	vars:= chi.URLParam(request,"id")
	a,err:=strconv.Atoi(vars)
	flag:=0
	if err !=nil{
		http.Error(writer,"ID of the expense not parsed",500)
	}
	for _, article := range expenses {
		if (article.Id==a) {
			flag=1
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
		}
		if err := render.Render(writer, request,responses.List1expense(exp)); err != nil {
			render.Render(writer, request,errrs.ErrRender(err))
			return
		}
	}
	if flag==0{
		http.Error(writer,"expense with ID "+vars+" not found",500)
	}


}

func ListAllExpense(writer http.ResponseWriter, request *http.Request) {
	if err := render.RenderList(writer, request, responses.NewExpensesResponse(&expenses)); err != nil {
		render.Render(writer, request, errrs.ErrRender(err))
		return
	}
}


func DeleteExpense(writer http.ResponseWriter, request *http.Request) {
	vars:= chi.URLParam(request,"id")
	a,err:=strconv.Atoi(vars)
	if err !=nil{
		http.Error(writer,"ID of the expense not parsed",500)
	}
	for index, article := range expenses {
		if (article.Id==a) {
			expenses =append(expenses[:index], expenses[index+1:]...)
			j, _ := json.Marshal(article)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)

			_, _ = fmt.Fprintf(writer, `{"success": true, "data": %v}`, string(j))
		}
	}

}
