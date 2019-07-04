package books

import (
	"errors"
	"github.com/asahasrabuddhe/rest-api/logger"
	"github.com/asahasrabuddhe/rest-api/renderers"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type Resource struct {}

func (res Resource) Routes() chi.Router {
    r := chi.NewRouter()

	r.Post("/", res.Create)
	r.Get("/", res.GetAll)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(res.Context)

		r.Get("/", res.GetOne)
		r.Put("/", res.Update)
		r.Delete("/", res.Delete)
	})

	return r
}

func (res Resource) Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := chi.URLParam(r, "id"); id != "" {
			if idInt, err := strconv.Atoi(id); err != nil {
				_ = render.Render(w, r, &renderers.ErrorResponse{
					Err:            err,
					HTTPStatusCode: 422,
					StatusText:     "unable to parse id",
				})
			} else {
				logger.LogEntrySetField(r, "book_id", idInt)

				// TODO - Add logic to retrieve resource for given ID, create new context with resource value and return

				_ = render.Render(w, r, &renderers.ErrorResponse{
					Err:            errors.New("resource not found"),
					HTTPStatusCode: 404,
					StatusText:     "resource not found",
				})
			}
		}
	})
}

func (res Resource) Create(writer http.ResponseWriter, request *http.Request) {
	var req Create

	err := render.Bind(request, &req)
	if err != nil {
		logger.LogEntrySetField(request, "error", err)
		return
	}

	// TODO - Add logic to persist resource

	_ = render.Render(writer, request, NewBookResponse(req.Book))
}

func (res Resource) GetOne(writer http.ResponseWriter, request *http.Request) {
	book := request.Context().Value("book").(Book)
	_ = render.Render(writer, request, NewBookResponse(&book))
}

func (res Resource) GetAll(writer http.ResponseWriter, request *http.Request) {
    var books Books

    // TODO - Add logic to populate above slice with data and return

	_ = render.Render(writer, request, NewBooksResponse(&books))
}

func (res Resource) Update(writer http.ResponseWriter, request *http.Request) {
	book := request.Context().Value("book").(Book)

	var req Update

	err := render.Bind(request, &req)
	if err != nil {
		logger.LogEntrySetField(request, "error", err)
		return
	}

	// TODO - Add logic to persist updated record

	_ = render.Render(writer, request, NewBookResponse(&book))
}

func (res Resource) Delete(writer http.ResponseWriter, request *http.Request) {
	 _ = request.Context().Value("book").(Book)

	// TODO - Add logic to delete record
}