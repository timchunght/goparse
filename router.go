package main

import (
	"net/http"
	"io/ioutil"
	// "reflect"
	"encoding/json"
	"goparse/Godeps/_workspace/src/github.com/gorilla/mux"
	"bytes"
	// "fmt"
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
    body, _ := ioutil.ReadAll(r.Body)
    var params map[string]interface{}
    // fmt.Println(string(body))
		_ = json.Unmarshal(body, &params)
    // fmt.Println(reflect.TypeOf(body))
    if params["_method"] == "GET" {
    	r.Method = "GET"   
    }
    delete(params, "_method")
    delete(params, "_ApplicationId")
    delete(params, "_ClientVersion")
    delete(params, "_InstallationId")
    delete(params, "_JavaScriptKey")
    delete(params, "_SessionToken")
    delete(params, "_MasterKey")
    newBody, _ := json.Marshal(params)
    r.Body = ioutil.NopCloser(bytes.NewReader(newBody))
    
   	next.ServeHTTP(w, r)
  })

  // if (req.method === 'POST' && req.body._method) {
  //   req.originalMethod = req.method;
  //   req.method = req.body._method;
  //   delete req.body._method;
  // }
  // next();
}
