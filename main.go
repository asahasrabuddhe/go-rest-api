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
	"go-rest-api/errrs"
	"go-rest-api/requests"
	"go-rest-api/responses"
	"go-rest-api/types"
	"log"
	"net/http"
	"time"
)
var db *gorm.DB
var expenses types.Expenses
var temp types.Expense
var req requests.CreateExpenseRequest
var err error
func main() {
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Expense?charset=utf8&parseTime=True")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}
	if(!db.HasTable(&types.Expense{}) ) {
		db.AutoMigrate(&types.Expense{})
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", CreateExpense)
		r.Get("/", ListAllExpense)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(ArticleCtx)
			r.Get("/", ListOneExpense)
			r.Put("/", UpdateExpense)
			r.Delete("/", DeleteExpense)
		})
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expenseID := chi.URLParam(r, "id")

		db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Expense?charset=utf8&parseTime=True")
		if err != nil {
			fmt.Println(err)
		}else{
			fmt.Println("Connection established")
		}
		var temp types.Expense
		Db:= db.Table("expenses").Where("id = ?", expenseID).Find(&temp)

		if Db.RowsAffected == 0{
			err=errors.New("ID not Found")
			render.Render(w, r, errrs.ErrRender(err))
			return
		} else{
			ctx := context.WithValue(r.Context(), "expense", Db )
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
func CreateExpense(writer http.ResponseWriter, request *http.Request) {
	err = render.Bind(request, &req)
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Expense?charset=utf8&parseTime=True")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}
	temp:=*req.Expense
	temp.CreatedOn=time.Now()
	temp.UpdatedOn=time.Now()
	db.Create(&temp)
	render.Render(writer, request, responses.List1expense(req.Expense))
}
func UpdateExpense(writer http.ResponseWriter, request *http.Request) {
	db := request.Context().Value("expense").(*gorm.DB)
	err:= render.Bind(request,&req)
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
	temp:=*req.Expense
	temp.UpdatedOn=time.Now()
	Db:= db.Update(&temp)
			if(Db.RowsAffected == 0){
				err:=errors.New("Expense not found")
				render.Render(writer,request,errrs.ErrRender(err))
				return
			}else{
				_=render.Render(writer, request, responses.List1expense(&temp))
			}
}
func ListOneExpense(writer http.ResponseWriter, request *http.Request) {
	db := request.Context().Value("expense").(*gorm.DB)
	Db:= db.Find(&temp)
	if(Db.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}else{
		_=render.Render(writer, request, responses.List1expense(&temp))
	}
}
func ListAllExpense(writer http.ResponseWriter, request *http.Request) {
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Expense?charset=utf8&parseTime=True")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}
	db.Find(&expenses)
	err = render.Render(writer, request, responses.NewExpensesResponse(&expenses))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
}
func DeleteExpense(writer http.ResponseWriter, request *http.Request) {
	db := request.Context().Value("expense").(*gorm.DB)
	Db:= db.Delete(&temp)
	if(Db.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}else{
		_=render.Render(writer, request, responses.List1expense(&temp))
	}
}
