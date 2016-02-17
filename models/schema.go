package models

import (
	// "errors"
	"encoding/json"
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"goparse/connection"
	"log"
	// "fmt"
	// "time"
)

func SchemaCreate(schema map[string]interface{}, className string) error {

	c, session := connection.GetCollection("_SCHEMA")
	defer session.Close()
	schema["_id"] = className
	schema["_metadata"] = metadata()
	err := c.Insert(schema)
	if err != nil {
		log.Fatal(err)
		return err

	}
	return err

}

// func (object Schema) Find(id string) (Schema, error) {

// 	c, session := connection.GetCollection("_SCHEMA")
// 	defer session.Close()
// 	result := Schema{}
// 	if bson.IsObjectIdHex(id) {
// 		err := c.FindId(bson.ObjectIdHex(id)).One(&result)
// 		// _ = "breakpoint"
// 		if err != nil {
// 			// Return empty Schema object and err if there is an error
// 			return Schema{}, err
// 		}
// 	} else {
// 		err := errors.New("Invalid id")

// 		return Schema{}, err
// 	}
// 	return result, nil
// }

func SchemaQuery(query bson.M) (map[string]interface{}, error) {

	c, session := connection.GetCollection("_SCHEMA")
	defer session.Close()
	var result map[string]interface{}
	err := c.Find(query).One(&result)
	// _ = "breakpoint"
	if err != nil {
		// Return empty object and err if there is an error
		return result, err
	}

	return result, err
}

func SchemaIndex() ([]interface{}, error) {

	c, session := connection.GetCollection("_SCHEMA")
	defer session.Close()
	var results []interface{}
	err := c.Find(bson.M{}).All(&results)
	// _ = "breakpoint"
	if err != nil {
		// Return empty object and err if there is an error
		return results, err
	}

	return results, err
}

func metadata() map[string]interface{} {
	metadata := map[string]interface{}{}
	b := []byte(`
	{
		"class_permissions": {
    "addField": {
      "*": true
    },
    "create": {
      "*": true
    },
    "delete": {
      "*": true
    },
    "find": {
      "*": true
    },
    "get": {
      "*": true
    },
    "readUserFields": [],
    "update": {
      "*": true
    },
    "writeUserFields": []
    }
  }`)
	_ = json.Unmarshal(b, &metadata)
	// fmt.Println(err)
	return metadata
}

// func (object Schema) Destroy(id string) error {

// 	c, session := connection.GetCollection("_SCHEMA")
// 	defer session.Close()
// 	return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
// }

func SchemaUpdate(schemaUpdates bson.M, className string) error {

	doc := bson.M{"$set": schemaUpdates}
	c, session := connection.GetCollection("_SCHEMA")
	defer session.Close()
	query := bson.M{"_id": className}
	err := c.Update(query, doc)
	return err
	// Upon successful update, we retrive the updated object
	// from db and return it. WARNING: this is an additional query

}
