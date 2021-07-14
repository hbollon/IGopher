import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import astor from "@/plugins/astilectron";
import titleMixin from "@/mixins/titleMixin";

import mitt, { Emitter } from 'mitt';
const emitter: Emitter = mitt();
export default emitter;

import "@/bootstrap/css/bootstrap.min.css";
import "@/bootstrap/js/bootstrap.min.js";

const app = createApp(App);
app
  .provide('emitter', emitter)
  .use(astor, {
      debug: true,
      emitter: emitter
    })
  .mixin(titleMixin)
  .use(router)
  .mount("#app");
