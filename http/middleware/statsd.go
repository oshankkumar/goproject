package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/alexcesaro/statsd.v2"
)

func StatsD(client *statsd.Client) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if client == nil {
				next.ServeHTTP(w, r)
				return
			}

			responseWriter, ok := w.(negroni.ResponseWriter)
			if !ok {
				responseWriter = negroni.NewResponseWriter(w)
			}

			name := strings.ReplaceAll(reqRouteName(r), " ", ".")

			client.Increment(fmt.Sprintf("http.req.%s", name))
			defer func(begin time.Time) {
				client.Timing(fmt.Sprintf("http.latency.%s.ms", name), int(time.Since(begin)/time.Millisecond))
				client.Increment(fmt.Sprintf("http.resp.%s.status.%d", name, responseWriter.Status()))
			}(time.Now())

			next.ServeHTTP(responseWriter, r)
		})
	}
}
