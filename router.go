package main

import (
	"net/http"
	"io/ioutil"
	// "reflect"
	"encoding/json"
	"goparse/Godeps/_workspace/src/github.com/gorilla/mux"
	"bytes"
  "log"
	"fmt"
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

	router := mux.NewRouter() //.StrictSlash(true)
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
    log.Println("========INCOMING REQUEST========")
    body, _ := ioutil.ReadAll(r.Body)
    // output/log original body
    log.Println(fmt.Sprintf("ORIGINAL BODY: %v", string(body)))
    // output/log original request
    log.Println(fmt.Sprintf("ORIGINAL REQUEST: %v", r))
    var params map[string]interface{}
    err := json.Unmarshal(body, &params)
    if err == nil {
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
    }
   	next.ServeHTTP(w, r)
  })

  // if (req.method === 'POST' && req.body._method) {
  //   req.originalMethod = req.method;
  //   req.method = req.body._method;
  //   delete req.body._method;
  // }
  // next();
}
