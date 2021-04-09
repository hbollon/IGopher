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
    let dmBotLaunchBtn = document.querySelector('#dmBotLaunchBtn')
    let dmBotLaunchIcon = document.querySelector('#dmBotLaunchIcon')
    let dmBotLaunchSpan = document.querySelector('#dmBotLaunchSpan')
        // Listen to messages sent by Go
    astilectron.onMessage(function(message) {
        // Process message
        if (message.msg === "bot crash") {
            if (dmBotRunning) {
                dmBotRunning = false;
                dmBotLaunchBtn.classList.add('btn-success');
                dmBotLaunchBtn.classList.remove('btn-danger');
                dmBotLaunchIcon.classList.add('fa-rocket');
                dmBotLaunchIcon.classList.remove('fa-spinner', 'fa-spin');
                dmBotLaunchSpan.innerHTML = 'Launch !';
                toastr.error(message.payload)
            }
        }
    });
})