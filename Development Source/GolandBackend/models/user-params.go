package models

/*
Simple modifiable struct to allow us to pass any
number of important preferences about any user
or device.
*/
type UserPreferences struct {
	DeviceID  string `json:"device_id"`
	IsVisible bool   `json:"is_visible"`
	MarkerURL string `json:"marker_url"`
}
