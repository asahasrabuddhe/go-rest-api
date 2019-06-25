package types

import (
	"time"
)

type Expense struct {
	Id          int       ` json:"id" bson:"id"`
	Description string    `json:"description" bson:"description"`
	Type        string    `json:"type" bson:"type"`
	Amount      float64   `json:"amount" bson:"amount"`
	CreatedOn   time.Time `json:"created_on" bson:"created_on"`
	UpdatedOn   time.Time `json:"updated_on" bson:"updated_on"`
}

type Expenses []Expense
