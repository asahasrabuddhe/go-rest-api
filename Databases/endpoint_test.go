package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetall(t *testing.T) {
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/Expense1?charset=utf8&parseTime=True")
	set := &Mysql{db}
	req, err := http.NewRequest("GET", "/expenses", nil)
	if err != nil {
		fmt.Println(err)
	}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(set.GetAll)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":5,"description":"test3","type":" Cash","amount":3,"created_on":"2019-07-24T06:29:37Z","updated_on":"2019-07-24T06:29:37Z"},{"id":6,"description":"test4","type":" Cash","amount":4,"created_on":"2019-07-24T06:30:14Z","updated_on":"2019-07-24T06:30:14Z"}]`
	if strings.Compare(rec.Body.String(),expected)==0{
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}

}