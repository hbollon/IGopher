<template>
  <SettingsPanel />
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { inject } from "vue";
import SettingsPanel from "@/components/SettingsPanel.vue";
import * as config from "@/config";
import { Astor } from "@/plugins/astilectron";

@Options({
  title: "Settings",
  components: {
    SettingsPanel,
  },
  mounted() {
    const astor: Astor = inject("astor");

    config.ready(() => {
      astor.onIsReady(function() {
        config.getIgopherConfig(astor, fillInputs);

        /// Buttons
        document
          .querySelector("#resetGlobalDefaultSettingsBtn")
          .addEventListener("click", function() {
            astor.trigger("resetGlobalDefaultSettings", {}, function(
              message: any
            ) {
              if (message.status === config.SUCCESS) {
                config.Toast.fire({
                  icon: "success",
                  title: message.msg,
                });
                config.getIgopherConfig(astor, fillInputs);
              } else {
                config.Toast.fire({
                  icon: "error",
                  title: "Unknown error during global settings reset",
                });
              }
            });
          });

        document
          .querySelector("#clearBotDataBtn")
          .addEventListener("click", function() {
            astor.trigger("clearAllData", {}, function(message: any) {
              if (message.status === config.SUCCESS) {
                config.Toast.fire({
                  icon: "success",
                  title: message.msg,
                });
                config.getIgopherConfig(astor, fillInputs);
              } else {
                config.Toast.fire({
                  icon: "error",
                  title: message.msg,
                });
              }
            });
          });

        // document.querySelector("#reinstallDependenciesBtn").addEventListener("click", function() {
        //     astilectron.sendMessage({ "msg": "reinstallDependencies" }, function(message) {
        //         if (message.status === SUCCESS) {
        //             iziToast.success({
        //     message: message.msg,
        // });
        //         } else {
        //             toastr.error('Unknown error during dependencies reinstallation');
        //         }
        //     });
        // });

        /// Forms
        // Settings view
        document
          .querySelector("#igCredentialsFormBtn")
          .addEventListener("click", function(e) {
            let message = { msg: "igCredentialsForm", payload: {} };
            let form = document.querySelector(
              "#igCredentialsForm"
            ) as HTMLFormElement;
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              let formData = new FormData(form);
              message.payload = config.serialize(formData);
              astor.trigger(message.msg, message.payload, function(
                message: any
              ) {
                if (message.status === config.SUCCESS) {
                  config.Toast.fire({
                    icon: "success",
                    title: message.msg,
                  });
                } else {
                  config.Toast.fire({
                    icon: "error",
                    title: "Error during settings saving: " + message.msg,
                  });
                }
              });
            }

            form.classList.add("was-validated");
          });

        document
          .querySelector("#quotasFormBtn")
          .addEventListener("click", function(e) {
            let message = { msg: "quotasForm", payload: {} };
            let form = document.querySelector("#quotasForm") as HTMLFormElement;
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              let formData = new FormData(form);
              message.payload = config.serialize(formData);
              astor.trigger(message.msg, message.payload, function(
                message: any
              ) {
                if (message.status === config.SUCCESS) {
                  config.Toast.fire({
                    icon: "success",
                    title: message.msg,
                  });
                } else {
                  config.Toast.fire({
                    icon: "error",
                    title: "Error during settings saving: " + message.msg,
                  });
                }
              });
            }

            form.classList.add("was-validated");
          });

        document
          .querySelector("#schedulerFormBtn")
          .addEventListener("click", function(e) {
            let message = { msg: "schedulerForm", payload: {} };
            let form = document.querySelector(
              "#schedulerForm"
            ) as HTMLFormElement;
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              let formData = new FormData(form);
              message.payload = config.serialize(formData);
              astor.trigger(message.msg, message.payload, function(
                message: any
              ) {
                if (message.status === config.SUCCESS) {
                  config.Toast.fire({
                    icon: "success",
                    title: message.msg,
                  });
                } else {
                  config.Toast.fire({
                    icon: "error",
                    title: "Error during settings saving: " + message.msg,
                  });
                }
              });
            }

            form.classList.add("was-validated");
          });

        document
          .querySelector("#blacklistFormBtn")
          .addEventListener("click", function(e) {
            let message = { msg: "blacklistForm", payload: {} };
            let form = document.querySelector(
              "#blacklistForm"
            ) as HTMLFormElement;
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              let formData = new FormData(form);
              message.payload = config.serialize(formData);
              astor.trigger(message.msg, message.payload, function(
                message: any
              ) {
                if (message.status === config.SUCCESS) {
                  config.Toast.fire({
                    icon: "success",
                    title: message.msg,
                  });
                } else {
                  config.Toast.fire({
                    icon: "error",
                    title: "Error during settings saving: " + message.msg,
                  });
                }
              });
            }

            form.classList.add("was-validated");
          });

        document
          .querySelector("#proxyFormBtn")
          .addEventListener("click", function(e) {
            let message = { msg: "proxyForm", payload: {} };
            let form = document.querySelector("#proxyForm") as HTMLFormElement;
            if (!form.checkValidity()) {
              e.preventDefault();
              e.stopPropagation();
            } else {
              let formData = new FormData(form);
              message.payload = config.serialize(formData);
              astor.trigger(message.msg, message.payload, function(
                message: any
              ) {
                if (message.status === config.SUCCESS) {
                  config.Toast.fire({
                    icon: "success",
                    title: message.msg,
                  });
                } else {
                  config.Toast.fire({
                    icon: "error",
                    title: "Error during settings saving: " + message.msg,
                  });
                }
              });
            }

            form.classList.add("was-validated");
          });

        document
          .querySelector("#proxyAuthCheck")
          .addEventListener("change", function() {
            let checkbox = document.querySelector(
              "#proxyAuthCheck"
            ) as HTMLInputElement;
            let divAuthInputs = document.querySelector("#proxyAuthInputs");
            if (checkbox.checked) {
              checkbox.value = "true";
              if (divAuthInputs.classList.contains("d-none"))
                divAuthInputs.classList.remove("d-none");
            } else {
              checkbox.value = "false";
              if (!divAuthInputs.classList.contains("d-none"))
                divAuthInputs.classList.add("d-none");
            }

            let authInputs = document.querySelectorAll(".auth-proxy");
            for (const element of authInputs as NodeListOf<HTMLInputElement>) {
              if (checkbox.checked) {
                if (element.required !== true) element.required = true;
              } else {
                if (element.required !== false) element.required = false;
              }
            }
          });
      });
    });
  },
})
export default class Settings extends Vue {}

