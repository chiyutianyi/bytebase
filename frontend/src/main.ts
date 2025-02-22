import axios from "axios";
import moment from "moment";
import { createApp } from "vue";
import App from "./App.vue";
import "./assets/css/inter.css";
import "./assets/css/tailwind.css";
import {
  BBAlert,
  BBAttention,
  BBAvatar,
  BBButtonAdd,
  BBButtonConfirm,
  BBCheckbox,
  BBContextMenu,
  BBModal,
  BBNotification,
  BBOutline,
  BBSelect,
  BBStepBar,
  BBStepTab,
  BBSwitch,
  BBTab,
  BBTabFilter,
  BBTable,
  BBTableCell,
  BBTableHeaderCell,
  BBTableSearch,
  BBTabPanel,
  BBTextField,
} from "./bbkit";
import dataSourceType from "./directives/data-source-type";
// @ts-ignore
import highlight from "./directives/highlight";
import { router } from "./router";
import { store } from "./store";
import {
  databaseSlug,
  dataSourceSlug,
  environmentName,
  environmentSlug,
  humanizeTs,
  instanceName,
  instanceSlug,
  isDev,
  isRelease,
  projectName,
  projectSlug,
  registerStoreWithActivityUtil,
  registerStoreWithRoleUtil,
  sizeToFit,
  urlfy,
} from "./utils";

registerStoreWithRoleUtil(store);
registerStoreWithActivityUtil(store);

console.debug("dev:", isDev());
console.debug("release:", isRelease());

axios.defaults.timeout = 10000;
axios.interceptors.request.use((request) => {
  if (isDev() && request.url!.startsWith("/api")) {
    console.debug(
      request.method?.toUpperCase() + " " + request.url + " request",
      JSON.stringify(request, null, 2)
    );
  }

  return request;
});

axios.interceptors.response.use(
  (response) => {
    if (isDev() && response.config.url!.startsWith("/api")) {
      console.debug(
        response.config.method?.toUpperCase() +
          " " +
          response.config.url +
          " response",
        JSON.stringify(response.data, null, 2)
      );
    }
    return response;
  },
  (error) => {
    if (error.response) {
      // When receiving 401 and is returned by our server, it means the current
      // login user's token becomes invalid. Thus we force a logout.
      // We could receive 401 when calling external service such as VCS provider,
      // in such case, we shouldn't logout.
      if (error.response.status == 401) {
        const host = store.getters["actuator/info"]().host;
        if (error.response.request.responseURL.startsWith(host))
          store.dispatch("auth/logout").then(() => {
            router.push({ name: "auth.signin" });
          });
      }

      if (error.response.data.message) {
        store.dispatch("notification/pushNotification", {
          module: "bytebase",
          style: "CRITICAL",
          title: error.response.data.message,
        });
      }
      return;
    } else if (error.code == "ECONNABORTED") {
      store.dispatch("notification/pushNotification", {
        module: "bytebase",
        style: "CRITICAL",
        title: "Connecting server timeout. Make sure the server is running.",
      });
      return;
    }

    throw error;
  }
);

// A global hook to collect errors to a central service
// app.config.errorHandler = function (err, vm, info) {
// };

// We need to restore the basic info in order to perform route authentication.
// Even using the <suspense>, it's still too late, thus we do the fetch here.
// We use finally because we always want to mount the app regardless of the error.
Promise.all([
  store.dispatch("actuator/fetchInfo"),
  store.dispatch("setting/fetchSetting"),
  store.dispatch("plan/fetchCurrentPlan"),
  store.dispatch("auth/restoreUser"),
]).finally(() => {
  const app = createApp(App);

  // Allow template to access various function
  app.config.globalProperties.window = window;
  app.config.globalProperties.console = console;
  app.config.globalProperties.moment = moment;
  app.config.globalProperties.humanizeTs = humanizeTs;
  app.config.globalProperties.isDev = isDev();
  app.config.globalProperties.isRelease = isRelease();
  app.config.globalProperties.sizeToFit = sizeToFit;
  app.config.globalProperties.urlfy = urlfy;
  app.config.globalProperties.environmentName = environmentName;
  app.config.globalProperties.environmentSlug = environmentSlug;
  app.config.globalProperties.projectName = projectName;
  app.config.globalProperties.projectSlug = projectSlug;
  app.config.globalProperties.instanceName = instanceName;
  app.config.globalProperties.instanceSlug = instanceSlug;
  app.config.globalProperties.databaseSlug = databaseSlug;
  app.config.globalProperties.dataSourceSlug = dataSourceSlug;

  app
    // Need to use a directive on the element.
    // The normal hljs.initHighlightingOnLoad() won't work because router change would cause vue
    // to re-render the page and remove the event listener required for
    .directive("highlight", highlight)
    .directive("data-source-type", dataSourceType)
    .use(store)
    .use(router)
    .component("BBAlert", BBAlert)
    .component("BBAttention", BBAttention)
    .component("BBAvatar", BBAvatar)
    .component("BBButtonAdd", BBButtonAdd)
    .component("BBButtonConfirm", BBButtonConfirm)
    .component("BBCheckbox", BBCheckbox)
    .component("BBContextMenu", BBContextMenu)
    .component("BBModal", BBModal)
    .component("BBNotification", BBNotification)
    .component("BBOutline", BBOutline)
    .component("BBSelect", BBSelect)
    .component("BBStepBar", BBStepBar)
    .component("BBStepTab", BBStepTab)
    .component("BBSwitch", BBSwitch)
    .component("BBTab", BBTab)
    .component("BBTabPanel", BBTabPanel)
    .component("BBTable", BBTable)
    .component("BBTableCell", BBTableCell)
    .component("BBTableHeaderCell", BBTableHeaderCell)
    .component("BBTableSearch", BBTableSearch)
    .component("BBTabFilter", BBTabFilter)
    .component("BBTextField", BBTextField)
    .mount("#app");
});
