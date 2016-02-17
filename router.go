package main

import (
	"net/http"
	// "io/ioutil"
	"goparse/Godeps/_workspace/src/github.com/gorilla/mux"
)

// TODOs: add middleware
// var allowMethodOverride = function(req, res, next) {
//   if (req.method === 'POST' && req.body._method) {
//     req.originalMethod = req.method;
//     req.method = req.body._method;
//     delete req.body._method;
//   }
//   next();
// };
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}

func allowMethodOverride(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // body, _ := ioutil.ReadAll(r.Body)

    // if body["_method"] != "/" {
    //   return
    // }
    // r.Method = "GET"
    r.Method = "GET"
   	next.ServeHTTP(w, r)
  })

  // if (req.method === 'POST' && req.body._method) {
  //   req.originalMethod = req.method;
  //   req.method = req.body._method;
  //   delete req.body._method;
  // }
  // next();
}
