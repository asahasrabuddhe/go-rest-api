package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"strconv"

	//v1 "github.com/go-chi/chi/_examples/versions/presenter/v1"
	"github.com/go-chi/chi/middleware"
	//"github.com/pkg/errors"
	"log"

	//"google.golang.org/genproto/googleapis/type/date"
	"io/ioutil"
	//"log"
	"net/http"
	//"context"
)

type Expense struct {
	Id          float64   `json:"id"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	//CreatedOn   date.Date `json:"created_on" `
	//UpdatedOn   date.Date `json:"updated_on"`
}

type Expenses []Expense

var expenses Expenses

//var expense = []*Expense{
//	{Id: 1, Description: "First", Type:"shopping", Amount: 1500.00},
//	{Id: 2, Description: "Second", Type:"Car", Amount: 1500000.00},
//}

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
	})//https://github.com/asahasrabuddhe/go-rest-api.git

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

	//if val, ok := data["id"].(float64); ok {
	//	expense.Id = val
	//}

	if val, ok := data["description"].(string); ok {
		expense.Description = val
	}

	if val, ok := data["type"].(string); ok {
		expense.Type = val
	}

	if val, ok := data["amount"].(float64); ok {
		expense.Amount = val
		expense.Id = val+1
	}

	expenses = append(expenses, *expense)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, _ = fmt.Fprintln(writer, `{"success": true}`)
}

func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
	vars := chi.URLParam(request,"id")
	key, _ := strconv.Atoi(vars)

	for _, expense := range expenses{
		if expense.Id == float64(key){
			json.NewEncoder(writer).Encode(expense)

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

}

//func dbRemoveExpense(id int) (*Expense,error){
//		for i, a := range expense {
//			if a.Id == id {
//				expense = append((expense)[:i], (expense)[i+1:]...)
//				return a, nil
//			}
//		}
//		return nil, errors.New("No Expense like this.")
//	}
//
//}

func DeleteExpense(writer http.ResponseWriter, request *http.Request) {

	//var err error
	//
	//article := request.Context.Value("expenses").(*Expenses)
	//
	//article ,err = dbRemoveExpense(expenses.id)
	//if err != nil{
	//	render.Render(writer, request, ErrInvalidRequest(err))
	//	return
	//}
	//
	//render.Render(writer, request, v1.NewExpenseResponse(article))

	////parse the path parameters
	vars := chi.URLParam(request,"id")
	//extract the id need to delete
	id, _ := strconv.Atoi(vars)

	for index, expense := range expenses {
		if expense.Id == float64(id) {
			expenses = append(expenses[:index], expenses[index+1:]...)
		}
	}
}
