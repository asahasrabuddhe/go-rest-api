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


	expense.Id = len(expenses) + 1

	if val, ok := data["description"].(string); ok {
		expense.Description = val
	}

	if val, ok := data["type"].(string); ok {
		expense.Type = val
	}

	if val, ok := data["amount"].(float64); ok {
		expense.Amount = val
	}

	expenses = append(expenses, *expense)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, _ = fmt.Fprintln(writer, `{"success": true}`)
}

func ListOneExpense(writer http.ResponseWriter, request *http.Request) {

	tempId,_ := strconv.Atoi( chi.URLParam(request, "id") )

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	for _,expense := range(expenses){
		if expense.Id == tempId {
			json.NewEncoder(writer).Encode(expense)
			return
		}
	}

}

func ListAllExpense(writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_ = encoder.Encode(expenses)
}

func UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	tempId,_ := strconv.Atoi( chi.URLParam(request, "id"))

	for key,expense := range expenses {
		if expense.Id == tempId {
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

			expense.Id = expenses[key].Id

			if val, ok := data["description"].(string); ok {
				expense.Description = val
			}else{
				expense.Description = expenses[key].Description
			}

			if val, ok := data["type"].(string); ok {
				expense.Type = val
			}else{
				expense.Type = expenses[key].Type
			}

			if val, ok := data["amount"].(float64); ok {
				expense.Amount = val
			}else{
				expense.Amount = expenses[key].Amount
			}

			expenses = append(expenses[:key], expenses[key + 1 :]...)
			expenses = append(expenses, *expense)

		}
	}

}

func DeleteExpense(writer http.ResponseWriter, request *http.Request) {
	tempId,_ := strconv.Atoi( chi.URLParam(request, "id") )

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	for key,expense := range expenses {
		if expense.Id == tempId {
			expenses = append(expenses[:key], expenses[key + 1 :]... )
			return
		}
	}
}