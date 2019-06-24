package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go-rest-api/ers"
	"go-rest-api/requests"
	"go-rest-api/response"
	"go-rest-api/types"
	"log"
	"net/http"
	"strconv"
	"upper.io/db.v3"
	"upper.io/db.v3/mongo"
)

var settings = mongo.ConnectionURL{
	Database:  `expense_db`,
	Host:      `127.0.0.1`,
}
var sess db.Database
var expenseCollection db.Collection
var expenses types.Expenses

func main() {
	var err error
	sess, err = mongo.Open(settings)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	defer sess.Close()

	expenseCollection = sess.Collection("expense")

	err = expenseCollection.Truncate()
	if err != nil {
		log.Fatalf("Truncate(): %q\n", err)
	}

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
			r.Use(ExpenseContext)
			r.Get("/", ListOneExpense)
			r.Put("/", UpdateExpense)
			r.Delete("/", DeleteExpense)
		})
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}

func ExpenseContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expenseID := chi.URLParam(r, "id")
		id, _:=strconv.Atoi(expenseID)
		result := expenseCollection.Find()
		var expenses types.Expenses
		err := result.All(&expenses)
		if err!= nil{
			log.Println(err)
		}
		for _, expense := range expenses {
			if expense.Id == id {
				ctx := context.WithValue(r.Context(), "expense", expense )
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}



func CreateExpense(writer http.ResponseWriter, request *http.Request) {

	var req requests.CreateExpenseRequest

	err := render.Bind(request, &req)
	if err!= nil {
		log.Println(err)
		return
	}
	_ , err = expenseCollection.Insert(*req.Expense)
	if err != nil{
		log.Println(err)
		return
	}
	j, err := json.Marshal(req.Expense)
	if err != nil{
		render.Render(writer,request,ers.ErrRender(err))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, _ = fmt.Fprintf(writer, `{"success": true, "data": %v}`, string(j))

	// to check whether expense was added or not
	//result := expenseCollection.Find()
	//var print types.Expenses
	//err = result.All(&print)
	//if err!= nil{
	//	log.Println(err)
	//}
	//for _, p := range print {
	//	fmt.Println(p)
	//}
}
func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(types.Expense)
	err := render.Render(writer, request, response.Listexpense(&expense))
	if err != nil{
		render.Render(writer,request,ers.ErrRender(err))
		return
	}
}

func ListAllExpense(writer http.ResponseWriter, request *http.Request) {
	result := expenseCollection.Find()
	var print types.Expenses
	err := result.All(&print)
	if err!= nil{
		log.Println(err)
	}
	err = render.Render(writer, request, response.ExpensesResponse(&print))
	if err != nil{
		render.Render(writer,request,ers.ErrRender(err))
		return
	}
}

func UpdateExpense(writer http.ResponseWriter, request *http.Request) {

	expense := request.Context().Value("expense").(types.Expense)

	var req requests.UpdateExpenseRequest

	err := render.Bind(request, &req)
	if err != nil {
		log.Println(err)
		return
	}
	id := expense.Id
	fmt.Println(id)
	result := expenseCollection.Find("id", id)
	err = result.One(&expense)
	if err!= nil{
		log.Println(err)
	}
	fmt.Println(result)
	err = result.Update(*req.Expense)

	//err = expenseCollection.UpdateReturning(*req.Expense)
	if err!= nil{
		log.Println(err)
	}
	//expenses[expense.Id-1] = *req.Expense


	err = render.Render(writer, request, response.Listexpense(&expense))
	if err != nil{
		render.Render(writer,request,ers.ErrRender(err))
		return
	}
}

func DeleteExpense(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(types.Expense)
	expenses = append(expenses[:expense.Id-1], expenses[expense.Id:]...)
}