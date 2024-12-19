package models

/*
A simple struct to help pass JSON data effectively
via HTTP protocol.
*/
type APIResponse struct {
	ResultList []Device `json:"result_list"`
}

type ServerEphemeralResponse struct {
	ServerPublicEphemeral string `json:"serverPublicEphemeral"`
	ServerSecretEphemeral string `json:"serverSecretEphemeral"`
}

type ServerSessionResponse struct {
	ServerSessionProof string `json:"serverSessionProof"`
}
