package controllers

import (
	// "fmt"
	"net/url"
	"encoding/json"
	"errors"
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	// "net/http"
	"fmt"
	"goparse/helpers"
)

// this function makes sure that all the params are present and their values not blank
func parseReqBodyParams(body []byte) (map[string]interface{}, error) {

	var params map[string]interface{}

	err := json.Unmarshal(body, &params)

	// fmt.Println(params)
	if err != nil {
		// fmt.Println(err)
		return params, errors.New("invalid JSON")
	}


	for key, _ := range params {
		if(!fieldNameIsValid(key)) {
			return params, errors.New(fmt.Sprintf("invalid field name: %s", key))
		}

		if key == "updatedAt" || key == "createdAt" || key == "objectId" {
			return params, errors.New(fmt.Sprintf("%s is an invalid field name", key))
		}
  }
	

	return params, nil
}

// convert the url-encoded query params into bson query object or return an errorMap if it errors out
func parseUrlEncodedQueryParams(rawQuery string) (bson.M, map[string]interface{}) {
	
	queryMap, _ := url.ParseQuery(rawQuery)
	var query map[string]interface{}
	for key, value := range queryMap {
		switch key {
		default:
			return bson.M{}, nil
		case "where":
			if len(value) == 1 {
				
				err := json.Unmarshal([]byte(value[0]), &query)
				if err != nil {
					return bson.M{}, map[string]interface{}{"code": helpers.INVALID_JSON, "error": "invalid JSON"}
				}

			} 
		case "order":
		case "limit":
		case "skip":
		case "keys":
		case "include":
		}
	}
	// fmt.Println(findObjectWithKey(query, "$select"))
	_ = parseWhereQuery(query)
	errMap := formatObjectQuery(query)
	// map can be used as bson.M for return
	return query, errMap
}

// convert the url-encoded query params into bson query object or return an errorMap if it errors out
func parseBodyQueryParams(body []byte) (bson.M, map[string]interface{}) {
	// return an empty bson hash map if the body is empty
	if string(body) == "" {
		
		return bson.M{}, nil
	} else {
		var queryMap map[string]interface{}

		err := json.Unmarshal(body, &queryMap)
		if err != nil {
			return bson.M{}, map[string]interface{}{"code": helpers.INVALID_JSON, "error": "invalid JSON"}
		}
		var query map[string]interface{}
		for key, value := range queryMap {
			switch key {
			default:
				return bson.M{}, nil
			case "where":
				query = value.(map[string]interface{})
			case "order":
			case "limit":
			case "skip":
			case "keys":
			case "include":
			}
		}

		// findObjectWithKey(query, "$select").(map[string]interface{})["helloworld"] = "helloworld"
		// fmt.Println(query)
		parseWhereQuery(query)
		errMap := formatObjectQuery(query)
		
		return query, errMap
	}
	
}

func parseWhereQuery(query map[string]interface{}) map[string]interface{} {

	// check each key to see if it is valid
	// if the first level value of a key is a map and the keys are not part an action key "$action"
	// do inner value checking
	// fmt.Println("INSIDE parseWhereQuery")
	for fieldName, value := range query {
		if fieldNameIsValid(fieldName) {
			switch value.(type) {
			default:
			case map[string]interface{}:
				// fmt.Println("MAP QUERY")
				for innerFieldName, innerValue := range value.(map[string]interface{}) {
					if isMongoQueryActionKey(innerFieldName) {
						switch innerValue.(type) {
						case map[string]interface{}:
							if innerValue.(map[string]interface{})["__type"] == "Date" {
								innerValue, _ = parseDate(innerValue.(map[string]interface{}))
							}
						}
					} else {
						// if innerFieldName is not an action and the value has field __type
						// if _, ok := value.(map[string]interface{})["__type"]; ok {
						switch value.(map[string]interface{})["__type"] {
						case "Date":
							dateObject, _ := parseDate(value.(map[string]interface{}))
							query[fieldName] = dateObject
						}
						// }

						break
					}
					// fmt.Println(innerFieldName)
					// fmt.Println(innerValue)
				}
			}

		} else {
			return map[string]interface{}{"code": helpers.INVALID_QUERY, "error": fmt.Sprintf("Invalid key %s for find", fieldName)}
		}
	}

	fmt.Println(query)
	return nil
}

func formatObjectQuery(query bson.M) map[string]interface{} {
	

	// make sure that the keys of the query keys are valid
	for exposedParamKey, _ := range query {
		if(!fieldNameIsValid(exposedParamKey)) {
			return map[string]interface{}{"code": helpers.INVALID_KEY_NAME, "error": fmt.Sprintf("invalid field name: %s", exposedParamKey)}
		}
	}

	paramKeyMapping := map[string]string{"objectId": "_id", "updatedAt": "_updated_at", "createdAt": "_created_at"}
	for exposedParamKey, dbParamKey := range paramKeyMapping {
		// we do not allow querying using param key format that we use in database (prefixed with "_")
		if value, ok := query[exposedParamKey]; ok {
			query[dbParamKey] = value
			delete(query, exposedParamKey)
		}
	}
	return nil

}

func isMongoQueryActionKey(action string) bool {
	// Key	Operation
	// --------------
	// $lt	Less Than
	// $lte	Less Than Or Equal To
	// $gt	Greater Than
	// $gte	Greater Than Or Equal To
	// $ne	Not Equal To
	// $in	Contained In
	// $nin	Not Contained in
	// $exists	A value is set for the key
	// $select	This matches a value for a key in the result of a different query
	// $dontSelect	Requires that a key's value not match a value for a key in the result of a different query
	// $all	Contains all of the given values
	// $regex	Requires that a key's value match a regular expression
	actions := []string{"$lt", "$lte", "$gt", "$gte", "$ne", "$in", "$nin", "$exists", "$select", "$dontSelect", "$all", "$regex"}
	for _, value := range actions {
		if action == value {
			return true
		}
	}
	return false
}



// Finds a subobject that has the given key, if there is one.
// Returns nil otherwise.
// this only accepts 
func findObjectWithKey(root interface{}, key string) interface{} {
  // if (typeof root !== 'object') {
  //   return;
  // }

  switch root.(type) {
  default:
  	return nil
  case []interface{}:
  	for _, item := range root.([]interface{}) {
  		answer := findObjectWithKey(item, key)
  		if answer != nil {
  			return answer
  		}
  	}
  case map[string]interface{}:
  }
  
  if _, ok := root.(map[string]interface{})[key]; ok {
		
		return root
	}

	for subkey, _ := range root.(map[string]interface{}) {
		answer := findObjectWithKey(root.(map[string]interface{})[subkey], key)
		if answer != nil {
			return answer
		}
	}
  
  return nil
}