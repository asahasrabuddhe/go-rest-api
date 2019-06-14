package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/genproto/googleapis/type/date"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Expense struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	CreatedOn   date.Date `json:"created_on" `
	UpdatedOn   date.Date `json:"updated_on"`
}

type Expenses []Expense

var expenses Expenses

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
			r.Get("/", ListOneExpense)
			r.Put("/", UpdateExpense)
			r.Delete("/", DeleteExpense)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func CreateExpense(writer http.ResponseWriter, request *http.Request) {
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "unable to read request body", 500)
	}

	var data map[string]interface{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		http.Error(writer, "unable to parse json request body", 422)
	}

	expense := new(Expense)


	if val, ok := data["description"].(string); ok {
		expense.Description = val
	}

	if val, ok := data["type"].(string); ok {
		expense.Type = val
	}

	if val, ok := data["amount"].(float64); ok {
		expense.Amount = val
		expense.Id=int(val)
	}

	expenses = append(expenses, *expense)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, _ = fmt.Fprintln(writer, `{"success": true}`)
}

func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
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
			json.NewEncoder(writer).Encode(article)
		}
	}
	if flag==0{
		http.Error(writer,"expense with ID "+vars+" not found",500)
	}


}

func ListAllExpense(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_ = encoder.Encode(expenses)
}

func UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	vars := chi.URLParam(request, "id")
	a, err := strconv.Atoi(vars)
	if err != nil {
		http.Error(writer, "unable to read request body", 500)
	}
	str, _ := ioutil.ReadAll(request.Body)
	var temp Expense
	var temp1 Expenses
	json.Unmarshal(str, &temp)
	for index, articles := range expenses {
		if (articles.Id == a) {
			temp1 = append(expenses[:index], temp)
			//json.NewEncoder(w).Encode(s4v)
			temp1 = append(temp1, expenses[index+1:]...)
		}
	}

	encoder := json.NewEncoder(writer)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_ = encoder.Encode(temp1)
}
func DeleteExpense(writer http.ResponseWriter, request *http.Request) {

	vars:= chi.URLParam(request,"id")
	a,err:=strconv.Atoi(vars)
	flag:=0
	if err !=nil{
		http.Error(writer,"ID of the expense not parsed",500)
	}
	for index, article := range expenses {
		if (article.Id==a) {
			flag=1
			expenses =append(expenses[:index], expenses[index+1:]...)

			encoder := json.NewEncoder(writer)

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)

			_ = encoder.Encode(expenses)

		}
	}
	if flag==0{
		http.Error(writer,"expense with ID "+vars+" not found",500)
	}


}
