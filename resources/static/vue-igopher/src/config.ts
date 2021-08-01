import { Astor } from "./plugins/astilectron";
import Swal from 'sweetalert2'
export const bootstrap: any = require("@/bootstrap/js/bootstrap.min.js"); // eslint-disable-line
export const SUCCESS = "Success";
export const ERROR = "Error";
export var igopherConfig: any; // eslint-disable-line

export const Toast = Swal.mixin({
  toast: true,
  position: 'top-right',
  iconColor: 'white',
  customClass: {
    popup: 'colored-toast'
  },
  showConfirmButton: false,
  showCloseButton: true,
  timer: 3000,
  timerProgressBar: true
});

// Parse JSON Array to JSON Object
export function serialize(data: any) {
  const obj: any = {};
  for (const [key, value] of data) {
    if (obj[key] !== undefined) {
      if (!Array.isArray(obj[key])) {
        obj[key] = [obj[key]];
      }
      obj[key].push(value);
    } else {
      obj[key] = value;
    }
  }
  return obj;
}

// Wait for the DOM to be fully loaded
export const ready = (callback: any) => {
  if (document.readyState != "loading") callback();
  else document.addEventListener("DOMContentLoaded", callback);
};

export function getIgopherConfig(astor: Astor, callback?: () => void): void {
  // Get actual IGopher configuration to fill inputs
  astor.trigger("getConfig", {}, function(message: any) {
    if (message.status === SUCCESS) {
      igopherConfig = JSON.parse(message.msg);
      if (callback !== undefined)
        callback();
    } else {
      Toast.fire({
        icon: 'error',
        title: 'Error',
      });
    }
  });
}
