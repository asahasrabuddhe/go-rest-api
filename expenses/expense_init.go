package expenses

import (
	"github.com/asahasrabuddhe/rest-api/server"
)

func init() {
	Exp = []Expense{
		{Id: 1, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 2, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 3, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 4, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 5, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 6, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 7, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 8, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 9, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 10, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 11, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 12, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 13, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 14, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 15, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 16, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 17, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 18, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 19, Type: "cash", Amount: 123.456, Description: "something"},
		{Id: 20, Type: "cash", Amount: 123.456, Description: "something"},
	}

	server.Mount("/expenses", Resource{}.Routes())
}
