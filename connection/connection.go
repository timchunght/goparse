package connection

import (
	"log"
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"os"
	// "errors"
	"fmt"
)

var (
	Session *mgo.Session
	// Database *mgo.Database
	err error
)

func Connect() {

	url := os.Getenv("MONGO_URL")
	fmt.Println(url)
	if url == "" {
		log.Printf("Cannot Find MONGO_URL, setting it to mongodb://localhost:27017")
		url = "mongodb://localhost:27017/modernplanit_parse_production"
	}

	Session, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	Session.SetMode(mgo.Monotonic, true)
	log.Printf("%s\t%s", "connected to", url)
}

func GetCollection(collectionName string) (*mgo.Collection, *mgo.Session) {
	session := Session.Copy()
	// _ = "breakpoint"
	return session.DB("").C(collectionName), session
}
