package controllers

import (
	// "encoding/json"
	"goparse/Godeps/_workspace/src/github.com/gorilla/mux"
	"net/http"
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"goparse/helpers"
	"goparse/models"
	"regexp"
	// "time"
	"io/ioutil"
	"fmt"
	"errors"
	// "reflect"
)

func ObjectCreate(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	className := string(vars["className"])

	// Return error if the className is not valid
	if !classNameIsValid(className) {
		err := helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.OBJECT_NOT_FOUND, fmt.Sprintf("Invalid classname: %s, classnames can only have alphanumeric characters and _, and must start with an alpha character ", className))
		if err != nil {
			panic(err)
		}
		return
	}

	// parse body and return error if json -> map conversion returns error
	body, _ := ioutil.ReadAll(r.Body)
	// parseReqBodyParams ensures that all fields are valid (err equals nil)
	params, err := parseReqBodyParams(body)
	if err != nil {

		if err.Error() == "invalid JSON" {
			err := helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INVALID_JSON, err.Error())
			if err != nil {
				panic(err)
			}
		} else {
			err := helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INVALID_KEY_NAME, err.Error())
			if err != nil {
				panic(err)
			}
		}
		
		return

	}


	
	// if it reaches this stage, it means that both the className and fieldNames are legal
	object := map[string]interface{}{}
	schema, err := models.SchemaQuery(bson.M{"_id": className})
	classExists := true
	if err != nil {
		if err.Error() == "not found" {

			classExists = false
		} else {
			panic(err)
			return
		}
	}

	// LEGAL TYPES:
	// Boolean
	// String
	// Number
	// Date
	// Object
	// Array
	// ------ TO BE IMPLEMENTED
	// GeoPoint
	// File
	// Pointer
	// Relation
	if classExists {
		schemaUpdate := bson.M{}
		// This block of code assumes that we have the schema object for the collection
		// TODOS: implement the scenario in which the schema for the collection does not exist
		for fieldName, value := range params {
			
			// if field exists in schema
			if expectedFieldType, ok := schema[fieldName]; ok {

				// We want to make sure that value type matches the type in the schema collection
				switch v := value.(type) {
				default:
					// TODO:
					// RETURN UNIDENTIFIED JSON ERR
					// fmt.Println("unidentified")
					// fmt.Println(v)

					fieldType := "unidentified type"
					helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type for key %s, expected %s, but got %s", fieldName, expectedFieldType, fieldType))
					return
				case bool:
					fieldType := "boolean"
					if expectedFieldType == fieldType {
						object[fieldName] = v
					} else {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type for key %s, expected %s, but got %s", fieldName, expectedFieldType, fieldType))
						return
					}
				case string:
					fieldType := "string"
					if expectedFieldType == fieldType {
						object[fieldName] = v
					} else {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type for key %s, expected %s, but got %s", fieldName, expectedFieldType, fieldType))
						return
					}
				case int, int32, int64, float32, float64:
					fieldType := "number"
					if expectedFieldType == fieldType {
						object[fieldName] = v
					} else {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type for key %s, expected %s, but got %s", fieldName, expectedFieldType, fieldType))
						return
					}
				case map[string]interface{}:
					// TODOS: this can be either a Object, Date ("__type"), or GeoPoint, 
					fieldType, err := getFieldTypeFromMap(v)
					// if the fieldType is not a legal type, return error
					if err != nil {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, err.Error())
						return
					}

					// if the fieldType is a legalType but does not match the type in the schema, return error
					if expectedFieldType ==  fieldType {
						object[fieldName] = v
					} else {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type for key %s, expected %s, but got %s", fieldName, expectedFieldType, fieldType))
						return
					}
				case []interface{}:
					fieldType := "array"
					if expectedFieldType == fieldType {
						object[fieldName] = v
					} else {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type for key %s, expected %s, but got %s", fieldName, expectedFieldType, fieldType))
						return
					}
				case nil:
					object[fieldName] = v
				}
			} else {
				

				switch v := value.(type) {
				default:
					// TODOs: return json error
					fmt.Println("unidentified")
					fmt.Println(v)
				case bool:
					object[fieldName] = v
					schemaUpdate[fieldName] = "boolean"
				case string:
					object[fieldName] = v
					schemaUpdate[fieldName] = "string"
				case int, int32, int64, float32, float64:
					object[fieldName] = v
					schemaUpdate[fieldName] = "number"
				case map[string]interface{}:
					// TODOS: this can be either a Object, Date ("__type"), or GeoPoint, 
					fieldType, err := getFieldTypeFromMap(v)
					if err != nil {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, err.Error())
						return
					}
					object[fieldName] = v
					schemaUpdate[fieldName] = fieldType
				case []interface{}:
					object[fieldName] = v
					schemaUpdate[fieldName] = "array"
				case nil:
					object[fieldName] = v
				}
			}
			
	  }
	  fmt.Println(object)
	  
	  // if schemaUpdate is larger than 0, then we will update schema
	  if len(schemaUpdate) > 0 {
	  	fmt.Println(schemaUpdate)
	  }
	}
  

	_ = helpers.RenderJson(w, http.StatusOK, schema)
	


	// _, paramsPresent := requiredBodyParamsCheck(body, []string{"event_id", "name", "description"})
	// if paramsPresent == true {

	// 	err := json.Unmarshal(body, &trivia)
	// 	if err != nil {
	// 		err := helpers.RenderJsonErr(w, http.StatusBadRequest, "Error parsing params")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		return
	// 	}

	// 	// make sure the create method does not have error
	// 	err = (&trivia).Create()
	// 	if err != nil {
	// 		err := helpers.RenderJsonErr(w, http.StatusInternalServerError, "Error creating object")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		return
	// 	}

	// 	if err := helpers.RenderJson(w, http.StatusCreated, trivia); err != nil {
	// 		panic(err)
	// 	}

	// } else {
	// 	err := helpers.RenderJsonErr(w, http.StatusBadRequest, "Required params not found")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}

func ObjectShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	className := string(vars["className"])
	objectId := string(vars["objectId"])

	if !classNameIsValid(className) {
		err := helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.OBJECT_NOT_FOUND, fmt.Sprintf("Invalid classname: %s, classnames can only have alphanumeric characters and _, and must start with an alpha character ", className))
		if err != nil {
			panic(err)
		}
		return
	}

	object, err := models.FindObjectById(className, objectId)
	if err != nil {
		// If we didn't find it, 404
		err := helpers.RenderJsonErr(w, http.StatusNotFound, helpers.INVALID_CLASS_NAME, err.Error())
		if err != nil {
			panic(err)
		}
		return
	}

	err = helpers.RenderJson(w, http.StatusOK, object)
	if err != nil {
		panic(err)
		return
	}

}

// This function makes sure that the className is valid for all ReadWrite requests
func classNameIsValid(className string) bool {
	//Class names have the same constraints as field names, but also allow the previous additional names.
  return className == "_User" || className == "_Installation" || className == "_Session" || className == "_Role" || fieldNameIsValid(className)
  // TODO: Implement joinClassRegex
  // || joinClassRegex.test(className)
    
}

// Makes sure that the fieldName is legal using the regex
func fieldNameIsValid(fieldName string) bool {

	re, _ := regexp.Compile(`^[A-Za-z][A-Za-z0-9_]*$`)
	return re.Match([]byte(fieldName))
}

// This function returns the field type when the request contains a map
func getFieldTypeFromMap(fieldValue map[string]interface{}) (string, error) {

	if fieldType, ok := fieldValue["__type"]; ok {
		if fieldType == "Date" || fieldType == "GeoPoint" {
			return fieldType.(string), nil
		} else {
			return "", errors.New(fmt.Sprintf("invalid type: %s", fieldType))
		}
	} else {
		return "object", nil
	}
}
// func ObjectQuery(w http.ResponseWriter, r *http.Request) {

// 	query, err := getQuery([]string{"event_id"}, r)

// 	if err != nil {
// 		// If we encounter error parsing query params
// 		err = helpers.RenderJsonErr(w, http.StatusBadRequest, err.Error())
// 		if err != nil {
// 			panic(err)
// 		}

// 		return
// 	}

// 	trivias, err := models.Object{}.Query(query)
// 	if err != nil {
// 		// If we encounter error parsing query params
// 		err = helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())
// 		if err != nil {
// 			panic(err)
// 		}

// 		return
// 	}

// 	err = helpers.RenderJson(w, http.StatusOK, trivias)
// 	if err != nil {
// 		panic(err)
// 	}

// }

// func ObjectDestroy(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := string(vars["id"])
// 	err := models.Object{}.Destroy(id)

// 	if err != nil {
// 		// If we didn't find it, 404
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

// func ObjectUpdate(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := string(vars["id"])
// 	body, _ := ioutil.ReadAll(r.Body)

// 	paramKeys := []string{"event_id", "name", "description"}
// 	requiredParamKeys := paramKeys
// 	doc, err := getUpdateDocFromBody(body, requiredParamKeys, paramKeys)
// 	if err != nil {
// 		// If we didn't find it, 404
// 		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())

// 		if err != nil {
// 			panic(err)
// 		}

// 		return
// 	}

// 	trivia, err := models.Object{}.Update(id, doc)

// 	if err != nil {

// 		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())

// 		if err != nil {
// 			panic(err)
// 		}
// 		return
// 	}

// 	helpers.RenderJson(w, http.StatusOK, trivia)

// }
