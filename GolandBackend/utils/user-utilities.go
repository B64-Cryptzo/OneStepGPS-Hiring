package utils

import (
	"errors"
	"fmt"
	"golandbackend/models"
)

func ReadUserData() ([]map[string]interface{}, error) {
	var userList []map[string]interface{}
	err := readJSONFile("mockdb/userdb.json", &userList)
	return userList, err
}

func SaveUserData(userList []map[string]interface{}) error {
	return writeJSONFile("mockdb/userdb.json", userList)
}

func ensurePreferences(user map[string]interface{}) map[string]interface{} {
	preferences, ok := user["Preferences"].(map[string]interface{})
	if !ok {
		preferences = make(map[string]interface{})
		user["Preferences"] = preferences
	}
	return preferences
}

func RefreshPreferencesForUser(deviceList []models.Device, user map[string]interface{}) error {
	preferences := ensurePreferences(user)

	for _, device := range deviceList {
		deviceID := fmt.Sprintf("%v", device.DeviceID)
		if _, exists := preferences[deviceID]; !exists {
			preferences[deviceID] = map[string]interface{}{
				"markerUrl": "https://swiftrix.net/uploads/d296f72785e72e6e278d910fbcaf9176.png",
				"isVisible": true,
			}
		}
	}

	return SaveUserData([]map[string]interface{}{user})
}

func UpdatePreferencesForUser(user map[string]interface{}, updatedPreferences map[string]interface{}) error {
	preferences := ensurePreferences(user)

	for deviceID, newPreference := range updatedPreferences {
		if prefMap, ok := newPreference.(map[string]interface{}); ok {
			preferences[deviceID] = prefMap
		} else {
			return errors.New("invalid preference format")
		}
	}

	return SaveUserData([]map[string]interface{}{user})
}

func findUserByField(field, value string, userList []map[string]interface{}) (map[string]interface{}, error) {
	for _, user := range userList {
		if fieldValue, ok := user[field].(string); ok && fieldValue == value {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with %s '%s' not found", field, value)
}

func FindUserBySessionToken(sessionToken string, userList []map[string]interface{}) (map[string]interface{}, error) {
	return findUserByField("sessionToken", sessionToken, userList)
}

func FindUserByUsername(username string, userList []map[string]interface{}) (map[string]interface{}, error) {
	return findUserByField("username", username, userList)
}

func GetAPIKeyBySessionToken(sessionToken string, user map[string]interface{}) (string, error) {
	if token, ok := user["sessionToken"].(string); ok && token == sessionToken {
		if apiKey, ok := user["apiKey"].(string); ok {
			return apiKey, nil
		}
		return "", errors.New("apiKey is not a string")
	}
	return "", errors.New("invalid session token")
}
