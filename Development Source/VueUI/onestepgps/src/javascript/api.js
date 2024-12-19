import axios from "axios";

const GoBackendApiUrl = "http://localhost:8080/api";

export async function fetchDevices() {
  try {
    const sessionToken = localStorage.getItem("sessionToken");
    if (!sessionToken) throw new Error("Session token not found. Please log in.");

    const response = await axios.get(`${GoBackendApiUrl}/devices`, {
      headers: {
        Authorization: `${sessionToken}`,
      },
    });

    return response.data.result_list;
  } catch (error) {
    handleSessionError(error);
    throw error; 
  }
}

export async function fetchPreferences() {
  try {
    const sessionToken = localStorage.getItem("sessionToken");
    if (!sessionToken) throw new Error("Session token not found. Please log in.");

    const preferencesResponse = await fetch(`${GoBackendApiUrl}/preferences`, {
      headers: {
        Authorization: `${sessionToken}`,
      },
    });

    if (!preferencesResponse.ok) {
      const errorData = await preferencesResponse.json();
      if (errorData.message === "Invalid session token") {
        handleSessionError(new Error("Invalid session token"));
      }
      throw new Error(errorData.message || "Failed to fetch preferences.");
    }

    return await preferencesResponse.json();
  } catch (error) {
    handleSessionError(error); 
    throw error;
  }
}

export async function updateUserPreferences(preferences) {
  try {
    const sessionToken = localStorage.getItem("sessionToken");

    if (!sessionToken) throw new Error("Session token not found. Please log in.");

    const response = await fetch(`${GoBackendApiUrl}/update-user-preferences`, {
      method: 'POST',
      headers: {
        'Authorization': sessionToken,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(preferences),
    });

    if (!response.ok) {
      const error = await response.json();
      
      throw new Error('Failed to update preferences: ' + error);
    }

    const result = await response.json();

    return result;

  } catch (error) {
    console.error('Error while updating preferences:', error);
    throw error;
  }
}

function handleSessionError(error) {
  console.error("Session token is invalid: ", error);
  localStorage.removeItem("sessionToken"); 
  window.location.reload();
}