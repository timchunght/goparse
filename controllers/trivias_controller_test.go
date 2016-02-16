package controllers_test

import (
	"bytes"
	"encoding/json"
	"goparse"
	"goparse/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"goparse/connection"
	"goparse/controllers"
	"goparse/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// this function sets up the test dependencies like connecting to database and setting environment variables
func init() {
	os.Setenv("MONGO_URL", "mongodb://localhost:27017/modernplanit_test")
	os.Setenv("PORT", "8080")
	connection.Connect()
}

func TestTriviaCreateControllerMethodShouldReturnHttpStatusCreated(t *testing.T) {

	params := map[string]interface{}{"name": "What is your name", "description": "The first question", "event_id": "56b01e2520dba346eb1932fc"}
	object, _ := json.Marshal(params)
	request, _ := http.NewRequest("POST", "/trivias", bytes.NewReader(object))
	response := httptest.NewRecorder()
	controllers.TriviaCreate(response, request)

	assert.Equal(t, 201, response.Code, "They should be equal")
}

func TestTriviaCreateReturn400WhenRequiredParamsNotPresent(t *testing.T) {

	object, _ := json.Marshal(map[string]interface{}{"name": "Is tim the best?"})
	request, _ := http.NewRequest("POST", "/trivias", bytes.NewReader(object))
	response := httptest.NewRecorder()
	controllers.TriviaCreate(response, request)

	assert.Equal(t, 400, response.Code, "They should be equal")
}

func TestTriviaCreateReturn400WhenRequiredParamKeyPresentButValueBlank(t *testing.T) {

	object, _ := json.Marshal(map[string]interface{}{"name": "", "description": "Awesome blank question", "event_id": "56b01e2520dba346eb1932fc"})
	request, _ := http.NewRequest("POST", "/trivias", bytes.NewReader(object))
	response := httptest.NewRecorder()
	controllers.ContactCreate(response, request)

	assert.Equal(t, 400, response.Code, "They should be equal")
}

func TestTriviaCreateReturnsCorrectJsonResponse(t *testing.T) {

	params := map[string]interface{}{"name": "What is your name", "description": "The first question", "event_id": "56b01e2520dba346eb1932fc"}
	object, _ := json.Marshal(params)
	request, _ := http.NewRequest("POST", "/trivias", bytes.NewReader(object))
	response := httptest.NewRecorder()
	controllers.TriviaCreate(response, request)
	json := test.MapFromJSON(response.Body.Bytes())
	assert.Equal(t, params["name"], json["name"], "They should be equal")
	assert.Equal(t, params["description"], json["description"], "They should be equal")
	assert.Equal(t, params["event_id"], json["event_id"], "They should be equal")

}

func TestTriviaCreateReturnsJsonContentTypeHeader(t *testing.T) {

	request, _ := http.NewRequest("POST", "/trivias", bytes.NewReader([]byte("")))
	response := httptest.NewRecorder()
	controllers.TriviaCreate(response, request)
	content_type := response.HeaderMap["Content-Type"][0]
	assert.Equal(t, "application/json; charset=UTF-8", content_type, "They should be equal")

}

func TestTriviaShowReturnsCorrectJson(t *testing.T) {
	params := map[string]interface{}{"name": "What is your name", "description": "The first question", "event_id": "56b01e2520dba346eb1932fc"}
	id := createTrivia(params)
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	url := "http://localhost:8080/trivias/" + id
	// retrieve contact
	request, _ := http.NewRequest("GET", url, nil)
	// response = httptest.NewRecorder()
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve object into hash map
	json := test.MapFromJSON(body)
	assert.Equal(t, id, json["id"], "They should be equal")
	assert.Equal(t, params["name"], json["name"], "They should be equal")
	assert.Equal(t, params["description"], json["description"], "They should be equal")
	assert.Equal(t, params["event_id"], json["event_id"], "They should be equal")

}

func TestTriviaShowReturnsErrorWhenWrongId(t *testing.T) {

	// start server
	router := main.NewRouter()
	go http.ListenAndServe(":8080", router)

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	url := "http://localhost:8080/trivias/fakeid"
	// retrieve contact
	request, _ := http.NewRequest("GET", url, nil)
	// response = httptest.NewRecorder()
	res, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, 404, res.StatusCode, "They should be equal")

}
func TestTriviaQueryReturnsListOfObjects(t *testing.T) {
	var results []map[string]interface{}
	params := map[string]interface{}{"name": "What is your name", "description": "The first question", "event_id": "56b01e2520dba346eb1932fc"}
	// create contact and retrieve contactId
	_ = createTrivia(params)
	url := "http://localhost:8080/trivias?event_id=" + params["event_id"].(string)
	request, _ := http.NewRequest("GET", url, nil)
	// create a client and make the request
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// parse response body into []byte
	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	_ = json.Unmarshal(body, &results)

	for _, json := range results {
		assert.Equal(t, params["event_id"], json["event_id"], "They should be equal")
	}

}

func TestTriviaQueryReturns404RecordNotFoundIfWrongValueIsGiven(t *testing.T) {

	url := "http://localhost:8080/trivias?event_id=fakeId"
	request, _ := http.NewRequest("GET", url, nil)
	// create a client and make the request
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// parse response body into []byte
	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, 404, response.StatusCode, "They should be equal")
	assert.Equal(t, "Record not found", json["text"], "They should be equal")

}

func TestTriviaQueryReturns404RecordNotFoundIfWrongQueryParameterIsGiven(t *testing.T) {

	url := "http://localhost:8080/trivias?helloworld=fakeId"
	request, _ := http.NewRequest("GET", url, nil)
	// create a client and make the request
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// parse response body into []byte
	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, 400, response.StatusCode, "They should be equal")
	assert.Equal(t, "Missing required query parameter", json["text"], "They should be equal")

}

