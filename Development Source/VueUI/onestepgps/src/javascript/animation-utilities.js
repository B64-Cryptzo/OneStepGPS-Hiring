export function lerp(start, end, t) {
    return start + (end - start) * t;
  }
  
  export function animateStep(startTime, duration, updateFunction, callback) {
    const currentTime = performance.now();
    const elapsedTime = currentTime - startTime;
    const t = Math.min(elapsedTime / duration, 1);
  
    updateFunction(t);
  
    if (t < 1) {
      requestAnimationFrame(() => animateStep(startTime, duration, updateFunction, callback));
    } else if (callback) {
      callback();
    }
  }
  
  export function animateProperty(obj, property, startValue, endValue, duration) {
    const startTime = performance.now();
    
    animateStep(startTime, duration, (t) => {
      obj[property] = lerp(startValue, endValue, t);
    });
  }
  
  export function animateLatLng(map, targetLatLng, duration = 300) {
    const startLatLng = map.getCenter(); 
    const startTime = performance.now();
  
    animateStep(startTime, duration, (t) => {
      const interpolatedLat = lerp(startLatLng.lat(), targetLatLng.lat, t);
      const interpolatedLng = lerp(startLatLng.lng(), targetLatLng.lng, t);
      map.setCenter(new google.maps.LatLng(interpolatedLat, interpolatedLng));
    });
  }