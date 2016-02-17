package controllers

import (
	// "fmt"

	"encoding/json"
	"errors"
	"goparse/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"net/http"
	"fmt"
)

func getQuery(paramKeys []string, r *http.Request) (bson.M, error) {

	query := bson.M{}
	// fmt.Println(r.URL.Query())
	for _, paramKey := range paramKeys {

		if value, ok := r.URL.Query()[paramKey]; ok {
			if value[0] == "" {
				// return error if the value is blank
				return query, errors.New("Missing required query parameter")
			} else {

				query[paramKey] = value[0]
			}

		} else {
			// return false if the key does not exist
			return query, errors.New("Missing required query parameter")
		}
	}

	return query, nil

}

func getUpdateDocFromBody(body []byte, requiredParamKeys, paramKeys []string) (bson.M, error) {
	var params map[string]interface{}
	err := json.Unmarshal(body, &params)
	doc := bson.M{}
	updates := bson.M{}
	if err != nil {
		// fmt.Println(err)
		return doc, err
	}

	// if required param key is present but has blank value, return error
	for _, paramKey := range requiredParamKeys {
		if value, ok := params[paramKey]; ok {

			switch v := value.(type) {
			default:
				updates[paramKey] = v
			case string:
				if v == "" {
					// return false if the value is blank
					return doc, errors.New("Required params can't be blank")
				} else {
					updates[paramKey] = v
				}
			case []string:
				if len(v) == 0 {
					return doc, errors.New("Required params can't be blank")
				} else {
					updates[paramKey] = v
				}
			}

		}
	}

	//if all required params are checked, create update doc using the values
	for _, paramKey := range paramKeys {
		if value, ok := params[paramKey]; ok {
			updates[paramKey] = value
		}
	}

	doc["$set"] = updates
	return doc, err
}

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
