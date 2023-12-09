package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

type HttpError struct {
	StatusCode int    `json:"code"`
	Status     string `json:"status"`
	Error      string `json:"error"`
}

func (h HttpError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, h.StatusCode)
	return nil
}

func RenderError(err error, statusCode int) render.Renderer {
	return &HttpError{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Error:      err.Error(),
	}
}