func TestTriviaDestroyReturnsSuccessWhenTriviaExists(t *testing.T) {

	params := map[string]interface{}{"name": "What is your name", "description": "The first question", "event_id": "56b01e2520dba346eb1932fc"}
	id := createTrivia(params)

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	url := "http://localhost:8080/trivias/" + id
	// retrieve contact
	request, _ := http.NewRequest("DELETE", url, nil)
	// response = httptest.NewRecorder()
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, 200, response.StatusCode, "They should be equal")
	assert.Equal(t, "Successfully deleted", json["message"], "They should be equal")

}

func TestTriviaDestroyReturns404IfTriviaDoesNotExist(t *testing.T) {

	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	fakeId := "56b01e2520dba346eb1932fc"
	url := "http://localhost:8080/trivias/" + fakeId
	request, _ := http.NewRequest("DELETE", url, nil)
	// response = httptest.NewRecorder()
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, 404, response.StatusCode, "They should be equal")
	assert.Equal(t, "not found", json["text"], "They should be equal")

}

func TestTriviaUpdateReturnsCorrectJsonResponse(t *testing.T) {

	params := map[string]interface{}{"name": "What is your name", "description": "The first question", "event_id": "56b01e2520dba346eb1932fc"}
	updateParams := map[string]interface{}{"name": "Where are you from", "description": "Another Interesting question"}

	// create contact and retrieve contactId
	id := createTrivia(params)
	// set update url
	url := "http://localhost:8080/trivias/" + id
	// create http request object
	object, _ := json.Marshal(updateParams)
	request, _ := http.NewRequest("PUT", url, bytes.NewReader(object))
	// create a client and make the request
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// parse response body into []byte
	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, id, json["id"], "They should be equal")
	assert.Equal(t, updateParams["name"], json["name"], "They should be equal")
	assert.Equal(t, updateParams["description"], json["description"], "They should be equal")

}

func TestTriviaUpdateDoesNotRemovedFieldsNotMentioned(t *testing.T) {

	params := map[string]interface{}{"name": "What is your name", "description": "The first question", "event_id": "56b01e2520dba346eb1932fc"}
	updateParams := map[string]interface{}{"name": "Where are you from", "description": "Another Interesting question"}

	// create contact and retrieve contactId
	id := createTrivia(params)
	// set update url
	url := "http://localhost:8080/trivias/" + id
	// create http request object
	object, _ := json.Marshal(updateParams)
	request, _ := http.NewRequest("PUT", url, bytes.NewReader(object))
	// create a client and make the request
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// parse response body into []byte
	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, id, json["id"], "They should be equal")
	assert.Equal(t, updateParams["name"], json["name"], "They should be equal")
	assert.Equal(t, updateParams["description"], json["description"], "They should be equal")
	// this is a field in the create params but not the update params
	assert.Equal(t, params["event_id"], json["event_id"], "They should be equal")

}

func TestTriviaUpdateReturnsErrorWhenInvalidBsonIdIsGiven(t *testing.T) {

	updateParams := map[string]interface{}{"name": "Where are you from?", "description": "Another Interesting question"}

	// set update url
	url := "http://localhost:8080/trivias/fakeId"
	// create http request object
	object, _ := json.Marshal(updateParams)
	request, _ := http.NewRequest("PUT", url, bytes.NewReader(object))
	// create a client and make the request
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// parse response body into []byte
	body, _ := ioutil.ReadAll(response.Body)
	// fmt.Println(string(body))
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, "Invalid id", json["text"], "They should be equal")

	assert.Equal(t, 404, response.StatusCode, "They should be equal")

}

func TestTriviaUpdateReturnsErrorIfRecordDoesNotExist(t *testing.T) {

	updateParams := map[string]interface{}{"name": "Where are you from", "description": "Another Interesting question"}
	// set update url with an Id that satisfies the bson format
	url := "http://localhost:8080/trivias/56b01e2520dba346eb1932f4"
	// create http request object
	object, _ := json.Marshal(updateParams)
	request, _ := http.NewRequest("PUT", url, bytes.NewReader(object))
	// create a client and make the request
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// parse response body into []byte
	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, "not found", json["text"], "They should be equal")

	assert.Equal(t, 404, response.StatusCode, "They should be equal")

}

func TestTriviaUpdateRequiredFieldIsBlank(t *testing.T) {
	errUpdateParams := map[string]interface{}{"name": "", "description": "Another Interesting question"}
	// set update url with an Id that satisfies the bson format
	url := "http://localhost:8080/trivias/56b01e2520dba346eb1932f4"
	// create http request object
	object, _ := json.Marshal(errUpdateParams)
	request, _ := http.NewRequest("PUT", url, bytes.NewReader(object))
	// create a client and make the request
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// parse response body into []byte
	body, _ := ioutil.ReadAll(response.Body)
	// parse the newly retrieve contact into map hash
	json := test.MapFromJSON(body)
	assert.Equal(t, "Required params can't be blank", json["text"], "They should be equal")

	assert.Equal(t, 404, response.StatusCode, "They should be equal")

}

func createTrivia(params map[string]interface{}) string {

	// Create contact
	requestJson, _ := json.Marshal(params)
	request, _ := http.NewRequest("POST", "/trivias", bytes.NewReader(requestJson))
	response := httptest.NewRecorder()
	controllers.TriviaCreate(response, request)
	json := test.MapFromJSON(response.Body.Bytes())

	// start server
	router := main.NewRouter()
	go http.ListenAndServe(":8080", router)

	return json["id"].(string)

}
