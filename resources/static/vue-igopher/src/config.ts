export const iziToast: any = require("izitoast"); // eslint-disable-line
export const bootstrap: any = require("@/bootstrap/js/bootstrap.min.js"); // eslint-disable-line
export const SUCCESS = "Success";
export const ERROR = "Error";
export var igopherConfig: any; // eslint-disable-line

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
}

export function getIgopherConfig(astor: any, callback: () => void) {
    // Get actual IGopher configuration to fill inputs
    astor.trigger("getConfig", {}, (message: any) => {
        if (message.status === SUCCESS) {
            igopherConfig = JSON.parse(message.msg);
            console.log("getIgopherConfig: "+igopherConfig);
            callback();
        } else {
            iziToast.error({
                message: message.msg,
            });
        }
    });
}