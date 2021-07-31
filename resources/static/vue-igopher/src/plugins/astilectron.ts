declare const astilectron: any;
import { Emitter } from 'mitt';
export class Astor {
    skipWait: boolean;
    debug: boolean;
    isReady: boolean;
    emitter: Emitter;

    constructor(debug: boolean, skipWait: boolean, emitter: Emitter) {
        this.debug = debug !== undefined ? debug : false;
        this.skipWait = skipWait !== undefined ? skipWait : false;
        this.emitter = emitter;
        this.isReady = false;
    }

    init() {
        this.log('init');
        this.isReady = false;

        if (this.skipWait) {
            this.onAstilectronReady();
            return;
        }

        document.addEventListener('astilectron-ready', this.onAstilectronReady.bind(this));
    }

    onAstilectronReady() {
        this.log('astilectron is ready');
        astilectron.onMessage(this.onAstilectronMessage.bind(this));
        this.log('removing ready listener');
        document.removeEventListener('astilectron-ready', this.onAstilectronReady.bind(this));
        this.isReady = true;
    }

    onIsReady(callback: any) {
        const self = this;
        const delay = 100;
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
    }

    onAstilectronMessage(message: any) {
        if (Array.prototype.slice.call(arguments).length == 1) { // eslint-disable-line
            if (message) {
                this.log('GO -> Vue', message);
                this.emit(message.msg, message);
            }
        } else {
            const identifier = message;
            message = Array.prototype.slice.call(arguments)[1]; // eslint-disable-line
            if (message) {
                this.log('GO -> Vue', message);
                this.emit(identifier, message);
            }
        }
    }

    trigger(name: string, payload = {}, callback: any = null) {
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
    }

    listen(name: any, callback: any, once = false) {
        if (once) {
            this.log('listen once', {name: name, callback: callback});
            const wrappedHandler = (evt: any) => {
                callback(evt);
                this.emitter.off(name, wrappedHandler);
            }
            this.emitter.on(name, wrappedHandler);
        } else {
            this.log('listen', {name: name, callback: callback});
            this.emitter.on(name, callback);
        }
    }

    emit(name: any, payload = {}) {
        this.log('EMIT', {name: name, payload: payload});
        this.emitter.emit(name, payload);
    }

    remove(name: any, callback: any) {
        this.emitter.off(name, callback);
    }

    log(message: any, data?: any) {
        if (!this.debug) {
            return;
        }

        if (data) {
            console.log('ASTOR | ' + message, data);
        } else {
            console.log('ASTOR | ' + message);
        }
    }

    getScope() {
        return '#' + Math.random().toString(36).substr(2, 7);
    }
}

export default {
    install (Vue: any, options: any) {
        const { debug, skipWait, emitter } = options;
        const astor: Astor = new Astor(debug, skipWait, emitter)

        Vue.config.globalProperties.$astor = astor;
        Vue.config.globalProperties.$astor.debug = debug;
        Vue.config.globalProperties.$astor.skipWait = skipWait;
        Vue.config.globalProperties.$astor.init();

        Vue.provide('astor', astor);
    }
}