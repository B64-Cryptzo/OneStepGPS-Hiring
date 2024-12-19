/* global google */

import { animateLatLng, animateProperty } from './animation-utilities';

let markersMap = new Map(); 

export async function initializeMap(mapElementId) {
  return new google.maps.Map(document.getElementById(mapElementId), {
    zoomControl: true,
    scrollwheel: true,
    zoom: 5,
    center: { lat: 34.0522, lng: -118.2437 }, 
  });
}

export function centerMap(device, map) {
  if (device?.latest_device_point?.lat && device?.latest_device_point?.lng) {
    const targetLatLng = {
      lat: device.latest_device_point.lat,
      lng: device.latest_device_point.lng,
    };

    animateLatLng(map, targetLatLng); 
    map.setZoom(11); 
  }
}

export function addMarker(device, map, selectDeviceCallback) {
  if (markersMap.has(device.device_id)) {
    return;
  }
  
  if (!device?.latest_device_point?.lat ||
    !device?.latest_device_point?.lng) {
      return;
    }

    const mapMarker = new google.maps.Marker({
      position: {
        lat: device.latest_device_point.lat,
        lng: device.latest_device_point.lng,
      },
      map: map,
      icon: {
        url: device.preferences.markerUrl,
        scaledSize: new google.maps.Size(35, 35),
        labelOrigin: new google.maps.Point(20, -5),
      },
      label: {
        text: device.display_name,
        color: "#ffffff",       
        fontSize: "12px",       
        fontWeight: "bold",      
      },
    });

    mapMarker.setVisible(device?.preferences?.isVisible);
    
    markersMap.set(device.device_id, mapMarker);

    mapMarker.addListener("click", () => {
      if (selectDeviceCallback) {
        selectDeviceCallback(device);
      }
    });
}

export function removeMarker(deviceId) {
  const mapMarker = markersMap.get(deviceId);
  if (mapMarker) {
    mapMarker.setMap(null); 
    markersMap.delete(deviceId);
  }
}

export function updateMarker(device) {

  if (!device?.latest_device_point?.lat ||
    !device?.latest_device_point?.lng
    ){
      return;
    }

    const mapMarker = markersMap.get(device.device_id);

    mapMarker.setVisible(device?.preferences?.isVisible);
    
    if (mapMarker) {
      const oldPosition = mapMarker.getPosition();
      const newPosition = new google.maps.LatLng(
        device.latest_device_point.lat,
        device.latest_device_point.lng
      );

      const startTime = performance.now();
      const duration = 500; 

      animateProperty(mapMarker, oldPosition, newPosition, startTime, duration);
      mapMarker.setTitle(device.display_name);

      const newIcon = {
        url: device.preferences.markerUrl,
        scaledSize: new google.maps.Size(35, 35),
        labelOrigin: new google.maps.Point(20, -5),
      };
  
      mapMarker.setIcon(newIcon);
    }
  
}

export function clearAllMarkers() {
  markersMap.forEach((mapMarker) => {
    mapMarker.setMap(null);
  });
  markersMap.clear();
}