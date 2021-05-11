import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import astor from "@/plugins/astilectron";

import mitt from 'mitt';
const emitter = mitt(); 

const app = createApp(App);
app
  .provide('emitter', emitter)
  .use(astor, {
      debug: true,
    })
  .use(router)
  .mount("#app");
