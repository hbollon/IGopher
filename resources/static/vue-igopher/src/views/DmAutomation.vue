<template>
  <DownloadTracking v-if="this.downloading" @done="this.successDl()" @error="this.errorDl()"/>
  <DmAutomationPanel @showDlModal="this.showModalComp()"/>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import DmAutomationPanel from "@/components/DmAutomationPanel.vue";
import DownloadTracking from "@/components/DownloadTracking.vue";
import * as config from "@/config";

@Options({
  title: "DM Automation",
  components: {
    DownloadTracking,
    DmAutomationPanel,
  },
  data() {
    return {
      downloading: false,
    }
  },
  methods: {
    showModalComp() {
      this.downloading = true;
    },
    dissmissModalComp() {
      this.downloading = false;
    },
    errorDl() {
      this.dissmissModalComp()
      config.Toast.fire({
        icon: "error",
        title: "Error during bot launch! Check logs tab for more details.",
      });
    },
    successDl() {
      this.dissmissModalComp()
      config.Toast.fire({
        icon: "success",
        title: "Bot successfully launched!",
      });
    },
  },
})
export default class DmAutomation extends Vue {}

</script>

<style lang="scss"></style>
