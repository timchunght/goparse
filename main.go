package main

import (
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"goparse/connection"
	"log"
	"net/http"
	"os"
)

var (
	session *mgo.Session
	err     error
	// db *mgo.Database
	// col *mgo.Collection
)

func main() {
	StartServer()

}

func StartServer() {
	connection.Connect()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), allowMethodOverride(router)))

}
