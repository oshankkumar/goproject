package httputil

import (
	"net/http"

	"github.com/sudo-suhas/xgo/errors"

	"example.com/example/goproject/api/types"
	"example.com/example/goproject/pkg/i18n"
)

const (
	HeaderAcceptLanguage = "Accept-Language"
	HeaderUserLocale     = "X-User-Locale"
)

type Middleware func(handler Handler) Handler

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

func MakeHTTPHandler(handler Handler, tr *i18n.Translator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler.ServeHTTP(w, r)
		if err == nil {
			return
		}

		i18nKey := errors.UserMsg(err)
		lang := getLanguange(r)

		WriteJSON(w, errors.StatusCode(err), &types.ErrorResponse{
			Code:     errors.WhatKind(err).Code,
			Message:  tr.Message(lang, i18n.TranslationConfig{Key: i18nKey}),
			Title:    tr.Title(lang, i18n.TranslationConfig{Key: i18nKey}),
			Severity: "error",
		})
	}
}

func getLanguange(r *http.Request) string {
	if lang := r.Header.Get(HeaderUserLocale); lang != "" {
		return lang
	}

	if lang := r.Header.Get(HeaderAcceptLanguage); lang != "" {
		return lang
	}

	return "en-ID"
}

type HandlerChain struct {
	middlewares []Middleware
}

func NewHandlerChain(m ...Middleware) *HandlerChain {
	return &HandlerChain{middlewares: append([]Middleware(nil), m...)}
}

func (c *HandlerChain) Then(h Handler) Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}
	return h
}
