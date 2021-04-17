ready(() => {
    document.addEventListener('astilectron-ready', function() {

        /// Buttons
        document.querySelector("#resetGlobalDefaultSettingsBtn").addEventListener("click", function() {
            astilectron.sendMessage({ "msg": "resetGlobalDefaultSettings" }, function(message) {
                if (message.status === SUCCESS) {
                    iziToast.success({
                        message: message.msg,
                    });
                } else {
                    iziToast.error({
                        message: "Unknown error during global settings reset",
                    });
                }
            });
        });

        document.querySelector("#clearBotDataBtn").addEventListener("click", function() {
            astilectron.sendMessage({ "msg": "clearAllData" }, function(message) {
                if (message.status === SUCCESS) {
                    iziToast.success({
                        message: message.msg,
                    });
                } else {
                    iziToast.error({
                        message: message.msg,
                    });
                }
            });
        });

        // document.querySelector("#reinstallDependenciesBtn").addEventListener("click", function() {
        //     astilectron.sendMessage({ "msg": "reinstallDependencies" }, function(message) {
        //         if (message.status === SUCCESS) {
        //             iziToast.success({
        //     message: message.msg,
        // });
        //         } else {
        //             toastr.error('Unknown error during dependencies reinstallation');
        //         }
        //     });
        // });

        /// Forms
        // Settings view
        document.querySelector('#igCredentialsFormBtn').addEventListener("click", function(e) {
            let message = { "msg": "igCredentialsForm" };
            let form = document.querySelector('#igCredentialsForm');
            if (!form.checkValidity()) {
                e.preventDefault()
                e.stopPropagation()
            } else {
                if (typeof content !== "undefined") {
                    let formData = new FormData(form);
                    message.payload = serialize(formData);
                }
                astilectron.sendMessage(message, function(message) {
                    if (message.status === SUCCESS) {
                        iziToast.success({
                            message: message.msg,
                        });
                    } else {
                        iziToast.error({
                            title: "Error during settings saving!",
                            message: message.msg,
                        });
                    }
                });
            }

            form.classList.add('was-validated')
        });

        document.querySelector('#quotasFormBtn').addEventListener("click", function(e) {
            let message = { "msg": "quotasForm" };
            let form = document.querySelector('#quotasForm');
            if (!form.checkValidity()) {
                e.preventDefault()
                e.stopPropagation()
            } else {
                if (typeof content !== "undefined") {
                    let formData = new FormData(form);
                    message.payload = serialize(formData);
                }
                astilectron.sendMessage(message, function(message) {
                    if (message.status === SUCCESS) {
                        iziToast.success({
                            message: message.msg,
                        });
                    } else {
                        iziToast.error({    
                            title: "Error during settings saving!",
                            message: message.msg,
                        });
                    }
                });
            }

            form.classList.add('was-validated')
        });

        document.querySelector('#schedulerFormBtn').addEventListener("click", function(e) {
            let message = { "msg": "schedulerForm" };
            let form = document.querySelector('#schedulerForm');
            if (!form.checkValidity()) {
                e.preventDefault()
                e.stopPropagation()
            } else {
                if (typeof content !== "undefined") {
                    let formData = new FormData(form);
                    message.payload = serialize(formData);
                }
                astilectron.sendMessage(message, function(message) {
                    if (message.status === SUCCESS) {
                        iziToast.success({
                            message: message.msg,
                        });
                    } else {
                        iziToast.error({
                            title: "Error during settings saving!",
                            message: message.msg,
                        });
                    }
                });
            }

            form.classList.add('was-validated')
        });

        document.querySelector('#blacklistFormBtn').addEventListener("click", function(e) {
            let message = { "msg": "blacklistForm" };
            let form = document.querySelector('#blacklistForm');
            if (!form.checkValidity()) {
                e.preventDefault()
                e.stopPropagation()
            } else {
                if (typeof content !== "undefined") {
                    let formData = new FormData(form);
                    message.payload = serialize(formData);
                }
                astilectron.sendMessage(message, function(message) {
                    if (message.status === SUCCESS) {
                        iziToast.success({
                            message: message.msg,
                        });
                    } else {
                        iziToast.error({
                            title: "Error during settings saving!",
                            message: message.msg,
                        });
                    }
                });
            }

            form.classList.add('was-validated')
        });
    });
});