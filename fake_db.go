package Testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go-rest-api/Interfaces"
	"go-rest-api/errrs"
	"go-rest-api/requests"
	"go-rest-api/responses"
	"go-rest-api/types"
	"log"
	"net/http"
	"strconv"
)
var expenses types.Expenses
var expense types.Expense

type fake struct {
	Db Interfaces.Databases
}
func main() {
	expenses=types.Expenses{
		{Id: 1, Type: "test1", Amount: 123.456, Description: "something"},
		{Id: 2, Type: "test2", Amount: 123.456, Description: "something"},
	}
	var D1 Interfaces.Databases
	fakeobj := &fake{D1}
	handlereq(fakeobj)
}
func handlereq(db1 *fake) {
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
	log.Fatal(http.ListenAndServe(":8081", r))
}
func (db *fake)ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		flag:=1
		expenseID := chi.URLParam(r, "id")
		b,_:=strconv.Atoi(expenseID)
		fmt.Println(b)
		for _, exp := range expenses {

			if exp.Id == b {
				flag=0
				ctx := context.WithValue(r.Context(), "expense", exp )
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}

		if flag ==1{
			err=errors.New("ID not Found")
			render.Render(w, r,  errrs.ErrRender(err))
		}
	})
}
func (db *fake)Create(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest
	var err error
	err = render.Bind(request, &req)
	if err != nil {
		log.Println(err)
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
	expenses = append(expenses, *req.Expense)
	render.Render(writer, request, responses.List1expense(req.Expense))

}
func (db *fake)Update(writer http.ResponseWriter, request *http.Request) {
	var req requests.UpdateExpenseRequest
	err:= render.Bind(request,&req)
	if err != nil {
		log.Println(err)
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
	// for loop -> get index of expense
	for index, e := range expenses {
		if e.Id == expense.Id {
			expenses[index] = *req.Expense
		}
	}
	render.Render(writer, request, responses.List1expense(&expense))
}
func (db *fake) GetId(writer http.ResponseWriter, request *http.Request) {

	exp := request.Context().Value("expense").(types.Expense)
	_ = render.Render(writer, request, responses.List1expense(&exp))
}

func (db *fake) GetAll(writer http.ResponseWriter, request *http.Request) {
	err := render.Render(writer, request, responses.NewExpensesResponse(&expenses))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
}
func (db *fake)Delete(writer http.ResponseWriter, request *http.Request) {
	exp := request.Context().Value("expense").(types.Expense)
	expenses =append(expenses[:exp.Id], expenses[exp.Id+1:]...)

	err:=render.Render(writer,request,responses.NewExpensesResponse(&expenses))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
}
