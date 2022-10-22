package middleware

import "net/http"

var middlewares []func(http.HandlerFunc) http.HandlerFunc

func Use(mws ...func(http.HandlerFunc) http.HandlerFunc) {
	for _, mw := range mws {
		middlewares = append(middlewares, mw)
	}
}

func Handler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
