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
	// "errors"
)

func ObjectCreate(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	className := string(vars["className"])

	if !classNameIsValid(className) {
		err := helpers.RenderJsonErr(w, http.StatusBadRequest, helpers.OBJECT_NOT_FOUND, fmt.Sprintf("Invalid classname: %s, classnames can only have alphanumeric characters and _, and must start with an alpha character ", className))
		if err != nil {
			panic(err)
		}
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	_, err := parseReqBodyParams(body)
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

	_, _ = models.SchemaQuery(bson.M{"_id": className})
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
