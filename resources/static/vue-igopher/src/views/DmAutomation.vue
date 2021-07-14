<template>
  <DmAutomationPanel />
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { inject, onActivated } from 'vue'
import DmAutomationPanel from "@/components/DmAutomationPanel.vue";
import * as config from "@/config"
import { Astor } from "@/plugins/astilectron";

@Options({
  title: "DM Automation",
  components: {
    DmAutomationPanel,
  },
  mounted() {
    console.log("dmautomation")
    const astor: Astor = inject('astor');
    
    config.ready(() => {
      astor.onIsReady(function() {
        config.getIgopherConfig(astor, fillInputs);
        const dmBotLaunchBtn = document.querySelector("#dmBotLaunchBtn");
        const dmBotLaunchIcon = document.querySelector("#dmBotLaunchIcon");
        const dmBotLaunchSpan = document.querySelector("#dmBotLaunchSpan");
        const dmBotHotReloadBtn = document.querySelector("#dmBotHotReloadBtn");
        const dmBotHotReloadIcn = document.querySelector("#dmBotHotReloadIcn");

        // Dynamics buttons inits
        let dmBotRunning = sessionStorage.getItem("botState");
        if (dmBotRunning === "false" || dmBotRunning === null) {
          dmBotLaunchBtn.classList.add("btn-success");
          dmBotLaunchBtn.classList.remove("btn-danger");
          dmBotLaunchIcon.classList.add("fa-rocket");
          dmBotLaunchIcon.classList.remove("fa-spinner", "fa-spin");
          dmBotLaunchSpan.innerHTML = "Launch !";
        } else {
          dmBotLaunchBtn.classList.add("btn-danger");
          dmBotLaunchBtn.classList.remove("btn-success");
          dmBotLaunchIcon.classList.add("fa-skull-crossbones");
          dmBotLaunchIcon.classList.remove("fa-rocket");
          dmBotLaunchSpan.innerHTML = "Stop !";
        }

        /// Buttons
        dmBotLaunchBtn.addEventListener("click", function() {
          if (
            dmBotRunning === "false" ||
            dmBotRunning === null
          ) {
            astor.trigger("launchDmBot", {}, function(message: any) {
              if (message.status === config.SUCCESS) {
                config.iziToast.success({
                  message: message.msg,
                });

                dmBotRunning = "true";
                dmBotLaunchBtn.classList.add("btn-danger");
                dmBotLaunchBtn.classList.remove("btn-success");
                dmBotLaunchIcon.classList.add("fa-skull-crossbones");
                dmBotLaunchIcon.classList.remove("fa-rocket");
                dmBotLaunchSpan.innerHTML = "Stop !";
                sessionStorage.setItem("botState", "true");
              } else {
                config.iziToast.error({
                  message: message.msg,
                });
              }
            });
          } else {
            dmBotLaunchIcon.classList.add("fa-spinner", "fa-spin");
            dmBotLaunchIcon.classList.remove("fa-skull-crossbones");
            config.iziToast.info({
              message:
                "Stop procedure launched, the bot will stop once the current action is finished.",
            });
            astor.trigger("stopDmBot", {}, function(
              message: any
            ) {
              if (message.status === config.SUCCESS) {
                config.iziToast.success({
                  message: message.msg,
                });
                dmBotRunning = "false";
                dmBotLaunchBtn.classList.add("btn-success");
                dmBotLaunchBtn.classList.remove("btn-danger");
                dmBotLaunchIcon.classList.add("fa-rocket");
                dmBotLaunchIcon.classList.remove("fa-spinner", "fa-spin");
                dmBotLaunchSpan.innerHTML = "Launch !";
                sessionStorage.setItem("botState", "false");
              } else {
                dmBotLaunchIcon.classList.add("fa-skull-crossbones");
                dmBotLaunchIcon.classList.remove("fa-spinner", "fa-spin");
                config.iziToast.error({
                  message: message.msg,
                });
              }
            });
          }
        });

        dmBotHotReloadBtn.addEventListener("click", function() {
          dmBotHotReloadIcn.classList.add("fa-spinner", "fa-spin");
          dmBotHotReloadIcn.classList.remove("fa-fire");
          config.iziToast.info({
            message:
              "Hot reload launched, the bot will update once the current action is finished.",
          });
          astor.trigger("hotReloadBot", {}, function(
            message: any
          ) {
            if (message.status === config.SUCCESS) {
              config.iziToast.success({
                message: message.msg,
              });
              dmBotHotReloadIcn.classList.add("fa-fire");
              dmBotHotReloadIcn.classList.remove("fa-spinner", "fa-spin");
            } else {
              config.iziToast.error({
                message: message.msg,
              });
              dmBotHotReloadIcn.classList.add("fa-fire");
              dmBotHotReloadIcn.classList.remove("fa-spinner", "fa-spin");
            }
          });
        });

        /// Forms
        // Dm automation view
        document
          .querySelector("#dmSettingsFormBtn")
          .addEventListener("click", function(e) {
            let message = {msg: "dmSettingsForm", payload: {}};
            let form = document.querySelector("#dmSettingsForm") as HTMLFormElement;
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              let formData = new FormData(form);
              message.payload = config.serialize(formData);
              astor.trigger(message.msg, message.payload, function(message: any) {
                if (message.status === config.SUCCESS) {
                  config.iziToast.success({
                    message: message.msg,
                  });
                } else {
                  config.iziToast.error({
                    title: "Error during settings saving!",
                    message: message.msg,
                  });
                }
              });
            }

            form.classList.add("was-validated");
          });

        document
          .querySelector("#dmUserScrappingSettingsFormBtn")
          .addEventListener("click", function(e) {
            let message = {msg: "dmUserScrappingSettingsForm", payload: {}};
            let form = document.querySelector("#dmUserScrappingSettingsForm") as HTMLFormElement;
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              let formData = new FormData(form);
              message.payload = config.serialize(formData);
              astor.trigger(message.msg, message.payload, function(message: any) {
                if (message.status === config.SUCCESS) {
                  config.iziToast.success({
                    message: message.msg,
                  });
                } else {
                  config.iziToast.error({
                    title: "Error during settings saving!",
                    message: message.msg,
                  });
                }
              });
            }

            form.classList.add("was-validated");
          });
      });
    });
  }
})
export default class DmAutomation extends Vue {}

function fillInputs() {
  const dmTemplatesField = document.getElementById(
    "dmTemplates"
  ) as HTMLTextAreaElement
  dmTemplatesField.value = config.igopherConfig.auto_dm.dmTemplates.join(";");

  const greetingTemplateField = document.getElementById(
    "greetingTemplate"
  ) as HTMLInputElement
  greetingTemplateField.value = config.igopherConfig.auto_dm.greeting.greetingTemplate;

  const srcUsersField = document.getElementById(
    "srcUsers"
  ) as HTMLInputElement
  srcUsersField.value = config.igopherConfig.scrapper.srcUsers.join(";");

  const scrappingQuantityField = document.getElementById(
    "scrappingQuantity"
  ) as HTMLInputElement
  scrappingQuantityField.value = config.igopherConfig.scrapper.scrappingQuantity;
}
</script>

<style lang="scss">

</style>