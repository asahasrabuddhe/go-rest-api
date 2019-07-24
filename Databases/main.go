package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-rest-api/Interfaces"
	"go-rest-api/errrs"
	"go-rest-api/requests"
	"go-rest-api/responses"
	"go-rest-api/types"
	"log"
	"net/http"
	"time"
)

type Mysql struct {
	Db *gorm.DB
}
var err error
var db1 Interfaces.Databases
func main() {

	dba, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	dba.Exec("CREATE DATABASE IF NOT EXISTS"+" Expense1")
	dba.Close()

	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Expense1?charset=utf8&parseTime=True")

	if err != nil {
		fmt.Println(err)
	}
		if (!db.HasTable(&types.Expense{})) {
			db.AutoMigrate(&types.Expense{})
		}
		set := &Mysql{db}
		handlerequest(set)
}
func handlerequest(db1 Interfaces.Databases){
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", db1.Create)
		r.Get("/", db1.GetAll)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(db1.ArticleCtx)
			r.Get("/", db1.GetId)
			r.Put("/", db1.Update)
			r.Delete("/", db1.Delete)
		})
	})
	log.Fatal(http.ListenAndServe(":8080", r))

}
func (db *Mysql)ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var temp types.Expense
		expenseID := chi.URLParam(r, "id")
		DB:= db.Db.Table("expenses").Where("id = ?", expenseID).Find(&temp)
		fmt.Println(temp)
		if DB.RowsAffected == 0{
			err=errors.New("ID not Found")
			render.Render(w, r, errrs.ErrRender(err))
			return
		} else{
			ctx := context.WithValue(r.Context(), "expense", temp)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
func (db *Mysql)Create(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest
	err = render.Bind(request, &req)
	temp:=*req.Expense
	temp.CreatedOn=time.Now()
	temp.UpdatedOn=time.Now()
	db.Db.Create(&temp)
	render.Render(writer, request, responses.List1expense(&temp))
}
func (db *Mysql)Update(writer http.ResponseWriter, request *http.Request) {
	var req requests.UpdateExpenseRequest
	expe:= request.Context().Value("expense").(types.Expense)
	err:= render.Bind(request,&req)
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
	expe=*req.Expense
	expe.UpdatedOn=time.Now()
	dB:= db.Db.Update(&expe)
			if(dB.RowsAffected == 0){
				err:=errors.New("unable to update")
				render.Render(writer,request,errrs.ErrRender(err))
				return
			}else{
				_=render.Render(writer, request, responses.List1expense(&expe))
			}
}
func (db *Mysql) GetId(writer http.ResponseWriter, request *http.Request) {

	expe:= request.Context().Value("expense").(types.Expense)
		_=render.Render(writer, request, responses.List1expense(&expe))
}

func (db *Mysql) GetAll(writer http.ResponseWriter, request *http.Request) {
	var exp types.Expenses
	db.Db.Find(&exp)
	err = render.Render(writer, request, responses.NewExpensesResponse(&exp))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
}
func (db *Mysql)Delete(writer http.ResponseWriter, request *http.Request) {
	expe:= request.Context().Value("expense").(types.Expense)
	dB:=db.Db.Delete(&expe)
	if(dB.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}else{
		_=render.Render(writer, request, responses.List1expense(&expe))
	}
}
