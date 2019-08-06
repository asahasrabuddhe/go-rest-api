package main

import (
	"bytes"
	"fmt"
	"go-rest-api/Interfaces"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

)
var D2 Interfaces.Databases


func TestGetall(t *testing.T) {
	fakeobj := &fake{D2}
	req, err := http.NewRequest("GET", "/expenses", nil)
	if err != nil {
		fmt.Println(err)
	}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeobj.GetAll)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":1,"description":"test5","type":" Cash","amount":4,"created_on":"0001-01-01T00:00:00Z","updated_on":"0001-01-01T00:00:00Z"},{"id":2,"description":"test5","type":" Cash","amount":4,"created_on":"0001-01-01T00:00:00Z","updated_on":"0001-01-01T00:00:00Z"}]`

	if strings.Compare(rec.Body.String(),expected)==0{
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}

}
func TestGetId(t *testing.T) {

	fakeobj := &fake{D2}

	req, err := http.NewRequest("GET", "/expenses/1", nil)


	if err != nil {
		t.Fatal(err)
	}
	//q := req.URL.Query()
	//q.Add("id", "1")
	//req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeobj.ArticleCtx)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":1,"description":"something","type":"cash","amount":123.456,"created_on":"0001-01-01T00:00:00Z","updated_on":"0001-01-01T00:00:00Z"}`
	if strings.Compare(rec.Body.String(),expected)==0{
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}
func TestCreate(t *testing.T) {
	fakeobj := &fake{D2}

	var jsonStr = []byte(`{"id": 3, "type": "cash", "amount": 123.456, "description": "created"}`)
	req, err := http.NewRequest("POST", "/expenses", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeobj.Create)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":3,"description":"created","type":"cash","amount":123.456,"created_on":"0001-01-01T00:00:00Z","updated_on":"0001-01-01T00:00:00Z"}`
	if strings.Compare(rec.Body.String(),expected)==0{
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}

}
func TestUpdate(t *testing.T) {
	fakeobj := &fake{D2}

	var jsonStr = []byte(`{"id": 2, "type": "cash", "amount": 123.456, "description": "updated"}`)

	req, err := http.NewRequest("PUT", "/expenses", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "2")
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeobj.Update)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":2,"description":"updated","type":"cash","amount":123.456,"created_on":"0001-01-01T00:00:00Z","updated_on":"0001-01-01T00:00:00Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
func TestDelete(t *testing.T) {
	var D2 Interfaces.Databases
	fakeobj := &fake{D2}

	req, err := http.NewRequest("DELETE", "/expenses/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	//q := req.URL.Query()
	//q.Add("id", "2")
	//req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fakeobj.Delete)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":2,"description":"updated","type":"cash","amount":123.456,"created_on":"0001-01-01T00:00:00Z","updated_on":"0001-01-01T00:00:00Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}