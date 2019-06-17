package errrs

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrorResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(er error) render.Renderer {
	return &ErrorResponse{
		Err:            er,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      er.Error(),
	}
}