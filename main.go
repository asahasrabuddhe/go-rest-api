package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go-rest-api/expenseDB"
	"go-rest-api/interfaces"
	"go-rest-api/types"
	"log"
	"net/http"
)



var expenses types.Expenses

var mh *expenseDB.MongoHandler

func RouteHandler(db_i interfaces.Databases){
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", db_i.Create)
		r.Get("/", db_i.GetAll)
		//r.Put("/{id}",UpdateExpense)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(db_i.ArticleCtx)
			r.Get("/", db_i.GetId)
			r.Put("/", db_i.Update)
			r.Delete("/", db_i.Delete)
		})
	})
	log.Fatal(http.ListenAndServe(":8080", r))

}



func main() {

   // var db_i interfaces.Databases
	mongoDbConnection := "mongodb://localhost:27017"
	mh = expenseDB.NewHandler(mongoDbConnection)
    RouteHandler(mh)

}