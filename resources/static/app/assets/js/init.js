// iziToast configuration for notification system
iziToast.settings({
    position: 'topRight',
    timeout: 8000,
    closeOnClick: true,
    resetOnHover: false,
    transitionIn: 'flipInX',
    transitionOut: 'flipOutX',
});

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
                iziToast.error({
                    message: message.payload,
                });
            }
        }
    });
})