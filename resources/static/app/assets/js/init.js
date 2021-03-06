// Dynamics buttons inits
var dmBotRunning = sessionStorage.getItem("botState");
console.log(dmBotRunning);
if(dmBotRunning === "false" || dmBotRunning === null) {
    $('#dmBotLaunchBtn').addClass('btn-success').removeClass('btn-danger');
    $('#dmBotLaunchIcon').addClass('fa-rocket').removeClass('fa-skull-crossbones');
    $('#dmBotLaunchSpan').text('Launch !');
} else {
    $('#dmBotLaunchBtn').addClass('btn-danger').removeClass('btn-success');
    $('#dmBotLaunchIcon').addClass('fa-skull-crossbones').removeClass('fa-rocket');
    $('#dmBotLaunchSpan').text('Stop !');
}

// Toastr configuration for notification system
toastr.options = {
  "closeButton": true,
  "debug": false,
  "newestOnTop": true,
  "progressBar": false,
  "positionClass": "toast-top-right",
  "preventDuplicates": false,
  "onclick": null,
  "showDuration": "300",
  "hideDuration": "1000",
  "timeOut": "5000",
  "extendedTimeOut": "1000",
  "showEasing": "swing",
  "hideEasing": "linear",
  "showMethod": "fadeIn",
  "hideMethod": "fadeOut"
}

// Wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // Listen to messages sent by Go
    astilectron.onMessage(function(message) {
        // Process message
        if (message.msg === "bot crash") {
            if(dmBotRunning) {
                dmBotRunning = false;
                $('#dmBotLaunchBtn').addClass('btn-success').removeClass('btn-danger');
                $('#dmBotLaunchIcon').addClass('fa-rocket').removeClass('fa-skull-crossbones');
                $('#dmBotLaunchSpan').text('Launch !');
                toastr.error(message.payload)
            }
        }
    });
})