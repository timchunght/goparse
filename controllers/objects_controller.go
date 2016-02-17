package controllers

import (
	// "encoding/json"
	"errors"
	"fmt"
	"goparse/Godeps/_workspace/src/github.com/gorilla/mux"
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"goparse/helpers"
	"goparse/models"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
	// "reflect"
)

func ObjectCreate(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.Method)
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
	// GeoPoint
	// ------ TO BE IMPLEMENTED
	// File
	// Pointer
	// Relation
	if classExists {
		schemaUpdates := bson.M{}
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
					if expectedFieldType == fieldType {
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
					fieldType := "unidentified type"
					helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type: %s", fieldType))
					return
				case bool:
					object[fieldName] = v
					schemaUpdates[fieldName] = "boolean"
				case string:
					object[fieldName] = v
					schemaUpdates[fieldName] = "string"
				case int, int32, int64, float32, float64:
					object[fieldName] = v
					schemaUpdates[fieldName] = "number"
				case map[string]interface{}:
					// TODOS: this can be either a Object, Date ("__type"), or GeoPoint,
					fieldType, err := getFieldTypeFromMap(v)
					if err != nil {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, err.Error())
						return
					}
					object[fieldName] = v
					schemaUpdates[fieldName] = fieldType
				case []interface{}:
					object[fieldName] = v
					schemaUpdates[fieldName] = "array"
				case nil:
					object[fieldName] = v
				}
			}

		}
		fmt.Println(object)

		// if schemaUpdates is larger than 0, then we will update schema
		// TODO: Implement Schema Update
		if len(schemaUpdates) > 0 {
			err := models.SchemaUpdate(schemaUpdates, className)
			if err != nil {
				panic(err)
				return
			}
		}
	} else {

		schema = map[string]interface{}{}
		for fieldName, value := range params {

			switch v := value.(type) {
			default:
				fieldType := "unidentified type"
				helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type: %s", fieldType))
				return
			case bool:
				object[fieldName] = v
				schema[fieldName] = "boolean"
			case string:
				object[fieldName] = v
				schema[fieldName] = "string"
			case int, int32, int64, float32, float64:
				object[fieldName] = v
				schema[fieldName] = "number"
			case map[string]interface{}:
				// TODOS: this can be either a Object, Date ("__type"), or GeoPoint,
				fieldType, err := getFieldTypeFromMap(v)
				if err != nil {
					helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, err.Error())
					return
				}
				switch fieldType {
				default:
					object[fieldName] = v
				case "geopoint":
				case "date":
					v, err := parseDate(v)
					if err != nil {
						helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, err.Error())
						return
					}
					object[fieldName] = v
				}

				schema[fieldName] = fieldType

			case []interface{}:
				object[fieldName] = v
				schema[fieldName] = "array"
			case nil:
			}
		}

		err := models.SchemaCreate(schema, className)
		if err != nil {
			panic(err)
		}
		
	}

	err = models.ObjectCreate(object, className)

	if err != nil {
		panic(err)
	}
	helpers.RenderJson(w, http.StatusOK, map[string]interface{}{"createdAt": object["_created_at"], "objectId": object["_id"]})
	return

	// _ = helpers.RenderJson(w, http.StatusOK, schema)

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

