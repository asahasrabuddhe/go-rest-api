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
var expenses types.Expenses
var temp types.Expense
var req requests.CreateExpenseRequest
var err error
var connstr = "root:root@tcp(127.0.0.1:3306)/Expense?charset=utf8&parseTime=True"
var db1 Interfaces.Databases
func main() {

	db,err := gorm.Open("mysql", connstr)
	if err != nil {
		fmt.Println(err)
	}
	if (!db.HasTable(&types.Expense{})) {
		db.AutoMigrate(&types.Expense{})
	}
	set:=&Mysql{db}
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
		expenseID := chi.URLParam(r, "id")
		var temp types.Expense
		Db:= db.Db.Table("expenses").Where("id = ?", expenseID).Find(&temp)

		if Db.RowsAffected == 0{
			err=errors.New("ID not Found")
			render.Render(w, r, errrs.ErrRender(err))
			return
		} else{
			ctx := context.WithValue(r.Context(), "expense", Db)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
func (db *Mysql)Create(writer http.ResponseWriter, request *http.Request) {
	err = render.Bind(request, &req)
	temp:=*req.Expense
	temp.CreatedOn=time.Now()
	temp.UpdatedOn=time.Now()
	db.Db.Create(&temp)
	render.Render(writer, request, responses.List1expense(req.Expense))
}
func (db *Mysql)Update(writer http.ResponseWriter, request *http.Request) {
	db.Db = request.Context().Value("expense").(*gorm.DB)
	err:= render.Bind(request,&req)
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
	temp:=*req.Expense
	temp.UpdatedOn=time.Now()
	dB:= db.Db.Update(&temp)
			if(dB.RowsAffected == 0){
				err:=errors.New("Expense not found")
				render.Render(writer,request,errrs.ErrRender(err))
				return
			}else{
				_=render.Render(writer, request, responses.List1expense(&temp))
			}
}
func (db *Mysql) GetId(writer http.ResponseWriter, request *http.Request) {
	db.Db = request.Context().Value("expense").(*gorm.DB)
	dB:=db.Db.Find(&temp)
	if(dB.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}else{
		_=render.Render(writer, request, responses.List1expense(&temp))
	}
}
func (db *Mysql) GetAll(writer http.ResponseWriter, request *http.Request) {
	db.Db.Find(&expenses)
	err = render.Render(writer, request, responses.NewExpensesResponse(&expenses))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
}
func (db *Mysql)Delete(writer http.ResponseWriter, request *http.Request) {
	db.Db = request.Context().Value("expense").(*gorm.DB)
	dB:= db.Db.Delete(&temp)
	if(dB.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}else{
		_=render.Render(writer, request, responses.List1expense(&temp))
	}
}
