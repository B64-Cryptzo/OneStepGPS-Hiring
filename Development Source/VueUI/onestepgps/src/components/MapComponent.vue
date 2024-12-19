<template>
  <div id="app">
    <div class="main-page">
      <div id="map" class="google-map-container">
        
      </div>
      <div class="device-list-container">
        <div class="device-list-header">
          <span class="header-device-name">Device Name</span>
          <span class="header-device-id">Device ID</span>
          <span class="header-device-state">State</span>
          <span class="header-device-actions">Location</span>
          <span class="header-device-visibility">Visibility</span>
        </div>
        <div class="device-listbox">
          <div
            v-for="device in devices"
            :key="device.device_id"
            class="listbox-item"
            @click="selectDevice(device)"
          >
            <span class="device-name">{{ device.display_name }}</span>
            <span class="device-id">{{ device.device_id }}</span>
            <span class="device-state">{{ device.drive_status }}</span>
            <span class="magnifying-glass" @click="centerMarkerOnMap(device, map); $event.stopPropagation()">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
                <path
                  d="M416 208c0 45.9-14.9 88.3-40 122.7L502.6 457.4c12.5 12.5 12.5 32.8 0 45.3s-32.8 12.5-45.3 0L330.7 376c-34.4 25.2-76.8 40-122.7 40C93.1 416 0 322.9 0 208S93.1 0 208 0S416 93.1 416 208zM208 352a144 144 0 1 0 0-288 144 144 0 1 0 0 288z"
                />
              </svg>
            </span>
            <span class="eye-icon" @click="toggleDeviceVisibility(device); $event.stopPropagation()">
              <svg
                v-if="device.preferences.isVisible"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 640 512"
                class="hidden"
              >
              <path d="M288 32c-80.8 0-145.5 36.8-192.6 80.6C48.6 156 17.3 208 2.5 243.7c-3.3 7.9-3.3 16.7 0 24.6C17.3 304 48.6 356 95.4 399.4C142.5 443.2 207.2 480 288 480s145.5-36.8 192.6-80.6c46.8-43.5 78.1-95.4 93-131.1c3.3-7.9 3.3-16.7 0-24.6c-14.9-35.7-46.2-87.7-93-131.1C433.5 68.8 368.8 32 288 32zM144 256a144 144 0 1 1 288 0 144 144 0 1 1 -288 0zm144-64c0 35.3-28.7 64-64 64c-7.1 0-13.9-1.2-20.3-3.3c-5.5-1.8-11.9 1.6-11.7 7.4c.3 6.9 1.3 13.8 3.2 20.7c13.7 51.2 66.4 81.6 117.6 67.9s81.6-66.4 67.9-117.6c-11.1-41.5-47.8-69.4-88.6-71.1c-5.8-.2-9.2 6.1-7.4 11.7c2.1 6.4 3.3 13.2 3.3 20.3z"/>
              </svg>
              <svg
                v-else
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 640 512"
                class="visible"
              >
                <path d="M38.8 5.1C28.4-3.1 13.3-1.2 5.1 9.2S-1.2 34.7 9.2 42.9l592 464c10.4 8.2 25.5 6.3 33.7-4.1s6.3-25.5-4.1-33.7L525.6 386.7c39.6-40.6 66.4-86.1 79.9-118.4c3.3-7.9 3.3-16.7 0-24.6c-14.9-35.7-46.2-87.7-93-131.1C465.5 68.8 400.8 32 320 32c-68.2 0-125 26.3-169.3 60.8L38.8 5.1zM223.1 149.5C248.6 126.2 282.7 112 320 112c79.5 0 144 64.5 144 144c0 24.9-6.3 48.3-17.4 68.7L408 294.5c8.4-19.3 10.6-41.4 4.8-63.3c-11.1-41.5-47.8-69.4-88.6-71.1c-5.8-.2-9.2 6.1-7.4 11.7c2.1 6.4 3.3 13.2 3.3 20.3c0 10.2-2.4 19.8-6.6 28.3l-90.3-70.8zM373 389.9c-16.4 6.5-34.3 10.1-53 10.1c-79.5 0-144-64.5-144-144c0-6.9 .5-13.6 1.4-20.2L83.1 161.5C60.3 191.2 44 220.8 34.5 243.7c-3.3 7.9-3.3 16.7 0 24.6c14.9 35.7 46.2 87.7 93 131.1C174.5 443.2 239.2 480 320 480c47.8 0 89.9-12.9 126.2-32.5L373 389.9z"/>
              </svg>
            </span>
          </div>
        </div>
      </div>
    </div>

    <div class="selected-device-container">
      <div v-if="selectedDevice">
        <h2 class="selected-device-name">
          {{ selectedDevice.display_name || "No device selected" }}
        </h2>
        <p class="selected-device-vin">
          VIN: {{ selectedDevice.latest_device_point.params.vin || "Not Available" }}
        </p>
        <span class="selected-device-image" @click="handleUploadDeviceImageClick()">
          <img v-if="selectedDevice && selectedDevice.preferences.markerUrl" :src="selectedDevice.preferences.markerUrl" alt="Selected Device" />
          <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
            <path d="M448 80c8.8 0 16 7.2 16 16l0 319.8-5-6.5-136-176c-4.5-5.9-11.6-9.3-19-9.3s-14.4 3.4-19 9.3L202 340.7l-30.5-42.7C167 291.7 159.8 288 152 288s-15 3.7-19.5 10.1l-80 112L48 416.3l0-.3L48 96c0-8.8 7.2-16 16-16l384 0zM64 32C28.7 32 0 60.7 0 96L0 416c0 35.3 28.7 64 64 64l384 0c35.3 0 64-28.7 64-64l0-320c0-35.3-28.7-64-64-64L64 32zm80 192a48 48 0 1 0 0-96 48 48 0 1 0 0 96z"/>
          </svg>
          <label>click to change device icon</label>
        </span>
        <div id="deviceSpeed" class="gauge-container">
        </div>
        <div class="gas-icon">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
            <path d="M32 64C32 28.7 60.7 0 96 0L256 0c35.3 0 64 28.7 64 64l0 192 8 0c48.6 0 88 39.4 88 88l0 32c0 13.3 10.7 24 24 24s24-10.7 24-24l0-154c-27.6-7.1-48-32.2-48-62l0-64L384 64c-8.8-8.8-8.8-23.2 0-32s23.2-8.8 32 0l77.3 77.3c12 12 18.7 28.3 18.7 45.3l0 13.5 0 24 0 32 0 152c0 39.8-32.2 72-72 72s-72-32.2-72-72l0-32c0-22.1-17.9-40-40-40l-8 0 0 144c17.7 0 32 14.3 32 32s-14.3 32-32 32L32 512c-17.7 0-32-14.3-32-32s14.3-32 32-32L32 64zM96 80l0 96c0 8.8 7.2 16 16 16l128 0c8.8 0 16-7.2 16-16l0-96c0-8.8-7.2-16-16-16L112 64c-8.8 0-16 7.2-16 16z"/>
          </svg>
          <span> 
            {{
              selectedDevice.latest_device_point.device_state.fuel_percent !== undefined
                ? Math.round(selectedDevice.latest_device_point.device_state.fuel_percent) + "%"
                : "N/A" 
            }}
          </span>
        </div>
      </div>
      <div v-else>
        <p>No device selected</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount } from 'vue';
import { devices,
  map,
  selectedDevice,
  handleUploadDeviceImageClick, 
  selectDevice, 
  toggleDeviceVisibility, 
  centerMarkerOnMap,
  initializeWebComponent, 
  deinitializeWebComponent } from "@/javascript/MapScripts.js";

onMounted(async() => {
  initializeWebComponent();
});

onBeforeUnmount(async() =>{
  deinitializeWebComponent();
});

</script>

<style scoped src="@/style/MapPage.css"></style>