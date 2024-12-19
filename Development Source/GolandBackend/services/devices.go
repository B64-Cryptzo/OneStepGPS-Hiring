package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golandbackend/models"
)

func FetchDeviceDataFromOneStepAPI(apiURL, apiKey string) ([]models.Device, error) {
	constructedOneStepAPIUrl := fmt.Sprintf("%s?latest_point=true&api-key=%s", apiURL, apiKey)

	httpClient := &http.Client{Timeout: 10 * time.Second}

	clientResponse, errorLog := httpClient.Get(constructedOneStepAPIUrl)
	if errorLog != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", constructedOneStepAPIUrl, errorLog)
	}

	defer clientResponse.Body.Close()

	responseBody, errorLog := io.ReadAll(clientResponse.Body)
	if errorLog != nil {
		return nil, fmt.Errorf("failed to read response body: %w", errorLog)
	}

	var apiResponseJSON models.APIResponse
	if errorLog := json.Unmarshal(responseBody, &apiResponseJSON); errorLog != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", errorLog)
	}

	return apiResponseJSON.ResultList, nil
}
