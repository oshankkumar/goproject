package router

import (
	"net/http"

	"example.com/example/goproject/http/httputil"
)

type Router interface {
	Routes() []Route
}

type Route interface {
	Path() string
	Method() string
	Handler() httputil.Handler
}

type httpRouter struct {
	method  string
	path    string
	handler httputil.Handler
}

func (h httpRouter) Path() string {
	return h.path
}

func (h httpRouter) Method() string {
	return h.method
}

func (h httpRouter) Handler() httputil.Handler {
	return h.handler
}

func NewRoute(method, path string, handler httputil.Handler, middlewares ...httputil.Middleware) Route {
	return httpRouter{method: method, path: path, handler: httputil.NewHandlerChain(middlewares...).Then(handler)}
}

// NewGetRoute initializes a new route with the http method GET.
func NewGetRoute(path string, handler httputil.Handler, middlewares ...httputil.Middleware) Route {
	return NewRoute(http.MethodGet, path, handler, middlewares...)
}

// NewPostRoute initializes a new route with the http method POST.
func NewPostRoute(path string, handler httputil.Handler, middlewares ...httputil.Middleware) Route {
	return NewRoute(http.MethodPost, path, handler, middlewares...)
}

// NewPutRoute initializes a new route with the http method PUT.
func NewPutRoute(path string, handler httputil.Handler, middlewares ...httputil.Middleware) Route {
	return NewRoute(http.MethodPut, path, handler, middlewares...)
}
