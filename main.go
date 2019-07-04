package main

import (
	_ "github.com/asahasrabuddhe/rest-api/books"
	_ "github.com/asahasrabuddhe/rest-api/expenses"
	"github.com/asahasrabuddhe/rest-api/server"
)

func main() {
	server.Serve()
}
