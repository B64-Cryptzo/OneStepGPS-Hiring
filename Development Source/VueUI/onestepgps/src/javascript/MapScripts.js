import { ref } from "vue";
import { fetchDevices, fetchPreferences, updateUserPreferences } from "@/javascript/api";
import { initializeMap, addMarker, updateMarker, removeMarker, centerMap, clearAllMarkers } from "@/javascript/map-utilities";
import { prepareDeviceData, createFileInput, isPngImage, encodeImageFile } from "@/javascript/device-utilities";
import Gauge from "svg-gauge";

export let map;
export const devices = ref([]);
export const selectedDevice = ref(null);
export const speedGauge = ref(null);

let refreshInterval;

export function handleUploadDeviceImageClick() {
  const fileInput = createFileInput();
  fileInput.onchange = handleFileSelection;
  fileInput.click();
}

async function handleFileSelection(event) {
  const imageFileSelected = event.target.files[0];

  if (imageFileSelected && isPngImage(imageFileSelected)) {
    const imageUrlB64Encoded = await encodeImageFile(imageFileSelected);

    const updatedDeviceUserPreferences = createUpdatedPreferences(imageUrlB64Encoded);
    await updateUserPreferences(updatedDeviceUserPreferences);

    updateMarkerIcon(imageUrlB64Encoded);
  } else {
    alert("Please select a PNG image.");
  }
}

function createUpdatedPreferences(imageUrlB64Encoded) {
  return {
    [selectedDevice.value.device_id]: {
      isVisible: selectedDevice.value.preferences.isVisible,
      markerUrl: imageUrlB64Encoded,
    },
  };
}

function updateMarkerIcon(imageUrlB64Encoded) {
  if (selectedDevice.value && selectedDevice.value.marker) {
    selectedDevice.value.marker.setIcon(imageUrlB64Encoded);
  }
}

function initializeSpeedGaugeWidget() {
  const gaugeContainer = document.getElementById("deviceSpeed");
  if (!gaugeContainer) {
    return null;
  }

  return Gauge(gaugeContainer, {
    max: 150,
    dialStartAngle: 180,
    dialEndAngle: 0,
    color: (value) => (value < 30 ? "green" : value < 80 ? "yellow" : "red"),
    label: function (value) {
      return Math.round(value) + "mph";
    },
    value: 0,
  });
}

function initializeDeviceSpeedGauge() {
  const intervalId = setInterval(() => {
    const deviceSpeedElement = document.getElementById("deviceSpeed");

    if (deviceSpeedElement) {
      if (!speedGauge.value) {
        speedGauge.value = initializeSpeedGaugeWidget();
      }

      speedGauge.value.setValueAnimated(selectedDevice?.value.latest_device_point?.speed, 1);

      clearInterval(intervalId);
    }
  }, 100);
}

async function prepareDevices(userDeviceList, preferences) {
  return userDeviceList.map((device) => ({
    ...prepareDeviceData(device),
    preferences: preferences[device.device_id] || {
      isVisible: true,
      markerUrl: "https://swiftrix.net/uploads/d296f72785e72e6e278d910fbcaf9176.png",
    },
  }));
}

async function updateDeviceMarkers(deviceIDToDeviceMap) {
  devices.value.forEach((existingDevice) => {
    if (!deviceIDToDeviceMap.has(existingDevice.device_id)) {
      removeMarker(existingDevice.device_id);
    }
  });

  deviceIDToDeviceMap.forEach((newDevice) => {
    const existingDevice = devices.value.find(d => d.device_id === newDevice.device_id);
    
    if (existingDevice) {
      const hasPreferencesChanged = JSON.stringify(existingDevice.preferences) !== JSON.stringify(newDevice.preferences);
      const hasDeviceChanged = JSON.stringify(existingDevice) !== JSON.stringify(newDevice);

      if (hasPreferencesChanged || hasDeviceChanged) {
        updateMarker(newDevice);
      }
    } else {
      devices.value.push(newDevice);
      addMarker(newDevice, map, selectDevice);
    }
  });
}

const loadUserDevicesAndPreferences = async () => {
  try {
    const userDeviceList = await fetchDevices();
    const userPreferences = await fetchPreferences();

    const preparedDevices = await prepareDevices(userDeviceList, userPreferences);
    const deviceIDToDeviceMap = new Map(preparedDevices.map((device) => [device.device_id, device]));

    await updateDeviceMarkers(deviceIDToDeviceMap);

    if (!selectedDevice.value && devices.value.length > 0) {
      selectedDevice.value = devices.value[0];
      initializeDeviceSpeedGauge();
    }
  } catch (error) {
    console.error("Failed to load devices or preferences:", error);
  }
};

export const selectDevice = (device) => {
  selectedDevice.value = device;
  initializeDeviceSpeedGauge();
  centerMarkerOnMap(device, map);
};

export function centerMarkerOnMap(deviceWithMarker, mapToCenterOn) {
  centerMap(deviceWithMarker, mapToCenterOn);
}

export const toggleDeviceVisibility = async (device) => {
  const originalVisibility = device.preferences.isVisible;
  device.preferences.isVisible = !originalVisibility;

  try {
    if (device.preferences.isVisible) {
      addMarker(device, map, selectDevice);
    } else {
      removeMarker(device.device_id);
    }

    const updatedPreferences = {
      [device.device_id]: {
        isVisible: device.preferences.isVisible,
        markerUrl: device.preferences.markerUrl,
      },
    };

    await updateUserPreferences(updatedPreferences);
  } catch (error) {
    console.error(`Failed to toggle visibility or update preferences for device ${device.device_id}:`, error);
    device.preferences.isVisible = originalVisibility;
  }
};

function autoClickGoogleBilingCloseButton(selector = ".dismissButton") {
  const PageObserver = new MutationObserver(() => {
    const googleBillingCloseButton = document.querySelector(selector);
    if (googleBillingCloseButton) {
      googleBillingCloseButton.click();
      PageObserver.disconnect();
    }
  });

  PageObserver.observe(document.body, { childList: true, subtree: true });
}

export async function initializeWebComponent() {
  map = await initializeMap("map");

  clearAllMarkers();
  
  autoClickGoogleBilingCloseButton();

  loadUserDevicesAndPreferences();

  refreshInterval = setInterval(loadUserDevicesAndPreferences, 1000);

  speedGauge.value = initializeSpeedGaugeWidget();
}

export function deinitializeWebComponent() {
  clearInterval(refreshInterval);
}
