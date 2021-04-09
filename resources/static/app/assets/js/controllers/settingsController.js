ready(() => {
    document.addEventListener('astilectron-ready', function() {

        /// Buttons
        document.querySelector("#resetGlobalDefaultSettingsBtn").addEventListener("click", function() {
            astilectron.sendMessage({ "msg": "resetGlobalDefaultSettings" }, function(message) {
                if (message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error('Unknown error during global settings reset');
                }
            });
        });

        document.querySelector("#clearBotDataBtn").addEventListener("click", function() {
            astilectron.sendMessage({ "msg": "clearAllData" }, function(message) {
                if (message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg);
                }
            });
        });

        document.querySelector("#reinstallDependenciesBtn").addEventListener("click", function() {
            astilectron.sendMessage({ "msg": "reinstallDependencies" }, function(message) {
                if (message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error('Unknown error during dependencies reinstallation');
                }
            });
        });

        /// Forms
        // Settings view
        document.querySelector('#igCredentialsForm').addEventListener("click", function() {
            let message = { "msg": "igCredentialsForm" };
            let content = document.querySelector('#igCredentialsForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if (message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });

        document.querySelector('#quotasForm').addEventListener("click", function() {
            let message = { "msg": "quotasForm" };
            let content = document.querySelector('#quotasForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if (message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });

        document.querySelector('#schedulerForm').addEventListener("click", function(e) {
            let message = { "msg": "schedulerForm" };
            let content = document.querySelector('#schedulerForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if (message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });

        document.querySelector('#blacklistForm').addEventListener("click", function() {
            let message = { "msg": "blacklistForm" };
            let content = document.querySelector('#blacklistForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if (message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });
    });
});