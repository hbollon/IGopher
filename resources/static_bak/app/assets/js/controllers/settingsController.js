document.addEventListener('astilectron-ready', function() {
    getIgopherConfig()
});

ready(() => {
    document.addEventListener('astilectron-ready', function() {

        /// Buttons
        document.querySelector("#resetGlobalDefaultSettingsBtn").addEventListener("click", function() {
            astilectron.sendMessage({ "msg": "resetGlobalDefaultSettings" }, function(message) {
                if (message.status === SUCCESS) {
                    iziToast.success({
                        message: message.msg,
                    });
                    getIgopherConfig()
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
                    getIgopherConfig()
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

        document.querySelector('#proxyFormBtn').addEventListener("click", function(e) {
            let message = { "msg": "proxyForm" };
            let form = document.querySelector('#proxyForm');
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
        document.querySelector('#proxyAuthCheck').addEventListener("change", function(e) {
            let checkbox = document.querySelector('#proxyAuthCheck')
            let divAuthInputs = document.querySelector('#proxyAuthInputs');
            if (checkbox.checked) {
                checkbox.value = 'true';
                if (divAuthInputs.classList.contains('d-none'))
                    divAuthInputs.classList.remove('d-none');
            } else {
                checkbox.value = 'false';
                if (!divAuthInputs.classList.contains('d-none'))
                    divAuthInputs.classList.add('d-none');
            }

            let authInputs = document.querySelectorAll('.auth-proxy');
            authInputs.forEach(element => {
                if (checkbox.checked) {
                    if (element.required !== true)
                        element.required = true;
                } else {
                    if (element.required !== false)
                        element.required = false;
                }
            });

        });

    });
});

function fillInputs() {
    document.getElementById("dmHourInput").value = igopherConfig.quotas.dmHour;
    document.getElementById("dmDayInput").value = igopherConfig.quotas.dmDay;
    document.getElementById("beginAtInput").value = igopherConfig.schedule.beginAt;
    document.getElementById("endAtInput").value = igopherConfig.schedule.endAt;
    document.getElementById("ipInput").value = igopherConfig.webdriver.proxy.ip;
    document.getElementById("portInput").value = igopherConfig.webdriver.proxy.port;

    if (igopherConfig.webdriver.proxy.auth == "true") {
        document.getElementById("proxyAuthCheck").click();
        document.getElementById("proxyUsernameInput").value = igopherConfig.webdriver.proxy.username;
        document.getElementById("proxyPasswordInput").value = igopherConfig.webdriver.proxy.password;
    }

    document.getElementById("usernameInput").value = igopherConfig.account.username;
    document.getElementById("passwordInput").value = igopherConfig.account.password;
}