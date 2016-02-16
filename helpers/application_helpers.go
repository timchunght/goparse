package helpers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type JsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

//Wrapper function for RenderJson
func RenderJsonErr(w http.ResponseWriter, statusCode int, text string) error {
	return RenderJson(w, statusCode, JsonErr{Code: statusCode, Text: text})
}

func RenderJson(w http.ResponseWriter, statusCode int, object interface{}) error {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(object)
	return err

}
