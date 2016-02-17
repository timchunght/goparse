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
		"ObjectUpdate",
		[]string{"PUT"},
		"/classes/{className}/{objectId}",
		controllers.ObjectUpdate,
	},
	Route{
		"ObjectUpdate",
		[]string{"PUT"},
		"/classes/{className}",
		controllers.ObjectUpdate,
	},
	Route{
		"ObjectUpdate",
		[]string{"PUT"},
		"/classes/",
		controllers.ObjectUpdate,
	},
	Route{
		"ObjectQuery",
		[]string{"GET"},
		"/classes/{className}",
		controllers.ObjectQuery,
	},
	Route{
		"ObjectShow",
		[]string{"GET"},
		"/classes/{className}/{objectId}",
		controllers.ObjectShow,
	},
	Route{
		"ObjectDestroy",
		[]string{"DELETE"},
		"/classes/{className}/{objectId}",
		controllers.ObjectDestroy,
	},
	Route{
		"ObjectDestroy",
		[]string{"DELETE"},
		"/classes/{className}",
		controllers.ObjectDestroy,
	},
	Route{
		"ObjectDestroy",
		[]string{"DELETE"},
		"/classes/",
		controllers.ObjectDestroy,
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
