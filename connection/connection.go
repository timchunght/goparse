package connection

import (
	"log"
	// "modernplanit/Godeps/_workspace/src/github.com/joho/godotenv"
	"modernplanit/Godeps/_workspace/src/gopkg.in/mgo.v2"
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
		log.Printf("Cannot Find MONGO_URL, setting it to mongodb://localhost:27017/modernplanit_development")
		url = "mongodb://localhost:27017/modernplanit_development"
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

// func (dbConf *DbConfig) parseConfig() {

// 	if os.Getenv("GO_ENV") == "test" {
// 		dbConf.Database = "modernplanit_test"
// 	} else if os.Getenv("GO_ENV") != "production" {
// 		err := godotenv.Load()
// 		if err != nil {
// 			log.Fatal("Error loading .env file")
// 		}
// 	}
// 	// Expected Url to be in the following format
// 	// "mongodb://localhost:27017"
// 	// mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb

// 	// To Be Implemented
// 	// dbConf := DbConfig{Host: os.Getenv("MONGO_HOST"),
// 	// 									 Database: os.Getenv("MONGO_DB"),
// 	// 									 Username: os.Getenv("MONGO_USERNAME"),
// 	// 									 Password: os.Getenv("MONGO_PASSWORD"),
// 	// 									 TestDatabase: os.Getenv("MONGO_TEST_DB")}

// 	dbConf.Url = os.Getenv("MONGO_URL")
// 	if dbConf.Url == "" {
// 		dbConf.Url = "mongodb://localhost:27017/sample"
// 	}

// }
