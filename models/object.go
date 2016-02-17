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
	// retrieve schema map first
	schema, err := SchemaQuery(bson.M{"_id": className})
	if err != nil {
		return object, err
	}	
	_ = parseObject(object, schema)
				
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


func QueryObject(query bson.M, className string) ([]map[string]interface{}, error) {

	c, session := connection.GetCollection(className)
	defer session.Close()
	var objects []map[string]interface{}

	err := c.Find(query).All(&objects)
	if err != nil {
		return objects, errors.New("object not found for get")
	}

	// retrieve schema map first
	schema, err := SchemaQuery(bson.M{"_id": className})
	if err != nil {
		return objects, err
	}

	if len(objects) > 0 {
		for _, object := range objects {
			_ = parseObject(object, schema)
				
		} 
	
	}

	return objects, err
}

// parseObject assumes a good map
func parseObject(object, schema map[string]interface{}) error {
	object["objectId"] = object["_id"]
	object["createdAt"] = object["_created_at"]
	object["updatedAt"] = object["_updated_at"]
	delete(object, "_id")
	delete(object, "_created_at")
	delete(object, "_updated_at")
	
	for key, value := range object {
		switch schema[key] {
		default:
			// do nothing
		case "date":
			if object[key] != nil {
				object[key] = map[string]interface{}{"__type": "Date", "iso": value}
			}
		case "geopoint":
			if object[key] != nil {
				object[key] = map[string]interface{}{"__type": "GeoPoint", "latitude": value.([]interface{})[0], "longitude": value.([]interface{})[1]}
			}
		}
	}

	return nil

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

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
