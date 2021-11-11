<template>
  <LateralNav />
  <div class="d-flex flex-column w-100" id="content-wrapper">
    <div id="content">
      <NavBar />
      <router-view />
    </div>
    <Footer />
  </div>
  <a class="border rounded back-to-top"
    ><i class="fas fa-angle-up"></i
  ></a>
</template>

<script lang="ts">
import { Vue, Options } from "vue-class-component";
import { inject } from "vue";
import { Astor } from "./plugins/astilectron";
import LateralNav from "@/components/LateralNav.vue";
import NavBar from "@/components/NavBar.vue";
import Footer from "@/components/Footer.vue";
import * as config from "@/config";
import "@/theme"

@Options({
  components: {
    LateralNav,
    NavBar,
    Footer,
  },
  data() {
    return {
      astor: Astor,
    }
  },
  mounted() {
    this.astor = inject("astor");
    config.ready(() => {
      this.astor.onIsReady(() => {
        this.astor.listen("bot crash", () => {
          const dmBotLaunchBtn = document.querySelector("#dmBotLaunchBtn");
          const dmBotLaunchIcon = document.querySelector("#dmBotLaunchIcon");
          const dmBotLaunchSpan = document.querySelector("#dmBotLaunchSpan");

          if(dmBotLaunchBtn != undefined && dmBotLaunchIcon != undefined && dmBotLaunchSpan != undefined) {
            dmBotLaunchBtn.classList.add("btn-success");
            dmBotLaunchBtn.classList.remove("btn-danger");
            dmBotLaunchIcon.classList.add("fa-rocket");
            dmBotLaunchIcon.classList.remove("fa-spinner", "fa-spin");
            dmBotLaunchSpan.innerHTML = "Launch !";
          }
          sessionStorage.setItem("botState", "false");
        });
      });
    });
  },
})
export default class App extends Vue {}
</script>

<style lang="scss">
#app {
  font-family: Roboto, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  display: flex;
}

.back-to-top {
  display: none;
  background-color: rgba(90,92,105,.5);
  color: #fff;
  width: 2.75rem;
  height: 2.75rem;
  text-align: center;
  position: fixed;
  bottom: 1rem;
  right: 1rem;
  line-height: 2.35rem;
  transition: background-color .3s, 
    opacity .5s, visibility .5s;
  z-index: 1000;
}

.back-to-top i {
  font-size: 16px;
  font-weight: 800;
  vertical-align: middle;
}

.back-to-top:hover {
  cursor: pointer;
  color: #fff;
  background-color: rgb(68, 68, 68);
}

.colored-toast.swal2-icon-success {
  background-color: #a5dc86 !important;
}

.colored-toast.swal2-icon-error {
  background-color: #f27474 !important;
}

.colored-toast.swal2-icon-warning {
  background-color: #f8bb86 !important;
}

.colored-toast.swal2-icon-info {
  background-color: #3fc3ee !important;
}

.colored-toast.swal2-icon-question {
  background-color: #87adbd !important;
}

.colored-toast .swal2-title {
  color: white;
}

.colored-toast .swal2-close {
  color: white;
}

.colored-toast .swal2-content {
  color: white;
}
</style>