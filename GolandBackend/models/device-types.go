package models

/*
Contains the primary data from the
converted struct for the OneStepGPS
API response. This information is used
to pass the data to the logged in users
of the Vue UI.
*/

type Params struct {
	VIN string `json:"vin"` // Vehicle VIN number (if applicable)
}

type DeviceState struct {
	DriveStatus string  `json:"drive_status"` // Whether vehicle is currently driving (if applicable)
	FuelPercent float64 `json:"fuel_percent"` // Vehicle's current fuel level as a percent (if applicable)
}

type LatestDevicePoint struct {
	Lat         float64     `json:"lat"`          // Geographic information in order to get an exact location (Latitude)
	Lng         float64     `json:"lng"`          // Geographic information in order to get an exact location (Longitude)
	Speed       float64     `json:"speed"`        // Normalized Vehicle Speed to help with management (if applicable)
	DeviceState DeviceState `json:"device_state"` // DeviceState contains the drive status and fuel percent
	Params      Params      `json:"params"`       // Params contains the Vehicle VIN number
}

type Device struct {
	DeviceID          string            `json:"device_id"`           // Device ID to idenfity individual devices
	DisplayName       string            `json:"display_name"`        // Associated moniker for the specified device
	ActiveState       string            `json:"active_state"`        // Current status of the specified device
	LatestDevicePoint LatestDevicePoint `json:"latest_device_point"` // LatestDevicePoint contains more specific data we want to expose
}
