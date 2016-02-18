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
	// "net/url"
	// "reflect"
)

func ObjectCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I AM INSIDE CREATE")
	vars := mux.Vars(r)
	className := string(vars["className"])
	if string(vars["objectId"]) != "" {
		helpers.RenderJson(w, http.StatusBadRequest, map[string]string{"error": "Cannot create a specific object id"})
		return
	}
	// Return error if the className is not valid
	if !classNameIsValid(className) {
		err := helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.OBJECT_NOT_FOUND, fmt.Sprintf("Invalid classname: %s, classnames can only have alphanumeric characters and _, and must start with an alpha character ", className))
		if err != nil {
			panic(err)
		}
		return
	}

	// parse body into []byte
	body, _ := ioutil.ReadAll(r.Body)
	// parseReqBodyParams ensures that all fields are valid (err equals nil)
	// and return error if json -> map conversion returns error
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
					fieldType := "unidentified type"
					helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.INCORRECT_TYPE, fmt.Sprintf("invalid type for key %s, expected %s, but got %s", fieldName, expectedFieldType, fieldType))
					return
				case bool:
					fieldType := "boolean"
					// if expectedFieldType and actual fieldType are the same, we will est the object map's value using the fieldName
					// else return error
					// same logic in bool, string, number, and array
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
						// this function sets the various special fieldType fields in the object
						// types: Object, Date, GeoPoint
						errHash := setSpecialFieldTypeFields(object, fieldName, v, fieldType)
						if errHash != nil {
							helpers.RenderJson(w, http.StatusBadRequest, errHash)
							return
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

				// when fieldName does not exist in current schema
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
					// this function sets the various special fieldType fields in the object
					// types: Object, Date, GeoPoint
					errHash := setSpecialFieldTypeFields(object, fieldName, v, fieldType)
					if errHash != nil {
						helpers.RenderJson(w, http.StatusBadRequest, errHash)
						return
					}
					schemaUpdates[fieldName] = fieldType
				case []interface{}:
					object[fieldName] = v
					schemaUpdates[fieldName] = "array"
				case nil:
					// when we get null, save it regardless, just don't record it in the schema
					object[fieldName] = v
				}
			}

		}
		fmt.Println(object)

		// if schemaUpdates is larger than 0, then we will update schema
		if len(schemaUpdates) > 0 {
			err := models.SchemaUpdate(schemaUpdates, className)
			if err != nil {
				panic(err)
				return
			}
		}
	} else {

		// this block of code executes when the schema for the className does not exist
		// again, className and all fieldNames can be assumed legal
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
				// this function sets the various special fieldType fields in the object
				// types: Object, Date, GeoPoint
				errHash := setSpecialFieldTypeFields(object, fieldName, v, fieldType)
				if errHash != nil {
					helpers.RenderJson(w, http.StatusBadRequest, errHash)
					return
				}
				schema[fieldName] = fieldType

			case []interface{}:
				object[fieldName] = v
				schema[fieldName] = "array"
			case nil:
				object[fieldName] = v
				// do nothing when it is nil, this is a placeholder case to prevent error
			}
		}

		err := models.SchemaCreate(schema, className)
		// there is nothing we can do if the schema fails to create when it passes all checks
		// TODOs: implement a failsafe for this
		if err != nil {
			panic(err)
			return
		}

	}

	// we can only panic and return if the model fails to create
	// TODOs: implement a failsafe for this
	err = models.ObjectCreate(object, className)
	if err != nil {
		panic(err)
	}
	// by now, the _created_at and _id fields should have been updated. All go maps are passed by reference
	helpers.RenderJson(w, http.StatusOK, map[string]interface{}{"createdAt": object["_created_at"], "objectId": object["_id"]})
	return

}

func ObjectUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	className := string(vars["className"])
	objectId := string(vars["objectId"])

	if className == "" {
		err := helpers.RenderJson(w, http.StatusBadRequest, map[string]string{"error": "missing class name"})
		if err != nil {
			panic(err)
		}
		return
	} else if objectId == "" {
		err := helpers.RenderJson(w, http.StatusBadRequest, map[string]string{"error": "Cannot update without a specific object id"})
		if err != nil {
			panic(err)
		}
		return
	}
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
						// this function sets the various special fieldType fields in the object
						// types: Object, Date, GeoPoint
						errHash := setSpecialFieldTypeFields(object, fieldName, v, fieldType)
						if errHash != nil {
							helpers.RenderJson(w, http.StatusBadRequest, errHash)
							return
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

					// key checking sequence "__type", "__op"
					// this function sets the various special fieldType fields in the object
					// types: Object, Date, GeoPoint
					errHash := setSpecialFieldTypeFields(object, fieldName, v, fieldType)
					if errHash != nil {
						helpers.RenderJson(w, http.StatusBadRequest, errHash)
						return
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

	// try to update object; return error if not successful
	object["_updated_at"] = time.Now()
	err = models.ObjectUpdate(object, objectId, className)
	if err != nil {
		_ = helpers.RenderJsonErr(w, http.StatusNotFound, helpers.OBJECT_NOT_FOUND, "object not found for update")
		return
	}

	helpers.RenderJson(w, http.StatusOK, map[string]interface{}{"updatedAt": object["_updated_at"]})
	return

}

// still need to write parser to match special fields: Object, GeoPoint, Date
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

	// check body first to see if queries are found
	body, _ := ioutil.ReadAll(r.Body)
	queries, err := parseBodyQueryParams(body)
	fmt.Println(queries)
	if err != nil {
		// print some errors here
		return
	}
	// Second: check UrlEncodedQuery
	// parseUrlEncodedQueryParams(r.URL.RawQuery)
	// queries, _ := url.ParseQuery(r.URL.RawQuery)

	// fmt.Println(queries["where"])
	// fmt.Println(len(queries["where"]))
	// fmt.Println(queries["limit"])
	// fmt.Println(len(queries["limit"]))
	// fmt.Println(queries["skip"])
	// fmt.Println(len(queries["skip"]))
	// fmt.Println(queries["keys"])
	// fmt.Println(len(queries["keys"]))
	// fmt.Println(queries["order"])
	// fmt.Println(len(queries["order"]))

	vars := mux.Vars(r)
	className := string(vars["className"])
	objects, _ := models.QueryObject(map[string]interface{}{}, className)
	// body, _ := ioutil.ReadAll(r.Body)
	// fmt.Println(r)
	// fmt.Println(string(body))
	_ = helpers.RenderJson(w, http.StatusOK, objects)
	return
}

func ObjectDestroy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	className := string(vars["className"])
	objectId := string(vars["objectId"])
	if className == "" {
		err := helpers.RenderJson(w, http.StatusBadRequest, map[string]string{"error": "missing class name"})
		if err != nil {
			panic(err)
		}
		return
	} else if objectId == "" {
		err := helpers.RenderJson(w, http.StatusBadRequest, map[string]string{"error": "Cannot delete without a specific object id"})
		if err != nil {
			panic(err)
		}
		return
	}

	err := models.ObjectDestroy(objectId, className)

	if err != nil {
		// If we didn't find it, 404
		err := helpers.RenderJsonErr(w, http.StatusNotFound, helpers.OBJECT_NOT_FOUND, "object not found for delete")
		if err != nil {
			panic(err)
		}

		return
	}

	// If the destroy action has no error, return the success message
	err = helpers.RenderJson(w, http.StatusOK, map[string]string{})
	if err != nil {
		panic(err)
	}

}

// this function handles the special fieldTypes: Object, GeoPoint, Date (Pointer and Relation to be considered)
func setSpecialFieldTypeFields(object map[string]interface{}, fieldName string, fieldValue map[string]interface{}, fieldType string) map[string]interface{} {
	switch fieldType {
	default:
		object[fieldName] = fieldValue
	case "geopoint":
		geoPoint, err := parseGeoPoint(fieldValue)
		object[fieldName] = geoPoint
		if err != nil {
			return map[string]interface{}{"code": helpers.INCORRECT_TYPE, "error": err.Error()}
		}
	case "date":
		dateObject, err := parseDate(fieldValue)
		if err != nil {
			return map[string]interface{}{"code": helpers.INCORRECT_TYPE, "error": err.Error()}
		}
		object[fieldName] = dateObject
	}
	return nil
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
func parseGeoPoint(geoPoint map[string]interface{}) ([]interface{}, error) {

	values := make([]interface{}, 2, 2)
	for _, key := range []string{"latitude", "longitude"} {
		switch geoPoint[key].(type) {
		default:
			return values, errors.New(fmt.Sprintf("Wrong format code 107: TODO: %v", geoPoint[key]))
		case int, int32, int64, float32, float64:
			value := geoPoint[key].(float64)
			switch key {
			case "latitude":
				if float64(value) <= float64(90) && float64(value) >= float64(-90) {
					// we want the original item when storing the point
					values[0] = geoPoint[key]
				} else {
					return values, errors.New(fmt.Sprintf("Latitude must be in [-90, 90]: %v", geoPoint[key]))
				}
			case "longitude":
				if float64(value) < float64(180) && float64(value) >= float64(-180) {
					// we want the original item when storing the point
					values[1] = geoPoint[key]
				} else {
					return values, errors.New(fmt.Sprintf("Longitude must be in [-180, 180): %v", geoPoint[key]))
				}
			}
		}
	}
	return values, nil
}

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
