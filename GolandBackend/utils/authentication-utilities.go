package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golandbackend/models"
	"log"
	"net/http"
	"os/exec"
)

func GenerateSessionToken() string {
	tokenBytes := make([]byte, 16)
	_, errorLog := rand.Read(tokenBytes)
	if errorLog != nil {
		fmt.Println("Error generating session token:", errorLog)
		return ""
	}

	return hex.EncodeToString(tokenBytes)
}

func GenerateServerEphemeral(verifier, salt string) (*models.ServerEphemeralResponse, error) {
	subServerUrl := "http://localhost:8082/generate-server-ephemeral"

	requestBody := map[string]interface{}{
		"verifier": verifier,
		"salt":     salt,
	}

	return MakeSecuredPostRequest[models.ServerEphemeralResponse](subServerUrl, requestBody)
}

func DeriveServerSession(serverSecretEphemeral, clientPublicEphemeral, salt, username, verifier, clientSessionProof string) (*models.ServerSessionResponse, error) {
	subServerUrl := "http://localhost:8082/derive-server-session"

	requestBody := map[string]interface{}{
		"serverSecretEphemeral": serverSecretEphemeral,
		"clientPublicEphemeral": clientPublicEphemeral,
		"salt":                  salt,
		"username":              username,
		"verifier":              verifier,
		"clientSessionProof":    clientSessionProof,
	}

	return MakeSecuredPostRequest[models.ServerSessionResponse](subServerUrl, requestBody)
}

func SRPGenerateServerEphemeral(responseWriter http.ResponseWriter, authenticationRequest models.AuthenticationRequest, users []map[string]interface{}) {
	user, errorLog := FindUserByUsername(authenticationRequest.Username, users)
	if HandleError(responseWriter, "User not found", errorLog, http.StatusInternalServerError) {
		return
	}

	userVerifier := user["verifier"].(string)
	userSalt := user["salt"].(string)

	serverEphemeral, errorLog := GenerateServerEphemeral(userVerifier, userSalt)
	if HandleError(responseWriter, "Unable to generate Ephemeral", errorLog, http.StatusInternalServerError) {
		return
	}

	serverEphemeralData := map[string]interface{}{
		"serverPublicEphemeral": serverEphemeral.ServerPublicEphemeral,
		"serverSecretEphemeral": serverEphemeral.ServerSecretEphemeral,
	}

	user["serverEphemeral"] = serverEphemeralData

	if errorLog := SaveUserData(users); errorLog != nil {
		HandleError(responseWriter, "Error saving user data", errorLog, http.StatusInternalServerError)
		return
	}

	serverAuthenticationResponse := map[string]interface{}{
		"identifier":            2,
		"salt":                  userSalt,
		"serverPublicEphemeral": serverEphemeral.ServerPublicEphemeral,
	}

	json.NewEncoder(responseWriter).Encode(serverAuthenticationResponse)
}

func SRPCreateServerSession(responseWriter http.ResponseWriter, authenticationRequest models.AuthenticationRequest, users []map[string]interface{}) {
	user, errorLog := FindUserByUsername(authenticationRequest.Username, users)
	if HandleError(responseWriter, "User not found", errorLog, http.StatusInternalServerError) {
		return
	}

	serverEphemeral, ok := user["serverEphemeral"].(map[string]interface{})
	if !ok || len(serverEphemeral) == 0 {
		HandleError(responseWriter, "Server ephemeral not initialized", nil, http.StatusInternalServerError)
		return
	}

	userSalt, ok := user["salt"].(string)
	if !ok {
		HandleError(responseWriter, "Invalid or missing 'salt'", nil, http.StatusInternalServerError)
		return
	}

	userIdentifier, ok := user["username"].(string)
	if !ok {
		HandleError(responseWriter, "Invalid or missing 'username'", nil, http.StatusInternalServerError)
		return
	}

	userVerifier, ok := user["verifier"].(string)
	if !ok {
		HandleError(responseWriter, "Invalid or missing 'verifier'", nil, http.StatusInternalServerError)
		return
	}

	clientPublicEphemeral := authenticationRequest.ClientPublicEphemeral
	if clientPublicEphemeral == "" {
		HandleError(responseWriter, "Invalid or missing 'clientPublicEphemeral'", nil, http.StatusInternalServerError)
		return
	}

	clientSessionProof := authenticationRequest.ClientSessionProof
	if clientSessionProof == "" {
		HandleError(responseWriter, "Invalid or missing 'clientSessionProof'", nil, http.StatusInternalServerError)
		return
	}

	serverSessionResponse, errorLog := DeriveServerSession(
		serverEphemeral["serverSecretEphemeral"].(string),
		clientPublicEphemeral,
		userSalt,
		userIdentifier,
		userVerifier,
		clientSessionProof)
	if HandleError(responseWriter, "Unable to generate Ephemeral", errorLog, http.StatusInternalServerError) {
		return
	}

	user["sessionToken"] = GenerateSessionToken()

	if errorLog := SaveUserData(users); errorLog != nil {
		HandleError(responseWriter, "Error saving user data", errorLog, http.StatusInternalServerError)
		return
	}

	serverAuthenticationResponse := map[string]interface{}{
		"identifier":         4,
		"sessionToken":       user["sessionToken"],
		"serverSessionProof": serverSessionResponse.ServerSessionProof,
	}
	json.NewEncoder(responseWriter).Encode(serverAuthenticationResponse)
}

func StartSubServerSRP() *exec.Cmd {
	cmd := exec.Command("node", "./srpJS/srp.js")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	if errorLog := cmd.Start(); errorLog != nil {
		log.Fatalf("Failed to start srp.js: %v", errorLog)
	}

	return cmd
}

func TerminateProcess(cmd *exec.Cmd) {
	if errorLog := cmd.Process.Kill(); errorLog != nil {
		log.Printf("Failed to kill process: %v", errorLog)
	} else {
		log.Println("Process terminated.")
	}
}
