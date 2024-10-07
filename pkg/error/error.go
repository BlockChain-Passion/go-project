package errors

import (
	"net/http"

	"github.com/go-chi/render"
)

type IErrors interface {
	Render(rw http.ResponseWriter, r *http.Request) error
	NoError(tid string) render.Renderer
	GenericError(tId string, code int, msg string) render.Renderer
	ErrorConflict(tId string) render.Renderer
	ErrorForbiddenRequest(tid string) render.Renderer
	ErrorInvalidRequest(tid string) render.Renderer
	ErrorInternal(tId string) render.Renderer
	ErrorRenderRequest(tid string) render.Renderer
}

type ServiceErrors struct {
	RequestID      string `json:"requestId"`
	HttpStatusCode int    `json:"code"`
	StatusText     string `json:"status"`
}

func (s *ServiceErrors) Render(rw http.ResponseWriter, r *http.Request) error {
	render.Status(r, s.HttpStatusCode)
	return nil
}

func (s *ServiceErrors) NoError(tId string) render.Renderer {
	s.RequestID = tId
	s.HttpStatusCode = 200
	s.StatusText = "complete"
	return s
}

func (s *ServiceErrors) GenericError(tId string, code int, msg string) render.Renderer {
	s.RequestID = tId
	s.HttpStatusCode = code
	s.StatusText = msg
	return s
}

func (s *ServiceErrors) ErrorConflict(tId string) render.Renderer {
	s.RequestID = tId
	s.HttpStatusCode = 409
	s.StatusText = "conflict"
	return s
}

func (s *ServiceErrors) ErrorForbiddenRequest(tId string) render.Renderer {
	s.RequestID = tId
	s.HttpStatusCode = 401
	s.StatusText = "forbidden"
	return s
}

func (s *ServiceErrors) ErrorInvalidRequest(tId string) render.Renderer {
	s.RequestID = tId
	s.HttpStatusCode = 400
	s.StatusText = "Invalid Request"
	return s
}

func (s *ServiceErrors) ErrorInternal(tId string) render.Renderer {
	s.RequestID = tId
	s.HttpStatusCode = 500
	s.StatusText = "internal error"
	return s
}

func (s *ServiceErrors) ErrorRenderRequest(tId string) render.Renderer {
	s.RequestID = tId
	s.HttpStatusCode = 422
	s.StatusText = "error rendering response"
	return s
}
