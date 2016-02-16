package main

import (
	"net/http"

	"goparse/Godeps/_workspace/src/github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		// methods := make([]string, 10)
		// methods = append(methods, route.Method)
		// methods = append(methods, "GET")
		router.
			Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
