package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func WithRecover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("Recovered from panic: %+v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		handler.ServeHTTP(w, req)
	})
}
