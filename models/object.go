package models

import (
	"errors"
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"goparse/connection"
	"log"
	"time"
	"math/rand"
	"fmt"
)

func ObjectCreate(object map[string]interface{}, className string) error {

	c, session := connection.GetCollection(className)
	defer session.Close()
	currentTime := time.Now()
	object["_id"] = RandomString(10)
	object["_created_at"] = currentTime
	object["_updated_at"] = currentTime
	fmt.Println(object)
	err := c.Insert(object)
	if err != nil {
		log.Fatal(err)
		return err

	}
	return err
}

func FindObjectById(objectId, className string) (map[string]interface{}, error) {

	c, session := connection.GetCollection(className)
	defer session.Close()
	object := map[string]interface{}{}

	err := c.Find(bson.M{"_id": objectId}).One(&object)
	if err != nil {
		return object, errors.New("object not found for get")
	}
	object["objectId"] = object["_id"]

	object["createdAt"] = object["_created_at"]
	object["updatedAt"] = object["_updated_at"]
	delete(object, "_id")
	delete(object, "_created_at")
	delete(object, "_updated_at")
	// _ = "breakpoint"
	return object, err
}

func ObjectUpdate(objectUpdates bson.M, objectId, className string) error {
	
	doc := bson.M{"$set": objectUpdates}
	c, session := connection.GetCollection(className)
	defer session.Close()
	query := bson.M{"_id": objectId}
	err := c.Update(query, doc)
	return err
	
}

func ObjectDestroy(objectId, className string) error {

	c, session := connection.GetCollection(className)
	defer session.Close()
	return c.Remove(bson.M{"_id": objectId})
}


// func SchemaQuery(query bson.M) (map[string]interface{}, error) {

// 	c, session := connection.GetCollection("_SCHEMA")
// 	defer session.Close()
// 	var result map[string]interface{}
// 	err := c.Find(query).One(&result)
// 	// _ = "breakpoint"
// 	if err != nil {
// 		// Return empty object and err if there is an error
// 		return result, err
// 	}

// 	return result, err
// }

// func SchemaIndex() ([]interface{}, error) {

// 	c, session := connection.GetCollection("_SCHEMA")
// 	defer session.Close()
// 	var results []interface{}
// 	err := c.Find(bson.M{}).All(&results)
// 	// _ = "breakpoint"
// 	if err != nil {
// 		// Return empty object and err if there is an error
// 		return results, err
// 	}

// 	return results, err
// }


// func (object Schema) Update(id string, doc bson.M) (Schema, error) {
// 	if bson.IsObjectIdHex(id) {
// 		c, session := connection.GetCollection("_SCHEMA")
// 		defer session.Close()
// 		query := bson.M{"_id": bson.ObjectIdHex(id)}
// 		doc["$set"].(bson.M)["updated_at"] = time.Now()
// 		err := c.Update(query, doc)
// 		// Upon successful update, we retrive the updated object
// 		// from db and return it. WARNING: this is an additional query
// 		if err == nil {
// 			return Schema{}.Find(id)
// 		} else {
// 			return Schema{}, err
// 		}
// 	}
// 	return Schema{}, errors.New("Invalid id")
// }

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
