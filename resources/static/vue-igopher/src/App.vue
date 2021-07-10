<template>
  <LateralNav />
  <div class="d-flex flex-column" id="content-wrapper">
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
</style>

<script lang="ts">
import { defineComponent } from "vue";
import "@/theme"

export default defineComponent({
  name: "App"
});
</script>
