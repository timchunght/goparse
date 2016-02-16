package controllers

import (
	"encoding/json"
	"modernplanit/Godeps/_workspace/src/github.com/gorilla/mux"
	"net/http"
	// "modernplanit/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"goparse/helpers"
	"goparse/models"
	// "time"
	"io/ioutil"
	// "fmt"
	// "errors"
)

func TriviaCreate(w http.ResponseWriter, r *http.Request) {
	var trivia models.Trivia
	body, _ := ioutil.ReadAll(r.Body)
	_, paramsPresent := requiredBodyParamsCheck(body, []string{"event_id", "name", "description"})
	if paramsPresent == true {

		err := json.Unmarshal(body, &trivia)
		if err != nil {
			err := helpers.RenderJsonErr(w, http.StatusBadRequest, "Error parsing params")
			if err != nil {
				panic(err)
			}
			return
		}

		// make sure the create method does not have error
		err = (&trivia).Create()
		if err != nil {
			err := helpers.RenderJsonErr(w, http.StatusInternalServerError, "Error creating object")
			if err != nil {
				panic(err)
			}
			return
		}

		if err := helpers.RenderJson(w, http.StatusCreated, trivia); err != nil {
			panic(err)
		}

	} else {
		err := helpers.RenderJsonErr(w, http.StatusBadRequest, "Required params not found")
		if err != nil {
			panic(err)
		}
	}
}

func TriviaShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := string(vars["id"])
	trivia, err := models.Trivia{}.Find(id)
	if err != nil {
		// If we didn't find it, 404
		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())
		if err != nil {
			panic(err)
		}

		return
	}

	err = helpers.RenderJson(w, http.StatusOK, trivia)
	if err != nil {
		panic(err)
	}

}

func TriviaQuery(w http.ResponseWriter, r *http.Request) {

	query, err := getQuery([]string{"event_id"}, r)

	if err != nil {
		// If we encounter error parsing query params
		err = helpers.RenderJsonErr(w, http.StatusBadRequest, err.Error())
		if err != nil {
			panic(err)
		}

		return
	}

	trivias, err := models.Trivia{}.Query(query)
	if err != nil {
		// If we encounter error parsing query params
		err = helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())
		if err != nil {
			panic(err)
		}

		return
	}

	err = helpers.RenderJson(w, http.StatusOK, trivias)
	if err != nil {
		panic(err)
	}

}

func TriviaDestroy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := string(vars["id"])
	err := models.Trivia{}.Destroy(id)

	if err != nil {
		// If we didn't find it, 404
		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())
		if err != nil {
			panic(err)
		}

		return
	}

	// If the destroy action has no error, return the success message
	err = helpers.RenderJson(w, http.StatusOK, map[string]string{"message": "Successfully deleted"})
	if err != nil {
		panic(err)
	}

}

func TriviaUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := string(vars["id"])
	body, _ := ioutil.ReadAll(r.Body)

	paramKeys := []string{"event_id", "name", "description"}
	requiredParamKeys := paramKeys
	doc, err := getUpdateDocFromBody(body, requiredParamKeys, paramKeys)
	if err != nil {
		// If we didn't find it, 404
		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())

		if err != nil {
			panic(err)
		}

		return
	}

	trivia, err := models.Trivia{}.Update(id, doc)

	if err != nil {

		err := helpers.RenderJsonErr(w, http.StatusNotFound, err.Error())

		if err != nil {
			panic(err)
		}
		return
	}

	helpers.RenderJson(w, http.StatusOK, trivia)

}
