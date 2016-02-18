package controllers

import (
	// "fmt"
	"net/url"
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

func parseBodyQueryParams(body []byte) (bson.M, error) {
	var params map[string]interface{}

	err := json.Unmarshal(body, &params)
	if err != nil {
		return bson.M{}, err
	}

	for key, _ := range params {
		switch key {
		default:
		case "where":
		case "order":
		case "limit":
		case "skip":
		case "keys":
		case "include":
		}
	}
	return bson.M{}, nil
}

func parseUrlEncodedQueryParams(rawQuery string) (bson.M, error) {
	
	queryMap, _ := url.ParseQuery(rawQuery)
	var query bson.M
	for key, value := range queryMap {
		switch key {
		default:
		case "where":
			if len(value) == 1 {
				
				err := json.Unmarshal([]byte(value[0]), &query)
				if err != nil {
					return bson.M{}, err
				}

			} 
		case "order":
		case "limit":
		case "skip":
		case "keys":
		case "include":
		}
	}

	fmt.Println(query)
	// fmt.Println(query["where"])
	// fmt.Println(len(query["where"]))
	// fmt.Println(query["limit"])
	// fmt.Println(len(query["limit"]))
	// fmt.Println(query["skip"])
	// fmt.Println(len(query["skip"]))
	// fmt.Println(query["keys"])
	// fmt.Println(len(query["keys"]))
	// fmt.Println(query["order"])
	// fmt.Println(len(query["order"]))

	// var params map[string]interface{}

	// err := json.Unmarshal(body, &params)
	// if err != nil {
	// 	return bson.M{}, err
	// }

	
	return query, nil
}



func formatObjectQuery(query bson.M) error {
	// paramKeyMapping := map[string]string{"objectId": "_id", "updatedAt": "_updated_at", "createdAt": "_created_at"}
	// for exposedParamKey, dbParamKey := range paramKeyMapping {
	// 	// we do not allow querying using param key format that we use in database (prefixed with "_")
	// 	delete(query, dbParamKey)

	// 	if value, ok := query[exposedParamKey]; ok {
	// 		query[dbParamKey] = value
	// 		delete(query, exposedParamKey)
	// 	}
	// }

	
	return nil

}