package main

import (
	expenseResource "github.com/asahasrabuddhe/rest-api/expenses/resource"
	"github.com/asahasrabuddhe/rest-api/server"
)

func main() {
	server.Initialize()

	server.Mount("/expenses", expenseResource.ExpenseResource{}.Routes())

	server.Serve()
}
