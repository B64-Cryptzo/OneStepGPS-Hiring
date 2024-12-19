package api

import (
	"encoding/json"
	"io"
	"net/http"

	"golandbackend/models"
	"golandbackend/services"
	"golandbackend/utils"
)

func DevicesHandler(apiURL string) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		utils.SetJSONResponseHeader(responseWriter)

		sessionToken := httpRequest.Header.Get("Authorization")
		if sessionToken == "" {
			http.Error(responseWriter, "Session token required", http.StatusUnauthorized)
			return
		}

		userList, errorLog := utils.ReadUserData()
		if utils.HandleError(responseWriter, "Unable to read user data", errorLog, http.StatusInternalServerError) {
			return
		}

		targetUser, errorLog := utils.FindUserBySessionToken(sessionToken, userList)
		if utils.HandleError(responseWriter, "Invalid session token", errorLog, http.StatusUnauthorized) {
			return
		}

		apiKey, errorLog := utils.GetAPIKeyBySessionToken(sessionToken, targetUser)
		if utils.HandleError(responseWriter, "Authorization Error", errorLog, http.StatusUnauthorized) {
			return
		}

		devices, errorLog := services.FetchDeviceDataFromOneStepAPI(apiURL, apiKey)
		if utils.HandleError(responseWriter, "Error fetching data", errorLog, http.StatusInternalServerError) {
			return
		}

		errorLog = utils.RefreshPreferencesForUser(devices, targetUser)
		if utils.HandleError(responseWriter, "Error initializing preferences", errorLog, http.StatusInternalServerError) {
			return
		}

		apiResponse := models.APIResponse{ResultList: devices}
		json.NewEncoder(responseWriter).Encode(apiResponse)
	}
}

func AuthenticateHandler() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var authenticationRequest models.AuthenticationRequest

		errorLog := json.NewDecoder(httpRequest.Body).Decode(&authenticationRequest)
		if utils.HandleError(responseWriter, "Invalid request payload", errorLog, http.StatusBadRequest) {
			return
		}

		userList, errorLog := utils.ReadUserData()
		if utils.HandleError(responseWriter, "Unable to read user data", errorLog, http.StatusInternalServerError) {
			return
		}

		switch authenticationRequest.Identifier {
		case 1:
			utils.SRPGenerateServerEphemeral(responseWriter, authenticationRequest, userList)

		case 3:
			utils.SRPCreateServerSession(responseWriter, authenticationRequest, userList)

		default:
			http.Error(responseWriter, "Invalid identifier", http.StatusBadRequest)
		}
	}
}

func PreferencesHandler() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		utils.SetJSONResponseHeader(responseWriter)

		sessionToken := httpRequest.Header.Get("Authorization")
		if sessionToken == "" {
			http.Error(responseWriter, "Session token required", http.StatusUnauthorized)
			return
		}

		userList, errorLog := utils.ReadUserData()
		if utils.HandleError(responseWriter, "Unable to read user data", errorLog, http.StatusInternalServerError) {
			return
		}

		targetUser, errorLog := utils.FindUserBySessionToken(sessionToken, userList)
		if utils.HandleError(responseWriter, "Invalid session token", errorLog, http.StatusUnauthorized) {
			return
		}

		preferences, ok := targetUser["Preferences"]
		if !ok {
			http.Error(responseWriter, "User's preferences not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(responseWriter).Encode(preferences)
	}
}

func UpdateUserPreferencesHandler() http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {

		if httpRequest.Method != http.MethodPost {
			http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionToken := httpRequest.Header.Get("Authorization")
		if sessionToken == "" {
			http.Error(responseWriter, "Session token required", http.StatusUnauthorized)
			return
		}

		userList, errorLog := utils.ReadUserData()
		if utils.HandleError(responseWriter, "Unable to read user data", errorLog, http.StatusInternalServerError) {
			return
		}

		targetUser, errorLog := utils.FindUserBySessionToken(sessionToken, userList)
		if utils.HandleError(responseWriter, "Invalid session token", errorLog, http.StatusUnauthorized) {
			return
		}

		var updatedPreferences map[string]interface{}
		body, errorLog := io.ReadAll(httpRequest.Body)
		if utils.HandleError(responseWriter, "Invalid request body", errorLog, http.StatusBadRequest) {
			return
		}

		errorLog = json.Unmarshal(body, &updatedPreferences)
		if utils.HandleError(responseWriter, "Invalid JSON format", errorLog, http.StatusBadRequest) {
			return
		}

		errorLog = utils.UpdatePreferencesForUser(targetUser, updatedPreferences)
		if utils.HandleError(responseWriter, "Failed to update preferences", errorLog, http.StatusInternalServerError) {
			return
		}

		errorLog = utils.SaveUserData(userList)
		if utils.HandleError(responseWriter, "Failed to update preferences", errorLog, http.StatusInternalServerError) {
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"message":"Preferences updated successfully"}`))
	})
}
