package middleware

import (
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"

	"example.com/example/goproject/pkg/context"
	"example.com/example/goproject/pkg/datastructure"
)

func WithLogger(ignorePaths ...string) mux.MiddlewareFunc {
	s := make(datastructure.StringSet)
	s.Add(ignorePaths...)

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s.IsPrefixOf(r.URL.Path) {
				handler.ServeHTTP(w, r)
				return
			}

			reqID := uuid.NewRandom()

			entry := logrus.WithFields(logrus.Fields{
				"reqID": reqID.String(),
			})

			ctx := context.WithReqID(r.Context(), reqID)
			ctx = context.WithLogger(ctx, entry)
			r = r.WithContext(ctx)

			w.Header().Set("X-Request-Id", reqID.String())

			res, ok := w.(negroni.ResponseWriter)
			if !ok {
				res = negroni.NewResponseWriter(w)
			}

			begin := time.Now()

			handler.ServeHTTP(res, r)

			entry.WithFields(logrus.Fields{
				"StartTime":  begin.Format(time.RFC3339),
				"Status":     res.Status(),
				"DurationMS": time.Since(begin).Milliseconds(),
				"Hostname":   r.Host,
				"Method":     r.Method,
				"Path":       r.URL.Path,
			}).Info("inbound http request")
		})
	}
}
