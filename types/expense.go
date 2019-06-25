package types

import (
	"google.golang.org/genproto/googleapis/type/date"
)

type Expense struct {
	Id          int       `json:"id" db:"id,omitempty" bson:"id"`
	Description string    `json:"description" db:"description" bson:"description"`
	Type        string    `json:"type" db:"type" bson:"type"`
	Amount      float64   `json:"amount" db:"amount" bson:"amount"`
	CreatedOn   date.Date `json:"created_on" db:"created_on" bson:"created_on"`
	UpdatedOn   date.Date `json:"updated_on" db:"updated on" bson:"updated_on"`
}

type Expenses []Expense




