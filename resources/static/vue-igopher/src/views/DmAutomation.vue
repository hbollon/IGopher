<template>
  <LateralNav />
  <div class="d-flex flex-column" id="content-wrapper">
    <div id="content">
      <NavBar />
      <DmAutomationPanel />
    </div>
    <Footer />
  </div>
  <a class="border rounded scroll-to-top" href="#page-top"
    ><i class="fas fa-angle-up"></i
  ></a>
</template>

<script>
import { inject } from 'vue'
import LateralNav from "@/components/LateralNav.vue";
import NavBar from "@/components/NavBar.vue";
import Footer from "@/components/Footer.vue";
import DmAutomationPanel from "@/components/DmAutomationPanel.vue";
import * as config from "@/config"
import "@/theme"

export default {
  components: {
    LateralNav,
    NavBar,
    Footer,
    DmAutomationPanel,
  },
  mounted() {
    const astor = inject('astor');

    config.ready(() => {
      document.addEventListener("astilectron-ready", function() {
        config.getIgopherConfig(astor, fillInputs);
        let dmBotLaunchBtn = document.querySelector("#dmBotLaunchBtn");
        let dmBotLaunchIcon = document.querySelector("#dmBotLaunchIcon");
        let dmBotLaunchSpan = document.querySelector("#dmBotLaunchSpan");
        let dmBotHotReloadBtn = document.querySelector("#dmBotHotReloadBtn");
        let dmBotHotReloadIcn = document.querySelector("#dmBotHotReloadIcn");

        // Dynamics buttons inits
        var dmBotRunning = sessionStorage.getItem("botState");
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
            dmBotRunning === false ||
            dmBotRunning === null
          ) {
            astor.trigger("launchDmBot", {}, function(message) {
              if (message.status === config.SUCCESS) {
                config.iziToast.success({
                  message: message.msg,
                });

                dmBotRunning = true;
                dmBotLaunchBtn.classList.add("btn-danger");
                dmBotLaunchBtn.classList.remove("btn-success");
                dmBotLaunchIcon.classList.add("fa-skull-crossbones");
                dmBotLaunchIcon.classList.remove("fa-rocket");
                dmBotLaunchSpan.innerHTML = "Stop !";
                sessionStorage.setItem("botState", true);
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
              message
            ) {
              if (message.status === config.SUCCESS) {
                config.iziToast.success({
                  message: message.msg,
                });
                dmBotRunning = false;
                dmBotLaunchBtn.classList.add("btn-success");
                dmBotLaunchBtn.classList.remove("btn-danger");
                dmBotLaunchIcon.classList.add("fa-rocket");
                dmBotLaunchIcon.classList.remove("fa-spinner", "fa-spin");
                dmBotLaunchSpan.innerHTML = "Launch !";
                sessionStorage.setItem("botState", false);
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
            message
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
            let message = {msg: "dmSettingsForm"};
            let form = document.querySelector("#dmSettingsForm");
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              if (typeof content !== "undefined") {
                let formData = new FormData(form);
                message.payload = config.serialize(formData);
              }
              astor.trigger(message.msg, message.payload, function(message) {
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
            let message = {msg: "dmUserScrappingSettingsForm"};
            let form = document.querySelector("#dmUserScrappingSettingsForm");
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              if (typeof content !== "undefined") {
                let formData = new FormData(form);
                message.payload = config.serialize(formData);
              }
              astor.trigger(message.msg, message.payload, function(message) {
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
  },
};

function fillInputs() {
  document.getElementById(
    "dmTemplates"
  ).value = config.igopherConfig.auto_dm.dmTemplates.join(";");
  document.getElementById("greetingTemplate").value =
    config.igopherConfig.auto_dm.greeting.greetingTemplate;
  document.getElementById(
    "srcUsers"
  ).value = config.igopherConfig.scrapper.srcUsers.join(";");
  document.getElementById("scrappingQuantity").value =
    config.igopherConfig.scrapper.scrappingQuantity;
}
</script>
