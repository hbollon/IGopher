declare const astilectron: any; // eslint-disable-line

/* eslint-disable */
export default {
    install (Vue: any, options: any) {
        const { debug, skipWait, emitter } = options;

        const astor = {
            skipWait: false,
            debug: false,
            isReady: false,
            init: function() {
                this.log('init');
                this.isReady = false;

                if (skipWait) {
                    this.onAstilectronReady();
                    return;
                }

                document.addEventListener('astilectron-ready', this.onAstilectronReady.bind(this));
            },
            onAstilectronReady: function() {
                this.log('astilectron is ready');
                astilectron.onMessage(this.onAstilectronMessage.bind(this));
                this.log('removing ready listener');
                document.removeEventListener('astilectron-ready', this.onAstilectronReady.bind(this));
                this.isReady = true;
            },
            onIsReady: function(callback: any) {
                let self = this;
                let delay = 100;
                if (!this.isReady) {
                    setTimeout( () => {
                        if (this.isReady) {
                            self.log('astor is ready');
                            callback();
                        } else {
                            self.onIsReady(callback);
                        }
                    }, delay);
                } else {
                    this.log('astor is ready');
                    callback();
                }
            },
            onAstilectronMessage: function(message: any) {
                if (Array.prototype.slice.call(arguments).length == 1) {
                    if (message) {
                        this.log('GO -> Vue', message);
                        this.emit(message.name, message);
                    }
                } else {
                    let identifier = message;
                    message = Array.prototype.slice.call(arguments)[1];
                    if (message) {
                        this.log('GO -> Vue', message);
                        this.emit(identifier, message);
                    }
                }
            },
            trigger: function(name: string, payload = {}, callback = null) {
                let logMessage = 'Vue -> GO';
                let identifier = name;

                if (callback !== null) {
                    logMessage = logMessage + ' (scoped)';
                    identifier = identifier + this.getScope();
                } 

                this.log(logMessage, {name: name, payload: payload});
                if (callback !== null) {
                    this.listen(identifier, callback, true);
                }
                astilectron.sendMessage({msg: name, payload: payload}, this.onAstilectronMessage.bind(this, identifier));
            },
            listen: function(name: any, callback: any, once = false) {
                if (once) {
                    this.log('listen once', {name: name, callback: callback});
                    const wrappedHandler = (evt: any) => {
                        callback(evt);
                        emitter.off(name, wrappedHandler);
                    }
                    emitter.on(name, wrappedHandler);
                } else {
                    this.log('listen', {name: name, callback: callback});
                    emitter.on(name, callback);
                }
            },
            emit: function(name: any, payload = {}) {
                this.log('EMIT', {name: name, payload: payload});
                emitter.emit(name, payload);
            },
            remove: function(name: any, callback: any) {
                emitter.off(name, callback);
            },
            log: function (message: any, data?: any) {
                if (!this.debug) {
                    return;
                }

                if (data) {
                    console.log('ASTOR | ' + message, data);
                } else {
                    console.log('ASTOR | ' + message);
                }
            },
            getScope: function() {
                return '#' + Math.random().toString(36).substr(2, 7);
            }
        }
        
        Vue.config.globalProperties.$astor = astor;
        Vue.config.globalProperties.$astor.debug = debug;
        Vue.config.globalProperties.$astor.skipWait = skipWait;
        Vue.config.globalProperties.$astor.init();

        Vue.provide('astor', astor);
    }   
}