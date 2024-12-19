package main

import (
	"golandbackend/api"
	"golandbackend/utils"
	"log"
	"net/http"
)

func main() {
	srpCmd := utils.StartSubServerSRP()
	defer utils.TerminateProcess(srpCmd)

	apiURL := "https://track.onestepgps.com/v3/api/public/device"

	http.Handle("/api/devices", api.DevicesHandler(apiURL))
	http.Handle("/api/authenticate", api.AuthenticateHandler())
	http.Handle("/api/preferences", api.PreferencesHandler())
	http.Handle("/api/update-user-preferences", api.UpdateUserPreferencesHandler())

	log.Println("Starting server on :8080...")
	if errorLog := http.ListenAndServe(":8080", utils.EnableCORS(http.DefaultServeMux)); errorLog != nil {
		log.Fatalf("Server failed: %v", errorLog)
	}
}
