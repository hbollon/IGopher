import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import astor from "@/plugins/astilectron";

import mitt, { Emitter } from 'mitt';
const emitter: Emitter = mitt();
export default emitter;

const app = createApp(App);
app
  .provide('emitter', emitter)
  .use(astor, {
      debug: true,
      emitter: emitter
    })
  .use(router)
  .mount("#app");