function fillInputs() {
  const dmHourField = document.getElementById(
    "dmHourInput"
  ) as HTMLTextAreaElement;
  dmHourField.value = config.igopherConfig.quotas.dmHour;

  const dmDayField = document.getElementById("dmDayInput") as HTMLInputElement;
  dmDayField.value = config.igopherConfig.quotas.dmDay;

  const quotasRadio = document.getElementById(
    config.igopherConfig.quotas.quotasActivation == "true" ? "quotasRadioEnabled" : "quotasRadioDisabled"
  ) as HTMLInputElement;
  quotasRadio.checked = true;

  const beginAtField = document.getElementById(
    "beginAtInput"
  ) as HTMLInputElement;
  beginAtField.value = config.igopherConfig.schedule.beginAt;

  const endAtField = document.getElementById("endAtInput") as HTMLInputElement;
  endAtField.value = config.igopherConfig.schedule.endAt;

  const schedulerRadio = document.getElementById(
    config.igopherConfig.schedule.scheduleActivation == "true" ? "schedulerRadioEnabled" : "schedulerRadioDisabled"
  ) as HTMLInputElement;
  schedulerRadio.checked = true;

  const blacklistRadio = document.getElementById(
    config.igopherConfig.blacklist.blacklistActivation == "true" ? "blacklistRadioEnabled" : "blacklistRadioDisabled"
  ) as HTMLInputElement;
  blacklistRadio.checked = true;

  const ipProxyField = document.getElementById("ipInput") as HTMLInputElement;
  ipProxyField.value = config.igopherConfig.webdriver.proxy.ip;

  const portProxyField = document.getElementById(
    "portInput"
  ) as HTMLInputElement;
  portProxyField.value = config.igopherConfig.webdriver.proxy.port;

  if (config.igopherConfig.webdriver.proxy.auth == "true") {
    document.getElementById("proxyAuthCheck").click();

    const proxyUsernameField = document.getElementById(
      "proxyUsernameInput"
    ) as HTMLInputElement;
    proxyUsernameField.value = config.igopherConfig.webdriver.proxy.username;

    const proxyPasswordField = document.getElementById(
      "proxyPasswordInput"
    ) as HTMLInputElement;
    proxyPasswordField.value = config.igopherConfig.webdriver.proxy.password;
  }

  const proxyRadio = document.getElementById(
    config.igopherConfig.webdriver.proxy.proxyActivation == "true" ? "proxyRadioEnabled" : "proxyRadioDisabled"
  ) as HTMLInputElement;
  proxyRadio.checked = true;

  const usernameField = document.getElementById(
    "usernameInput"
  ) as HTMLInputElement;
  usernameField.value = config.igopherConfig.account.username;

  const passwordField = document.getElementById(
    "passwordInput"
  ) as HTMLInputElement;
  passwordField.value = config.igopherConfig.account.password;
}
</script>

<style lang="scss"></style>
