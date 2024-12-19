package models

/*
Universal struct used in SRP authentication which
allows us to use a single API endpoint to accomplish
3 distinct steps during the SRP protocol
*/

type AuthenticationRequest struct {
	Identifier            int    `json:"identifier"`
	Username              string `json:"username,omitempty"`
	ClientPublicEphemeral string `json:"clientPublicEphemeral,omitempty"`
	ClientSessionProof    string `json:"clientSessionProof,omitempty"`
	ClientSessionKey      string `json:"clientSessionKey,omitempty"`
}
