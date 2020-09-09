package handler

import (
	"net/http"

	"example.com/example/goproject/api/types"
	"example.com/example/goproject/http/httputil"
	"example.com/example/goproject/http/router"
)

func Ping() router.Router {
	return ping{}
}

type ping struct{}

func (p ping) Routes() []router.Route {
	return []router.Route{router.NewGetRoute("/ping", p)}
}

func (p ping) ServeHTTP(w http.ResponseWriter, _ *http.Request) error {
	return httputil.WriteJSON(w, http.StatusOK, &types.PingResponse{Status: "OK"})
}
