export function prepareDeviceData(device) {
    device.latest_device_point = {
      lat: device.latest_device_point?.lat || null,
      lng: device.latest_device_point?.lng || null,
      speed: device.latest_device_point?.speed || null,
      device_state: device.latest_device_point?.device_state || null,
      params: device.latest_device_point?.params || null,
    };

    device.speed = device.latest_device_point.speed || 0;
    device.drive_status = device.latest_device_point.device_state?.drive_status || "unknown";
    device.vin = device.latest_device_point.params?.vin || "N/A";
    device.fuel_percent = device.latest_device_point.device_state?.fuel_percent ?? null;
    return device;
  }
  
  export function createFileInput() {
    const fileInput = document.createElement("input");
    fileInput.type = "file";
    fileInput.accept = "image/png";
    return fileInput;
  }

  export function isPngImage(file) {
    return file.type === "image/png";
  }
  
  export function encodeImageFile(file) {
    const reader = new FileReader();
  
    return new Promise((resolve, reject) => {
      reader.onload = () => {
        const base64Data = reader.result;
  
        resolve(base64Data);
      };
  
      reader.onerror = reject;
  
      reader.readAsDataURL(file);
    });
  }