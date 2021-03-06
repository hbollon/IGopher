const SUCCESS = "Success"
const ERROR = "Error"

// Parse JSON Array to JSON Object
$.fn.serializeObject = function(){
    var obj = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (obj[this.name] !== undefined) {
            if (!obj[this.name].push) {
                obj[this.name] = [obj[this.name]];
            }
            obj[this.name].push(this.value || '');
        } else {
            obj[this.name] = this.value || '';
        }
    });
    return obj;
};

// Messages sending routine to Go backend
$(document).ready(function(){ 
    document.addEventListener('astilectron-ready', function() {
        
        /// Buttons
        $("#resetGlobalDefaultSettingsBtn").on("click", function(){
            astilectron.sendMessage({"msg":"resetGlobalDefaultSettings"}, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error('Unknown error during global settings reset');
                }
            });
        }); 
        
        $("#clearBotDataBtn").on("click", function(){
            astilectron.sendMessage({"msg":"clearAllData"}, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg);
                }
            });
        }); 
        
        $("#reinstallDependenciesBtn").on("click", function(){
            astilectron.sendMessage({"msg":"reinstallDependencies"}, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error('Unknown error during dependencies reinstallation');
                }
            });
        }); 
        
        $("#dmBotLaunchBtn").on("click", function(){
            if(dmBotRunning === "false" || dmBotRunning === false || dmBotRunning === null) {
                astilectron.sendMessage({"msg":"launchDmBot"}, function(message) {
                    if(message.status === SUCCESS) {
                        toastr.success(message.msg);
                        dmBotRunning = true 
                        $('#dmBotLaunchBtn').addClass('btn-danger').removeClass('btn-success');
                        $('#dmBotLaunchIcon').addClass('fa-skull-crossbones').removeClass('fa-rocket');
                        $('#dmBotLaunchSpan').text('Stop !');
                        sessionStorage.setItem("botState", true)
                    } else {
                        toastr.error(message.msg);
                    }
                });
            } else {
                $('#dmBotLaunchIcon').addClass('fa-spinner').addClass('fa-spin').removeClass('fa-skull-crossbones');
                toastr.info("Stop procedure launched, the bot will stop once the current action is finished.")
                astilectron.sendMessage({"msg":"stopDmBot"}, function(message) {
                    if(message.status === SUCCESS) {
                        toastr.success(message.msg);
                        dmBotRunning = false
                        $('#dmBotLaunchBtn').addClass('btn-success').removeClass('btn-danger');
                        $('#dmBotLaunchIcon').addClass('fa-rocket').removeClass('fa-spinner').removeClass('fa-spin');
                        $('#dmBotLaunchSpan').text('Launch !');
                        sessionStorage.setItem("botState", false)
                    } else {
                        $('#dmBotLaunchIcon').addClass('fa-skull-crossbones').removeClass('fa-spinner').removeClass('fa-spin');
                        toastr.error(message.msg);
                    }
                });
            }
        }); 
        
        $("#dmBotHotReloadBtn").on("click", function(){
            $('#dmBotHotReloadIcn').addClass('fa-spinner').addClass('fa-spin').removeClass('fa-fire');
            toastr.info("Hot reload launched, the bot will update once the current action is finished.")
            astilectron.sendMessage({"msg":"hotReloadBot"}, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                    $('#dmBotHotReloadIcn').addClass('fa-fire').removeClass('fa-spinner').removeClass('fa-spin');
                } else {
                    toastr.error(message.msg);
                    $('#dmBotHotReloadIcn').addClass('fa-fire').removeClass('fa-spinner').removeClass('fa-spin');
                }
            });
        }); 
        
        /// Forms
        // Settings view
        $('#igCredentialsForm').submit(function() {
            let message = {"msg": "igCredentialsForm"};
            let content = $('#igCredentialsForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });
        
        $('#quotasForm').submit(function() {
            let message = {"msg": "quotasForm"};
            let content = $('#quotasForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });
        
        $('#schedulerForm').submit(function() {
            let message = {"msg": "schedulerForm"};
            let content = $('#schedulerForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });
        
        $('#blacklistForm').submit(function() {
            let message = {"msg": "blacklistForm"};
            let content = $('#blacklistForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });
        
        // Dm automation view
        $('#dmSettingsForm').submit(function() {
            let message = {"msg": "dmSettingsForm"};
            let content = $('#dmSettingsForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });
        
        $('#dmUserScrappingSettingsForm').submit(function() {
            let message = {"msg": "dmUserScrappingSettingsForm"};
            let content = $('#dmUserScrappingSettingsForm').serializeObject();
            if (typeof content !== "undefined") {
                message.payload = content;
            }
            astilectron.sendMessage(message, function(message) {
                if(message.status === SUCCESS) {
                    toastr.success(message.msg);
                } else {
                    toastr.error(message.msg, "Error during settings saving!");
                }
            });
            return false; // avoid page reload
        });
    });
});