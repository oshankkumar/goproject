package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"

	"example.com/example/goproject/pkg/datastructure"
)

func reqRouteName(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if nil == route {
		return "NotFoundHandler"
	}
	if n := route.GetName(); n != "" {
		return n
	}
	if n, _ := route.GetPathTemplate(); n != "" {
		return r.Method + " " + n
	}
	n, _ := route.GetHostTemplate()
	return r.Method + " " + n
}

func NewRelic(app *newrelic.Application, ignorePathWithPrefix ...string) mux.MiddlewareFunc {
	s := make(datastructure.StringSet)
	s.Add(ignorePathWithPrefix...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s.IsPrefixOf(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			name := reqRouteName(r)
			txn := app.StartTransaction(name)
			defer txn.End()
			txn.SetWebRequestHTTP(r)
			w = txn.SetWebResponse(w)
			r = newrelic.RequestWithTransactionContext(r, txn)
			next.ServeHTTP(w, r)
		})
	}
}