func ObjectUpdate(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.Method)
	vars := mux.Vars(r)
	className := string(vars["className"])
	objectId := string(vars["objectId"])
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
	classExists := true
	schema, err := models.SchemaQuery(bson.M{"_id": className})
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
	// GeoPoint
	// ------ TO BE IMPLEMENTED
	// File
	// Pointer
	// Relation
	if classExists {
		schemaUpdates := bson.M{}
		// This block of code assumes that we have the schema object for the collection
		// TODOS: implement the scenario in which the schema for the collection does not exist
		for fieldName, value := range params {

			// if field exists in schema
			if expectedFieldType, ok := schema[fieldName]; ok {
				// We want to make sure that value type matches the type in the schema collection
				switch v := value.(type) {
				default:
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
					if expectedFieldType == fieldType {
						// TODOs:
						// IMPLEMENT THE VARIOUS TYPES "Date", "GeoPoint"
						switch fieldType {
						default:
							object[fieldName] = v
						case "geopoint":
							object[fieldName] = v
						case "date":
							v, err := parseDate(v)
							if err != nil {
								helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, err.Error())
								return
							}
							object[fieldName] = v
						}
						
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
					fieldType := "unidentified type"
					helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type: %s", fieldType))
					return
				case bool:
					object[fieldName] = v
					schemaUpdates[fieldName] = "boolean"
				case string:
					object[fieldName] = v
					schemaUpdates[fieldName] = "string"
				case int, int32, int64, float32, float64:
					object[fieldName] = v
					schemaUpdates[fieldName] = "number"
				case map[string]interface{}:
					// TODOS: this can be either a Object, Date ("__type"), or GeoPoint,
					fieldType, err := getFieldTypeFromMap(v)
					if err != nil {
						_ = helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, err.Error())
						return
					}

					switch fieldType {
					default:
						object[fieldName] = v
					case "geopoint":
						object[fieldName] = v
					case "date":
						v, err := parseDate(v)
						if err != nil {
							helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, err.Error())
							return
						}
						object[fieldName] = v
					}
					schemaUpdates[fieldName] = fieldType
				case []interface{}:
					object[fieldName] = v
					schemaUpdates[fieldName] = "array"
				case nil:
					object[fieldName] = v
				}
			}

		}
		fmt.Println(object)

		// if schemaUpdates is larger than 0, then we will update schema
		// TODO: Implement Schema Update
		if len(schemaUpdates) > 0 {
			err := models.SchemaUpdate(schemaUpdates, className)
			if err != nil {
				panic(err)
				return
			}
		}
	} else {

		_ = helpers.RenderJsonErr(w, http.StatusNotFound, helpers.OBJECT_NOT_FOUND, "object not found for update")
		return
	}

	// try to update object and return error if not successful
	object["_updated_at"] = time.Now()
	err = models.ObjectUpdate(object, objectId, className)
	if err != nil {
		_ = helpers.RenderJsonErr(w, http.StatusNotFound, helpers.OBJECT_NOT_FOUND, "object not found for update")
		return
	}

	helpers.RenderJson(w, http.StatusOK, map[string]interface{}{"updatedAt": object["_updated_at"]})
	return

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

	object, err := models.FindObjectById(objectId, className)
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

func ObjectQuery(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(r)
	fmt.Println(string(body))
	return
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

		switch fieldType {
		default:
			return "", errors.New(fmt.Sprintf("invalid type: %s", fieldType))
		case "Date":
			if _, ok := fieldValue["iso"]; ok {
				// fmt.Println(fieldType)
				return strings.ToLower(fieldType.(string)), nil
			} else {
				return "", errors.New(fmt.Sprintf("Invalid date: %v", fieldValue))
			}
		case "GeoPoint":
			// if the fieldType is GeoPoint and latitude exists, return no error
			// else return an error that corresponds with code 111
			if _, ok := fieldValue["latitude"]; ok {
				if _, ok := fieldValue["longitude"]; ok {
					// fmt.Println(fieldType)
					return strings.ToLower(fieldType.(string)), nil
				} else {
					return "", errors.New("Invalid format for longitude")
				}
			} else {
				return "", errors.New("Invalid format for latitude")
			}
		}
	}

	return "object", nil

}

// parse GeoPoint map and makes sure that latitude is "Latitude must be in [-90, 90]: 123213.0"
// "Longitude must be in [-180, 180): 213213.2" (code 111)
// It also makes sure that the values are numbers (if they are not numbers, code 107)
// func parseGeoPoint(geoPoint map[string]interface{}) ([]float64, error) {

// 	return
// }

// parse Date map and makes sure that iso is a string and it is parsable (code 107 if not)
// "error": "invalid date: '2011-11-07T20:58:34.448Zs'" code 107
// sample date string: 2011-11-07T20:58:34.448Z
// myTime, err := time.Parse(time.RFC3339, "2015-03-13T22:05:08Z")
func parseDate(date map[string]interface{}) (time.Time, error) {
	switch dateString := date["iso"].(type) {
	default:
		return time.Time{}, errors.New("unexpected type of iso")
	case string:
		dateObject, err := time.Parse(time.RFC3339, dateString)
		if err != nil {
			return time.Time{}, errors.New(fmt.Sprintf("invalid date: '%s'", dateString))
		}

		return dateObject, err
	}

}

// ------------------------------LEGACY CODE------------------------------
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
