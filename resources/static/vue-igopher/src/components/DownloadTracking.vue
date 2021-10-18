<template>
  <div
    class="modal fade"
    id="dlModal"
    tabindex="-1"
    aria-labelledby="dlModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog modal-dialog-centered" style="width: 40vw; height: 50vw;">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="dlModalLabel">Dependencies Manager</h5>
        </div>
        <div class="modal-body text-center align-middle">
          Downloading required dependencies, please wait...
          <hr>
          <div
            class="row my-2 px-1"
            v-for="(dl, filename) in downloadTracking"
            :key="filename"
          >
            <div class="col-3" style="font-size: 12px">
              {{ filename }}
            </div>
            <div class="col-9 ps-2">
              <div class="progress" style="height: 20px;">
                <div
                  :id="'bar-' + filename"
                  class="progress-bar progress-bar-striped progress-bar-animated"
                  role="progressbar"
                  style="width: 0%;"
                  :aria-valuenow="dl.progress"
                  aria-valuemin="0"
                  aria-valuemax="100"
                >
                  {{ dl.progress }}%
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Options } from "vue-class-component";
import { inject } from "vue";
import { Astor } from "../plugins/astilectron";
import { bootstrap } from "@/config";

@Options({
  data() {
    return {
      astor: Astor,
      dlmodal: {},
      downloadTracking: {},
    };
  },
  methods: {
    updateProgress(payload: any) {
      this.downloadTracking = payload.payload;
      for (const key in this.downloadTracking) {
        const progress = Math.floor(this.downloadTracking[key].Progress);
        const bar = document.getElementById("bar-" + key);
        bar.setAttribute("aria-valuenow", "" + progress);
        bar.style.width = progress + "%";
        bar.innerHTML = progress + "%";
      }
    },
    closeModal() {
      this.dlModal.dispose();
      document.body.classList.remove("modal-open");
      const backdrop = document.querySelector(".modal-backdrop");
      backdrop.parentNode.removeChild(backdrop);
      this.astor.remove("downloads tracking", this.updateProgress);
      this.$emit("dlDone");
    },
  },
  mounted() {
    this.astor = inject("astor");
    this.astor.listen(
      "downloads tracking",
      this.updateProgress.bind(this),
      false
    );
    this.astor.listen("downloads done", this.closeModal.bind(this), true);

    this.dlModal = new bootstrap.Modal(document.getElementById("dlModal"), {
      backdrop: "static",
      keyboard: false,
    });
    this.dlModal.toggle();
  },
})
export default class DownloadTracking extends Vue {}
</script>

<style></style>
