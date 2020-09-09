package httputil

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	return err
}
