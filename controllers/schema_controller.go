package controllers

import (
	// "encoding/json"
	"goparse/Godeps/_workspace/src/github.com/gorilla/mux"
	"net/http"
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"goparse/helpers"
	"goparse/models"
	// "time"
	// "io/ioutil"
	// "fmt"
	// "errors"
)

// func SchemaCreate(w http.ResponseWriter, r *http.Request) {
// 	var trivia models.Schema
// 	body, _ := ioutil.ReadAll(r.Body)
// 	_, paramsPresent := requiredBodyParamsCheck(body, []string{"event_id", "name", "description"})
// 	if paramsPresent == true {

// 		err := json.Unmarshal(body, &trivia)
// 		if err != nil {
// 			err := helpers.RenderJsonErr(w, http.StatusBadRequest, "Error parsing params")
// 			if err != nil {
// 				panic(err)
// 			}
// 			return
// 		}

// 		// make sure the create method does not have error
// 		err = (&trivia).Create()
// 		if err != nil {
// 			err := helpers.RenderJsonErr(w, http.StatusInternalServerError, "Error creating object")
// 			if err != nil {
// 				panic(err)
// 			}
// 			return
// 		}

// 		if err := helpers.RenderJson(w, http.StatusCreated, trivia) err != nil {
// 			panic(err)
// 		}

// 	} else {
// 		err := helpers.RenderJsonErr(w, http.StatusBadRequest, "Required params not found")
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// func SchemaShow(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := string(vars["id"])
// 	trivia, err := models.Schema{}.Find(id)
// 	if err != nil {
// 		// If we didn"t find it, 404
// 		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())
// 		if err != nil {
// 			panic(err)
// 		}

// 		return
// 	}

// 	err = helpers.RenderJson(w, http.StatusOK, trivia)
// 	if err != nil {
// 		panic(err)
// 	}

// }

func SchemaQuery(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	className := string(vars["className"])

	result, err := models.SchemaQuery(bson.M{"_id": className})
	if err != nil {
		// If we encounter error parsing query params
		err = helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())
		if err != nil {
			panic(err)
		}

		return
	}

	err = helpers.RenderJson(w, http.StatusOK, result)
	if err != nil {
		panic(err)
	}

}

func SchemaIndex(w http.ResponseWriter, r *http.Request) {

	results, err := models.SchemaIndex()
	if err != nil {
		// If we encounter error parsing query params
		err = helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())
		if err != nil {
			panic(err)
		}

		return
	}

	err = helpers.RenderJson(w, http.StatusOK, results)
	if err != nil {
		panic(err)
	}

}

// func SchemaDestroy(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := string(vars["id"])
// 	err := models.Schema{}.Destroy(id)

// 	if err != nil {
// 		// If we didn"t find it, 404
// 		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())
// 		if err != nil {
// 			panic(err)
// 		}

// 		return
// 	}

// 	// If the destroy action has no error, return the success message
// 	err = helpers.RenderJson(w, http.StatusOK, map[string]string{"message": "Successfully deleted"})
// 	if err != nil {
// 		panic(err)
// 	}

// }

// func SchemaUpdate(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := string(vars["id"])
// 	body, _ := ioutil.ReadAll(r.Body)

// 	paramKeys := []string{"event_id", "name", "description"}
// 	requiredParamKeys := paramKeys
// 	doc, err := getUpdateDocFromBody(body, requiredParamKeys, paramKeys)
// 	if err != nil {
// 		// If we didn"t find it, 404
// 		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())

// 		if err != nil {
// 			panic(err)
// 		}

// 		return
// 	}

// 	trivia, err := models.Schema{}.Update(id, doc)

// 	if err != nil {

// 		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())

// 		if err != nil {
// 			panic(err)
// 		}
// 		return
// 	}

// 	helpers.RenderJson(w, http.StatusOK, trivia)

// }

func MongoFieldTypeToSchemaAPIType(typeStr string) map[string]string {
  if (string(typeStr[0]) == "*") {
    return map[string]string{ "type": "Pointer", "targetClass": string(typeStr[1:len(typeStr)])}
  }

  if (string(typeStr[0:len("relation<")]) == "relation<") {
  	return map[string]string{"type": "Relation", "targetClass": string(typeStr[len("relation<"):len(typeStr)])}
  
  }
  switch typeStr {
  case "number":   
  	return map[string]string{"type": "Number"}
  case "string":   
  	return map[string]string{"type": "String"}
  case "boolean":  
  	return map[string]string{"type": "Boolean"}
  case "date":     
  	return map[string]string{"type": "Date"}
  case "map":
  case "object":   
  	return map[string]string{"type": "Object"}
  case "array":    
  	return map[string]string{"type": "Array"}
  case "geopoint": 
  	return map[string]string{"type": "GeoPoint"}
  case "file":     
  	return map[string]string{"type": "File"}
  }
  return map[string]string{}
}