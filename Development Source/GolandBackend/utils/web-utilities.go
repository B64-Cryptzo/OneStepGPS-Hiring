package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

/*
This is important when testing over a local connection.
Generally, two locally hosted HTTP webservers cannot interact
with each other due to Cross Origin Resource Sharing policy.
We want to bypass this for the purposes of local development.
*/

func EnableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func SetJSONResponseHeader(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
}

func HandleError(responseWriter http.ResponseWriter, message string, errorLog error, statusCode int) bool {
	if errorLog != nil {
		http.Error(responseWriter, message+": "+errorLog.Error(), statusCode)
		return true
	}
	return false
}

func MakeSecuredPostRequest[T any](url string, requestBody map[string]interface{}) (*T, error) {
	jsonBody, _ := json.Marshal(requestBody)

	req, errorLog := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if errorLog != nil {
		return nil, errorLog
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 6IWpWBBtYzjQdCBj")

	client := &http.Client{}
	subServerResponse, errorLog := client.Do(req)
	if errorLog != nil {
		return nil, errorLog
	}
	defer subServerResponse.Body.Close()

	var response T
	if errorLog := json.NewDecoder(subServerResponse.Body).Decode(&response); errorLog != nil {
		return nil, errorLog
	}

	return &response, nil
}
