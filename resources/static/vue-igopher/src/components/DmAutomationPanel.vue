<template>
  <div class="container-fluid">
    <div class="d-flex d-xl-flex align-items-xl-center" style="margin-bottom: 20px;">
      <h3 class="text-dark d-flex d-xl-flex justify-content-start mb-1" style="margin-right: auto;">DM Automation</h3>
      <a class="btn btn-warning d-flex justify-content-end btn-icon-split" role="button"
      id="dmBotHotReloadBtn" style="margin-right: 5px;" @click="onClickHotReloadBtn"
        ><span class="text-white-50 icon"><i class="fas fa-fire" id="dmBotHotReloadIcn"></i></span><span class="text-white text">Hot reload</span></a
      ><a class="btn btn-success d-flex justify-content-end btn-icon-split" role="button" id="dmBotLaunchBtn" @click="onClickDmBotLaunchBtn"
        ><span class="text-white-50 icon"><i class="fas fa-rocket" id="dmBotLaunchIcon"></i></span
        ><span class="text-white text" id="dmBotLaunchSpan">Launch !</span></a
      >
    </div>
    <div class="row">
      <div class="col-lg-6 col-xl-7">
        <div class="card shadow mb-4">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h6 class="text-primary fw-bold m-0">Module Settings</h6>
          </div>
          <div class="card-body">
            <form id="dmSettingsForm" novalidate="">
              <div class="form-group mb-3">
                <label class="form-label">Message templates (separated by a semicolon):</label
                ><textarea class="form-control" id="dmTemplates" name="dmTemplates" required="">Hey ! What's up?</textarea>
                <div class="invalid-feedback">Invalid input!</div>
              </div>
              <hr />
              <div class="form-group mb-3">
                <label class="form-label">Greeting:&nbsp;</label>
                <div class="form-check">
                  <input type="radio" class="form-check-input" id="greetingRadioEnabled" value="true" name="greetingActivation" required /><label
                    class="form-check-label"
                    for="greetingRadioEnabled"
                    >Enabled</label
                  >
                </div>
                <div class="form-check disabled">
                  <input type="radio" class="form-check-input" id="greetingRadioDisabled" value="false" name="greetingActivation" required checked /><label
                    class="form-check-label"
                    for="greetingRadioDisabled"
                    >Disabled</label
                  >
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
                <small class="form-text text-muted"
                  >Will add custom greeting with user's username before message templates: ("Hello" as greeting template will produce -&gt; "Hello
                  &lt;username&gt;, [dm template]")</small
                >
              </div>
              <div class="form-group mb-3">
                <label class="form-label">Greeting template:</label
                ><input class="form-control" type="text" id="greetingTemplate" required="" name="greetingTemplate" value="Hello" />
                <div class="invalid-feedback">Invalid input!</div>
              </div>
              <div>
                <button class="btn btn-primary" id="dmSettingsFormBtn" type="button" @click="submitDmSettings">Save</button
                ><button class="btn btn-secondary" type="reset" style="margin-left: 5px;">Reset to default</button>
              </div>
            </form>
          </div>
        </div>
      </div>
      <div class="col-lg-6 col-xl-5">
        <div class="card shadow mb-4">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h6 class="text-primary fw-bold m-0">Users Scrapping</h6>
          </div>
          <div class="card-body">
            <form id="dmUserScrappingSettingsForm" novalidate="">
              <div class="form-group mb-3">
                <label class="form-label">Source Users:</label>
                <Multiselect
                  id="srcUsersSelect"
                  name="srcUsers"
                  v-model="srcUsrMultiSelect.value"
                  v-bind="srcUsrMultiSelect"
                ></Multiselect>
                <div class="invalid-feedback">Invalid input!</div>
                <small class="form-text text-muted"
                  >This is a tags input, just hit enter to validate a username.</small
                >
              </div>
              <div class="form-group mb-3">
                <label class="form-label">Fetch from:</label>
                <div class="form-check">
                  <input type="radio" class="form-check-input" id="formCheck-1" value="1" name="fetchUsersFrom" required checked /><label
                    class="form-check-label"
                    for="formCheck-1"
                    >Users followers</label
                  >
                </div>
                <div class="form-check disabled">
                  <input type="radio" class="form-check-input" id="formCheck-2" value="2" name="fetchUsersFrom" required disabled /><label
                    class="form-check-label"
                    for="formCheck-2"
                    >Users following (not available yet)</label
                  >
                </div>
                <div class="form-check disabled">
                  <input type="radio" class="form-check-input" id="formCheck-3" value="3" name="fetchUsersFrom" required disabled /><label
                    class="form-check-label"
                    for="formCheck-3"
                    >Users post likers (not available yet)</label
                  >
                </div>
                <div class="form-check disabled">
                  <input type="radio" class="form-check-input" id="formCheck-4" value="4" name="fetchUsersFrom" required disabled /><label
                    class="form-check-label"
                    for="formCheck-4"
                    >Users post commentators (not available yet)</label
                  >
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div class="form-group mb-3">
                <label class="form-label">Fetch quantity:</label
                ><input
                  class="form-control"
                  type="number"
                  id="scrappingQuantity"
                  min="1"
                  max="5000"
                  required=""
                  name="scrappingQuantity"
                  value="200"
                />
                <div class="invalid-feedback">Invalid input!</div>
              </div>
              <div>
                <button class="btn btn-primary" id="dmUserScrappingSettingsFormBtn" type="button" @click="submitScrapperSettings">Save</button
                ><button class="btn btn-secondary" type="reset" style="margin-left: 5px;">Reset to default</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Astor } from "@/plugins/astilectron";
