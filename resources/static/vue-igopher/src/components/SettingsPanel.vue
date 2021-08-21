<template>
  <div class="container-fluid">
    <h3 class="text-dark mb-2">Settings</h3>
    <div class="row">
      <div class="col-lg-6 col-xl-6">
        <div class="card shadow mb-4">
          <div
            class="card-header d-flex justify-content-between align-items-center"
          >
            <h6 class="text-primary fw-bold m-0">Interactions quotas</h6>
          </div>
          <div class="card-body">
            <p class="text-center">
              This module manages and limits the number of interactions
              performed per hour and per day, in order to respect the limits
              set. The default settings are those recommended to protect your
              Instagram account from restrictions.<br />This module is
              recommended.
            </p>
            <hr />
            <form id="quotasForm" novalidate="">
              <div class="form-group mb-3">
                <label class="form-label">Activation:&nbsp;</label>
                <div class="form-check">
                  <input
                    type="radio"
                    class="form-check-input"
                    id="quotasRadioEnabled"
                    value="true"
                    name="quotasActivation"
                    required
                    checked
                  /><label class="form-check-label" for="quotasRadioEnabled"
                    >Enabled</label
                  >
                </div>
                <div class="form-check disabled">
                  <input
                    type="radio"
                    class="form-check-input"
                    id="quotasRadioDisabled"
                    value="false"
                    name="quotasActivation"
                    required
                  /><label class="form-check-label" for="quotasRadioDisabled"
                    >Disabled</label
                  >
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <hr />
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Likes per day</span
                  ><input
                    id="likesPerDayInput"
                    type="number"
                    class="form-control"
                    name="likesPerDay"
                    required
                    min="1"
                    value="100"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Likes per hour</span
                  ><input
                    id="likesPerHourInput"
                    type="number"
                    class="form-control"
                    name="likesPerHour"
                    required
                    min="1"
                    value="10"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <hr />
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Follows per day</span
                  ><input
                    id="followsPerDayInput"
                    type="number"
                    class="form-control"
                    name="followsPerDay"
                    required
                    min="1"
                    value="80"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Follows per hour</span
                  ><input
                    id="followsPerHourInput"
                    type="number"
                    class="form-control"
                    name="followsPerHour"
                    required
                    min="1"
                    value="8"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <hr />
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Unfollows per day</span
                  ><input
                    id="unfollowsPerDayInput"
                    type="number"
                    class="form-control"
                    name="unfollowsPerDay"
                    required
                    min="1"
                    value="60"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Unfollows per hour</span
                  ><input
                    id="unfolllowsPerHourInput"
                    type="number"
                    class="form-control"
                    name="unfolllowsPerHour"
                    required
                    min="1"
                    value="6"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <hr />
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Direct messages per day</span
                  ><input
                    id="dmDayInput"
                    type="number"
                    class="form-control"
                    name="dmDay"
                    required
                    min="1"
                    value="40"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Direct messages per hour</span
                  ><input
                    id="dmHourInput"
                    type="number"
                    class="form-control"
                    name="dmHour"
                    required
                    min="1"
                    value="4"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div>
                <button
                  class="btn btn-primary"
                  id="quotasFormBtn"
                  type="button"
                >
                  Save</button
                ><button
                  class="btn btn-secondary"
                  type="reset"
                  style="margin-left: 5px;"
                >
                  Reset to default
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
      <div class="col-lg-6 col-xl-6">
        <div class="card shadow mb-4">
          <div
            class="card-header d-flex justify-content-between align-items-center"
          >
            <h6 class="text-primary fw-bold m-0">Scheduler</h6>
          </div>
          <div class="card-body">
            <p class="text-center">
              This module defines and manages the bot's operating hours.<br />It
              is recommended to activate it and leave a period of inactivity.
            </p>
            <hr />
            <form id="schedulerForm" novalidate="">
              <div class="form-group mb-3">
                <label class="form-label">Activation:&nbsp;</label>
                <div class="form-check">
                  <input
                    type="radio"
                    class="form-check-input"
                    id="schedulerRadioEnabled"
                    value="true"
                    name="scheduleActivation"
                    required
                    checked
                  /><label class="form-check-label" for="schedulerRadioEnabled"
                    >Enabled</label
                  >
                </div>
                <div class="form-check disabled">
                  <input
                    type="radio"
                    class="form-check-input"
                    id="schedulerRadioDisabled"
                    value="false"
                    name="scheduleActivation"
                    required
                  /><label class="form-check-label" for="schedulerRadioDisabled"
                    >Disabled</label
                  >
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <hr />
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Start time</span
                  ><input
                    id="beginAtInput"
                    class="form-control"
                    type="time"
                    name="beginAt"
                    required
                    value="08:00"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">End time</span
                  ><input
                    id="endAtInput"
                    class="form-control"
                    type="time"
                    name="endAt"
                    required
                    value="18:00"
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div>
                <button
                  class="btn btn-primary"
                  id="schedulerFormBtn"
                  type="button"
                >
                  Save</button
                ><button
                  class="btn btn-secondary"
                  type="reset"
                  style="margin-left: 5px;"
                >
                  Reset to default
                </button>
              </div>
            </form>
          </div>
        </div>
        <div class="card shadow mb-4">
          <div
            class="card-header d-flex justify-content-between align-items-center"
          >
            <h6 class="text-primary fw-bold m-0">Users blacklist</h6>
          </div>
          <div class="card-body">
            <p class="text-center">
              This module maintains a blacklist where are added the users with
              whom we have already interacted in order to avoid duplicates and
              spam.
            </p>
            <hr />
            <form id="blacklistForm" novalidate="">
              <div class="form-group mb-3">
                <label class="form-label">Activation:&nbsp;</label>
                <div class="form-check">
                  <input
                    type="radio"
                    class="form-check-input"
                    id="blacklistRadioEnabled"
                    value="true"
                    name="blacklistActivation"
                    required
                    checked
                  /><label class="form-check-label" for="blacklistRadioEnabled"
                    >Enabled</label
                  >
                </div>
                <div class="form-check disabled">
                  <input
                    type="radio"
                    class="form-check-input"
                    id="blacklistRadioDisabled"
                    value="false"
                    name="blacklistActivation"
                    required
                  /><label class="form-check-label" for="blacklistRadioDisabled"
                    >Disabled</label
                  >
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div>
                <button
                  class="btn btn-primary"
                  id="blacklistFormBtn"
                  type="button"
                >
                  Save</button
                ><button
                  class="btn btn-secondary"
                  type="reset"
                  style="margin-left: 5px;"
                >
                  Reset to default
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    <hr />
    <div class="row justify-content-center">
      <div class="col-lg-6 col-xl-6">
        <div class="card shadow mb-4">
          <div
            class="card-header d-flex justify-content-between align-items-center"
          >
            <h6 class="text-primary fw-bold m-0">Proxy</h6>
          </div>
          <div class="card-body">
            <div
              class="text-start"
              style="margin-right: 15px;margin-left: 15px;"
            >
              <p class="text-center">
                This proxy configuration will be used for all webdriver
                operations to ensure the security of your Instagram account.<br />We
                recommend that you use a proxy to further limit Instagram
                limitation risks and / or sanctions.
              </p>
              <hr />
              <form id="proxyForm" novalidate="">
                <div class="form-group mb-3">
                  <div class="input-group">
                    <span class="input-group-text">IP Adress</span
                    ><input
                      id="ipInput"
                      type="text"
                      class="form-control"
                      name="ip"
                      pattern="((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
                      required
                    />
                    <div class="invalid-feedback">Invalid input!</div>
                  </div>
                </div>
                <div class="form-group mb-3">
                  <div class="input-group">
                    <span class="input-group-text">Port</span
                    ><input
                      id="portInput"
                      type="number"
                      class="form-control"
                      name="port"
                      required
                      min="1"
                      max="65535"
                    />
                    <div class="invalid-feedback">Invalid input!</div>
                  </div>
                </div>
                <hr />
                <div class="row" style="margin-bottom: 12px;">
                  <div
                    class="col d-flex d-xxl-flex justify-content-center align-items-center"
                  >
                    <p style="margin-bottom: 0px;">Authentication:</p>
                  </div>
                  <div
                    class="col d-flex d-xxl-flex justify-content-center align-items-center justify-content-xxl-center align-items-xxl-center"
                  >
                    <input type="checkbox" id="proxyAuthCheck" name="auth" />
                  </div>
                </div>
                <div class="d-none" id="proxyAuthInputs">
                  <div class="form-group mb-3">
                    <div class="input-group">
                      <span class="input-group-text">Username</span
                      ><input
                        id="proxyUsernameInput"
                        type="text"
                        class="form-control auth-proxy"
                        name="username"
                      />
                      <div class="invalid-feedback">Invalid input!</div>
                    </div>
                  </div>
                  <div class="form-group mb-3">
                    <div class="input-group">
                      <span class="input-group-text">Password</span
                      ><input
                        id="proxyPasswordInput"
                        type="password"
                        class="form-control auth-proxy"
                        name="password"
                      />
                      <div class="invalid-feedback">Invalid input!</div>
                    </div>
                  </div>
                </div>
                <hr />
                <div class="form-group mb-3">
                  <label class="form-label">Activation:&nbsp;</label>
                  <div class="form-check disabled">
                    <input
                      type="radio"
                      class="form-check-input"
                      id="proxyRadioEnabled"
                      value="true"
                      name="proxyActivation"
                      required
                    /><label class="form-check-label" for="proxyRadioEnabled"
                      >Enabled</label
                    >
                  </div>
                  <div class="form-check">
                    <input
                      type="radio"
                      class="form-check-input"
                      id="proxyRadioDisabled"
                      value="false"
                      name="proxyActivation"
                      required
                      checked
                    /><label class="form-check-label" for="proxyRadioDisabled"
                      >Disabled</label
                    >
                    <div class="invalid-feedback">Invalid input!</div>
                  </div>
                </div>
                <div>
                  <button
                    class="btn btn-primary"
                    id="proxyFormBtn"
                    type="button"
                  >
                    Save</button
                  ><button
                    class="btn btn-secondary"
                    type="reset"
                    style="margin-left: 5px;"
                  >
                    Reset to default
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
        <div class="card shadow mb-4">
          <div
            class="card-header d-flex justify-content-between align-items-center"
          >
            <h6 class="text-danger fw-bold m-0">Danger Zone</h6>
          </div>
          <div class="card-body">
            <p class="text-center">
              Be careful with following settings!<br />Some changes cannot be
              undone!
            </p>
            <hr />
            <form id="igCredentialsForm" method="post" novalidate="">
              <p>Instagram credentials:</p>
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Username</span
                  ><input
                    id="usernameInput"
                    type="text"
                    class="form-control"
                    name="username"
                    required
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div class="form-group mb-3">
                <div class="input-group">
                  <span class="input-group-text">Password</span
                  ><input
                    id="passwordInput"
                    type="password"
                    class="form-control"
                    name="password"
                    required
                  />
                  <div class="invalid-feedback">Invalid input!</div>
                </div>
              </div>
              <div>
                <button
                  class="btn btn-primary"
                  id="igCredentialsFormBtn"
                  type="button"
                >
                  Save</button
                ><button
                  class="btn btn-secondary"
                  type="reset"
                  style="margin-left: 5px;"
                >
                  Reset to default
                </button>
              </div>
            </form>
            <hr />
            <div
              class="text-center"
              style="margin-right: 15px;margin-left: 15px;"
            >
              <div
                class="row d-xl-flex justify-content-center justify-content-xl-center align-items-xl-center"
                style="margin-bottom: 15px;"
              >
                <div class="col">
                  <button
                    class="btn btn-danger"
                    id="resetGlobalDefaultSettingsBtn"
                    type="button"
                  >
                    Reset default settings
                  </button>
                </div>
              </div>
              <div
                class="row d-xl-flex justify-content-center justify-content-xl-center align-items-xl-center"
                style="margin-bottom: 15px;"
              >
                <div class="col">
                  <button
                    class="btn btn-danger"
                    id="clearBotDataBtn"
                    type="button"
                  >
                    Clear cache &amp; data
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Vue } from "vue-class-component";

export default class SettingsPanel extends Vue {}
</script>

<style lang="scss"></style>
