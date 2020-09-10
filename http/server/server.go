package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"example.com/example/goproject/http/httputil"
	"example.com/example/goproject/http/router"
	"example.com/example/goproject/pkg/i18n"
)

type Middleware func(handler http.Handler) http.Handler

type Server struct {
	routers     []router.Router
	middlewares []mux.MiddlewareFunc
	httpServer  *http.Server
	addr        string
	tr          *i18n.Translator
}

func New(addr string, tr *i18n.Translator) *Server {
	return &Server{addr: addr, tr: tr}
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{Addr: s.addr, Handler: s.createMux()}
	return s.httpServer.ListenAndServe()
}

func (s *Server) createMux() *mux.Router {
	m := mux.NewRouter()

	for _, apiRouter := range s.routers {
		for _, r := range apiRouter.Routes() {
			f := httputil.MakeHTTPHandler(r.Handler(), s.tr)
			logrus.Debugf("Registering %s, %s", r.Method(), r.Path())
			m.Path(r.Path()).Methods(r.Method()).Handler(f)
		}
	}

	m.Use(s.middlewares...)

	return m
}

func (s *Server) Use(m ...mux.MiddlewareFunc) {
	s.middlewares = append(s.middlewares, m...)
}

func (s *Server) InitRouter(routers ...router.Router) {
	s.routers = append(s.routers, routers...)
}

func (s *Server) Stop(shutDownTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
	err := s.httpServer.Shutdown(ctx)
	cancel()
	return err
}