import { Vue, Options } from "vue-class-component";
import { inject } from 'vue';
import * as config from "@/config";

import Multiselect from '@vueform/multiselect'

@Options({
  components: {
    Multiselect,
  },
  data() {
    return {
      astor: Astor,
      srcUsrMultiSelect: {
        mode: 'tags',
        value: [],
        options: [],
        object: false,
        required: true,
        searchable: true,
        createTag: true,
        showOptions: false,
        caret: false,
      },
    }
  },
  computed: {
    sessionBotState: {
      get(): string {
        return sessionStorage.getItem("botState");
      },
      set(value: string) {
        return sessionStorage.setItem("botState", value);
      },
    }
  },
  methods: {
    fillInputs() {
      const dmTemplatesField = document.getElementById(
        "dmTemplates"
      ) as HTMLTextAreaElement;
      dmTemplatesField.value = config.igopherConfig.auto_dm.dmTemplates.join(";");

      const greetingTemplateField = document.getElementById(
        "greetingTemplate"
      ) as HTMLInputElement;
      greetingTemplateField.value =
        config.igopherConfig.auto_dm.greeting.greetingTemplate;

      const greetingRadio = document.getElementById(
        config.igopherConfig.auto_dm.greeting.greetingActivation == "true" ? "greetingRadioEnabled" : "greetingRadioDisabled"
      ) as HTMLInputElement;
      greetingRadio.checked = true;

      const scrappingQuantityField = document.getElementById(
        "scrappingQuantity"
      ) as HTMLInputElement;
      scrappingQuantityField.value =
        config.igopherConfig.scrapper.scrappingQuantity;

      config.igopherConfig.scrapper.srcUsers.forEach((username: string) => {
        this.srcUsrMultiSelect.options.push({ value: username, label: username });
        this.srcUsrMultiSelect.value.push(username);
      });
    },
    onClickDmBotLaunchBtn(): void {
      const dmBotLaunchBtn = document.querySelector("#dmBotLaunchBtn");
      const dmBotLaunchIcon = document.querySelector("#dmBotLaunchIcon");
      const dmBotLaunchSpan = document.querySelector("#dmBotLaunchSpan");
      if (this.sessionBotState === "false" || this.sessionBotState === null) {
        this.astor.trigger("launchDmBot", {}, (message: any) => {
          if (message.status === config.SUCCESS) {
            this.$emit("showDlModal");
            this.sessionBotState = "true";
            dmBotLaunchBtn.classList.add("btn-danger");
            dmBotLaunchBtn.classList.remove("btn-success");
            dmBotLaunchIcon.classList.add("fa-skull-crossbones");
            dmBotLaunchIcon.classList.remove("fa-rocket");
            dmBotLaunchSpan.innerHTML = "Stop !";
          } else {
            config.Toast.fire({
              icon: "error",
              title: message.msg,
            });
          }
        });
      } else {
        dmBotLaunchIcon.classList.add("fa-spinner", "fa-spin");
        dmBotLaunchIcon.classList.remove("fa-skull-crossbones");
        config.Toast.fire({
          icon: "info",
          title:
            "Stop procedure launched, the bot will stop once the current action is finished.",
        });

        this.astor.trigger("stopDmBot", {}, (message: any) => {
          if (message.status === config.SUCCESS) {
            config.Toast.fire({
              icon: "success",
              title: message.msg,
            });

            this.sessionBotState = "false";
            dmBotLaunchBtn.classList.add("btn-success");
            dmBotLaunchBtn.classList.remove("btn-danger");
            dmBotLaunchIcon.classList.add("fa-rocket");
            dmBotLaunchIcon.classList.remove("fa-spinner", "fa-spin");
            dmBotLaunchSpan.innerHTML = "Launch !";
          } else {
            dmBotLaunchIcon.classList.add("fa-skull-crossbones");
            dmBotLaunchIcon.classList.remove("fa-spinner", "fa-spin");
            config.Toast.fire({
              icon: "error",
              title: message.msg,
            });
          }
        });
      }
    },
    onClickHotReloadBtn(): void {
      const dmBotHotReloadIcn = document.querySelector("#dmBotHotReloadIcn");
      dmBotHotReloadIcn.classList.add("fa-spinner", "fa-spin");
      dmBotHotReloadIcn.classList.remove("fa-fire");
      config.Toast.fire({
        icon: "info",
        title:
          "Hot reload launched, the bot will update once the current action is finished.",
      });
      this.astor.trigger("hotReloadBot", {}, function(message: any) {
        if (message.status === config.SUCCESS) {
          config.Toast.fire({
            icon: "success",
            title: message.msg,
          });
          dmBotHotReloadIcn.classList.add("fa-fire");
          dmBotHotReloadIcn.classList.remove("fa-spinner", "fa-spin");
        } else {
          config.Toast.fire({
            icon: "error",
            title: message.msg,
          });
          dmBotHotReloadIcn.classList.add("fa-fire");
          dmBotHotReloadIcn.classList.remove("fa-spinner", "fa-spin");
        }
      });
    },
    submitDmSettings(e: Event): void {
      let message = { msg: "dmSettingsForm", payload: {} };
      let form = document.querySelector(
        "#dmSettingsForm"
      ) as HTMLFormElement;
      if (!form.checkValidity()) {
        e.preventDefault();
        e.stopPropagation();
      } else {
        let formData = new FormData(form);
        message.payload = config.serialize(formData);
        this.astor.trigger(message.msg, message.payload, function(
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
    },
    submitScrapperSettings(e: Event): void {
      let message: any = {};
      message.msg = "dmUserScrappingSettingsForm"
      let form = document.querySelector(
        "#dmUserScrappingSettingsForm"
      ) as HTMLFormElement;
      if (!form.checkValidity()) {
        e.preventDefault();
        e.stopPropagation();
      } else {
        let formData = new FormData(form);
        message.payload = config.serialize(formData);
        message.payload.srcUsers = this.srcUsrMultiSelect.value.join(';');
        this.astor.trigger(message.msg, message.payload, function(
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
    },
  },
  mounted() {
    this.astor = inject("astor");

    config.ready(() => {
      this.astor.onIsReady(() => {
        config.getIgopherConfig(this.astor, this.fillInputs.bind(this));
        const dmBotLaunchBtn = document.querySelector("#dmBotLaunchBtn");
        const dmBotLaunchIcon = document.querySelector("#dmBotLaunchIcon");
        const dmBotLaunchSpan = document.querySelector("#dmBotLaunchSpan");

        // Dynamics buttons inits
        if (this.sessionBotState === "false" || this.sessionBotState === null) {
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
      });
    });
  }
})
export default class DmAutomationPanel extends Vue {}
</script>

<style src="@vueform/multiselect/themes/default.css"></style>
