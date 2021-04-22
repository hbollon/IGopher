var igopherConfig;

ready(() => {
    document.addEventListener('astilectron-ready', function() {

        // Get actual IGopher configuration to fill inputs
        astilectron.sendMessage({ "msg": "getConfig" }, function(message) {
            if (message.status === SUCCESS) {
                igopherConfig = JSON.parse(message.msg);
                console.log(igopherConfig);
                fillInputs();
            } else {
                iziToast.error({
                    message: message.msg,
                });
            }
        });

        let dmBotLaunchBtn = document.querySelector('#dmBotLaunchBtn')
        let dmBotLaunchIcon = document.querySelector('#dmBotLaunchIcon')
        let dmBotLaunchSpan = document.querySelector('#dmBotLaunchSpan')
        let dmBotHotReloadBtn = document.querySelector('#dmBotHotReloadBtn')
        let dmBotHotReloadIcn = document.querySelector('#dmBotHotReloadIcn')

        // Dynamics buttons inits
        var dmBotRunning = sessionStorage.getItem("botState");
        if (dmBotRunning === "false" || dmBotRunning === null) {
            dmBotLaunchBtn.classList.add('btn-success');
            dmBotLaunchBtn.classList.remove('btn-danger');
            dmBotLaunchIcon.classList.add('fa-rocket');
            dmBotLaunchIcon.classList.remove('fa-spinner', 'fa-spin');
            dmBotLaunchSpan.innerHTML = 'Launch !';
        } else {
            dmBotLaunchBtn.classList.add('btn-danger');
            dmBotLaunchBtn.classList.remove('btn-success');
            dmBotLaunchIcon.classList.add('fa-skull-crossbones');
            dmBotLaunchIcon.classList.remove('fa-rocket');
            dmBotLaunchSpan.innerHTML = 'Stop !';
        }

        /// Buttons
        dmBotLaunchBtn.addEventListener("click", function() {
            if (dmBotRunning === "false" || dmBotRunning === false || dmBotRunning === null) {
                astilectron.sendMessage({ "msg": "launchDmBot" }, function(message) {
                    if (message.status === SUCCESS) {
                        iziToast.success({
                            message: message.msg,
                        });

                        dmBotRunning = true
                        dmBotLaunchBtn.classList.add('btn-danger');
                        dmBotLaunchBtn.classList.remove('btn-success');
                        dmBotLaunchIcon.classList.add('fa-skull-crossbones');
                        dmBotLaunchIcon.classList.remove('fa-rocket');
                        dmBotLaunchSpan.innerHTML = 'Stop !';
                        sessionStorage.setItem("botState", true)
                    } else {
                        iziToast.error({
                            message: message.msg,
                        });
                    }
                });
            } else {
                dmBotLaunchIcon.classList.add('fa-spinner', 'fa-spin');
                dmBotLaunchIcon.classList.remove('fa-skull-crossbones');
                iziToast.info({
                    message: "Stop procedure launched, the bot will stop once the current action is finished.",
                });
                astilectron.sendMessage({ "msg": "stopDmBot" }, function(message) {
                    if (message.status === SUCCESS) {
                        iziToast.success({
                            message: message.msg,
                        });
                        dmBotRunning = false
                        dmBotLaunchBtn.classList.add('btn-success');
                        dmBotLaunchBtn.classList.remove('btn-danger');
                        dmBotLaunchIcon.classList.add('fa-rocket');
                        dmBotLaunchIcon.classList.remove('fa-spinner', 'fa-spin');
                        dmBotLaunchSpan.innerHTML = 'Launch !';
                        sessionStorage.setItem("botState", false)
                    } else {
                        dmBotLaunchIcon.classList.add('fa-skull-crossbones');
                        dmBotLaunchIcon.classList.remove('fa-spinner', 'fa-spin');
                        iziToast.error({
                            message: message.msg,
                        });
                    }
                });
            }
        });

        dmBotHotReloadBtn.addEventListener("click", function() {
            dmBotHotReloadIcn.classList.add('fa-spinner', 'fa-spin');
            dmBotHotReloadIcn.classList.remove('fa-fire');
            iziToast.info({
                message: "Hot reload launched, the bot will update once the current action is finished.",
            });
            astilectron.sendMessage({ "msg": "hotReloadBot" }, function(message) {
                if (message.status === SUCCESS) {
                    iziToast.success({
                        message: message.msg,
                    });
                    dmBotHotReloadIcn.classList.add('fa-fire');
                    dmBotHotReloadIcn.classList.remove('fa-spinner', 'fa-spin');
                } else {
                    iziToast.error({
                        message: message.msg,
                    });
                    dmBotHotReloadIcn.classList.add('fa-fire');
                    dmBotHotReloadIcn.classList.remove('fa-spinner', 'fa-spin');
                }
            });
        });

        /// Forms
        // Dm automation view
        document.querySelector('#dmSettingsFormBtn').addEventListener("click", function(e) {
            let message = { "msg": "dmSettingsForm" };
            let form = document.querySelector('#dmSettingsForm');
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

        document.querySelector('#dmUserScrappingSettingsFormBtn').addEventListener("click", function(e) {
            let message = { "msg": "dmUserScrappingSettingsForm" };
            let form = document.querySelector('#dmUserScrappingSettingsForm');
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

function fillInputs() {
    document.getElementById("dmTemplates").value = igopherConfig.auto_dm.dmTemplates.join(";");
    document.getElementById("greetingTemplate").value = igopherConfig.auto_dm.greeting.greetingTemplate;
    document.getElementById("srcUsers").value = igopherConfig.scrapper.srcUsers.join(";");
    document.getElementById("scrappingQuantity").value = igopherConfig.scrapper.scrappingQuantity;
}