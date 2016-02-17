package main

import (
	"goparse/controllers"
	"net/http"
)

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	// Route{
	// 	"TriviaCreate",
	// 	[]string{"POST"},
	// 	"/trivias",
	// 	controllers.TriviaCreate,
	// },
	// Route{
	// 	"TriviaQuery",
	// 	[]string{"GET"},
	// 	"/trivias",
	// 	controllers.TriviaQuery,
	// },
	// Route{
	// 	"TriviaShow",
	// 	[]string{"GET"},
	// 	"/trivias/{id}",
	// 	controllers.TriviaShow,
	// },
	// Route{
	// 	"TriviaDestroy",
	// 	[]string{"DELETE"},
	// 	"/trivias/{id}",
	// 	controllers.TriviaDestroy,
	// },
	// Route{
	// 	"TriviaUpdate",
	// 	[]string{"PUT"},
	// 	"/trivias/{id}",
	// 	controllers.TriviaUpdate,
	// },
	Route{
		"ObjectCreate",
		[]string{"POST"},
		"/classes/{className}",
		controllers.ObjectCreate,
	},
	Route{
		"ObjectShow",
		[]string{"GET"},
		"/classes/{className}/{objectId}",
		controllers.ObjectShow,
	},
	Route{
		"SchemaIndex",
		[]string{"GET"},
		"/schemas",
		controllers.SchemaIndex,
	},
	Route{
		"SchemaQuery",
		[]string{"GET"},
		"/schemas/{className}",
		controllers.SchemaQuery,
	},
}
