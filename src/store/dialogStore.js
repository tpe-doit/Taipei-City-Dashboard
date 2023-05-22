import { defineStore } from "pinia";

export const useDialogStore = defineStore("dialog", {
  state: () => ({
    dialogs: {
      addDashboard: false,
      notificationBar: false,
      addComponent: false,
      dashboardSettings: false,
      moreInfo: false,
      userSettings: false,
    },
    notification: {
      status: "",
      message: "",
    },
    moreInfoContent: null,
  }),
  getters: {},
  actions: {
    showDialog(dialog) {
      this.dialogs[dialog] = true;
    },
    hideAllDialogs() {
      const keys = Object.keys(this.dialogs);
      for (let i = 0; i < keys.length; i++) {
        this.dialogs[keys[i]] = false;
      }
      this.moreInfoContent = null;
    },
    showNotification(status, message) {
      this.showDialog("notificationBar");
      this.notification = {
        status: status,
        message: message,
      };
      setTimeout(() => {
        this.dialogs.notificationBar = false;
      }, 3000);
    },
    showMoreInfo(content) {
      this.dialogs.moreInfo = true;
      this.moreInfoContent = content;
    },
  },
});
