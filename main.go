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
<<<<<<< HEAD
	Id          float64   `json:"id"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
=======
	Id          float64 `json:"id"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
>>>>>>> f784dea5bce3611fd7081e9feee701e0abb382e5
	//CreatedOn   date.Date `json:"created_on" `
	//UpdatedOn   date.Date `json:"updated_on"`
}

type Expenses []Expense

var (
	expenses Expenses
	expense1 Expenses
)

//var expense = []*Expense{
//	{Id: 1, Description: "First", Type:"shopping", Amount: 1500.00},
//	{Id: 2, Description: "Second", Type:"Car", Amount: 1500000.00},
//}

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
<<<<<<< HEAD
	})//https://github.com/asahasrabuddhe/go-rest-api.git
=======
	}) //https://github.com/asahasrabuddhe/go-rest-api.git
>>>>>>> f784dea5bce3611fd7081e9feee701e0abb382e5

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
<<<<<<< HEAD
		expense.Id = val+1
=======
		expense.Id = val + 1
>>>>>>> f784dea5bce3611fd7081e9feee701e0abb382e5
	}

	expenses = append(expenses, *expense)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, _ = fmt.Fprintln(writer, `{"success": true}`)
}

func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
<<<<<<< HEAD
	vars := chi.URLParam(request,"id")
	key, _ := strconv.Atoi(vars)

	for _, expense := range expenses{
		if expense.Id == float64(key){
			json.NewEncoder(writer).Encode(expense)
			return//infinte loop
			//ctrl +I
=======
	vars := chi.URLParam(request, "id")
	key, _ := strconv.Atoi(vars)

	for _, expense := range expenses {
		if expense.Id == float64(key) {
			json.NewEncoder(writer).Encode(expense)
>>>>>>> f784dea5bce3611fd7081e9feee701e0abb382e5

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

	////parse the path parameters
	vars := chi.URLParam(request, "id")
	//extract the id need to delete
	id, _ := strconv.Atoi(vars)
	str, _ := ioutil.ReadAll(request.Body)
	var expense3 Expense
	var expense4 Expenses
	json.Unmarshal(str, &expense3)
	for index, exp := range expenses {
		if exp.Id == float64(id) {
			expense4 = append(expense1[:index], expense3)
			//json.NewEncoder(w).Encode(s4v)
			expense4 = append(expense4, expense1[index+1:]...)
		}
	}
	json.NewEncoder(writer).Encode(expense4)
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
<<<<<<< HEAD
	vars := chi.URLParam(request,"id")
=======
	vars := chi.URLParam(request, "id")
>>>>>>> f784dea5bce3611fd7081e9feee701e0abb382e5
	//extract the id need to delete
	id, _ := strconv.Atoi(vars)

	for index, expense := range expenses {
		if expense.Id == float64(id) {
			expenses = append(expenses[:index], expenses[index+1:]...)
		}
	}
}
