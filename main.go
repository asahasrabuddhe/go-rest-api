package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go-rest-api/expenseDb"
	//"go-rest-api/expenseHandler"
	"go-rest-api/interfaces"
	"go-rest-api/types"
	"log"
	"net/http"
)

var expenses types.Expenses
var mh *expenseDb.MongoHandler

func main() {

	mongoDbConnection := "mongodb://localhost:27017"
	mh = expenseDb.NewHandler(mongoDbConnection)
	r := registerRoutes(mh)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func registerRoutes(db_i interfaces.Databases) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", db_i.Create)
		r.Get("/", db_i.GetAll)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(db_i.ExpenseCtx)
			r.Get("/", db_i.GetOne)
			r.Put("/", db_i.Update)
			r.Delete("/", db_i.Delete)
		})
	})
	return r
}
