import { defineStore } from "pinia";
import axios from "axios";
import router from "../router";
import { parseTimeFun } from "../assets/utilityFunctions/parseTime";
import { v4 as uuidv4 } from "uuid";
import { useDialogStore } from "./dialogStore";

export const useContentStore = defineStore("content", {
  state: () => ({
    dashboards: [],
    currentDashboard: {
      mode: null,
      id: null,
      index: null,
      name: null,
      type: null,
      content: [],
    },
  }),
  getters: {
    customizedDashboards() {
      return this.dashboards.filter((item) => item.type === "customized");
    },
    fixedDashboards() {
      return this.dashboards.filter((item) => item.type === "fixed");
    },
  },
  actions: {
    // Will get all of the dashboards that belong to the user. Will also reroute the user to the first dashboard in the list
    setDashboards() {
      axios
        .get("/api_server/manager/topic/all")
        .then((rs) => {
          this.dashboards = rs.data.data.filter((item) => item.id !== 63918);
          if (!this.currentDashboard.type || !this.currentDashboard.id) {
            this.currentDashboard.id = this.dashboards[0].id;
            this.currentDashboard.type = this.dashboards[0].type;
            router.replace({
              query: {
                id: this.dashboards[0].id,
                type: this.dashboards[0].type,
              },
            });
          }
          this.setCurrentDashboardContent();
        })
        .catch((e) => {
          console.log(e);
        });
    },
    // Called by the router. Will provide the store with query parameters 'id' and 'type'. If either isn't present, the setDashboards method will be called
    setRouteParams(mode, id, type) {
      this.currentDashboard.mode = mode;
      if (
        this.currentDashboard.id === id &&
        this.currentDashboard.type === type
      ) {
        return;
      }
      this.currentDashboard.id = id;
      this.currentDashboard.type = type;
      if (this.dashboards.length === 0) {
        this.setDashboards();
        return;
      }
      if (!id || !type) {
        this.setDashboards();
        return;
      }
      this.setCurrentDashboardContent();
    },
    // Will call an api to get the general info of a dashboard. This function will also filter out null components
    setCurrentDashboardContent() {
      axios
        .get(
          `/api_server/manager/topcomp/${this.currentDashboard.type}/${this.currentDashboard.id}`
        )
        .then((rs) => {
          this.currentDashboard.index = rs.data.index;
          this.currentDashboard.name = rs.data.name;
          this.currentDashboard.content = rs.data.components.filter(
            (comp) => comp.request_list !== null
          );
          this.setCurrentDashboardChartData();
        })
        .catch((e) => {
          if (e.response.status === 500) {
            const emptyDashboard = this.dashboards.filter(
              (item) =>
                item.id == this.currentDashboard.id &&
                item.type === this.currentDashboard.type
            );
            this.currentDashboard.index = emptyDashboard[0].index;
            this.currentDashboard.name = emptyDashboard[0].name;
            this.currentDashboard.content = [];
            return;
          }
          router.replace({ path: "/Dashboard" });
        });
    },
    // Will look through the components and call every requested api (POST) to get chart info for each component.
    // This function will save the chart configs of the components that don't have api paths
    setCurrentDashboardChartData() {
      for (let i = 0; i < this.currentDashboard.content.length; i++) {
        for (
          let j = 0;
          j < this.currentDashboard.content[i].request_list.length;
          j++
        ) {
          if (j === 0) {
            this.currentDashboard.content[i].chartData = [];
            this.currentDashboard.content[i].chartType = "";
          }
          this.currentDashboard.content[i].chartType +=
            this.currentDashboard.content[i].request_list[j].type;
          let payload = null;
          if (this.currentDashboard.content[i].request_list[j].form_data) {
            payload = this.parseRequestFormData(
              this.currentDashboard.content[i].request_list[j].form_data
            );
          }
          if (this.currentDashboard.content[i].request_list[j].path) {
            axios
              .post(
                `/api_server${this.currentDashboard.content[i].request_list[j].path}`,
                payload
              )
              .then((rs) => {
                this.currentDashboard.content[i].chartData.push(rs.data);
              })
              .catch((e) => {
                console.log(e);
              });
          } else {
            this.currentDashboard.content[i].chartData.push(
              this.currentDashboard.content[i].request_list[j].config
            );
          }
        }
      }
    },
    // This function will parse the form_data of each component and produce a payload that is to be attached to the POST request for chart info
    parseRequestFormData(form_data) {
      const timeData = parseTimeFun(form_data);
      const otherData = { ...form_data };
      if (otherData.time_step) delete otherData.time_step;
      if (otherData.time_unit) delete otherData.time_unit;
      if (otherData.last) delete otherData.last;
      const allData = {
        ...timeData,
        ...otherData,
      };

      const formData = new FormData();
      Object.keys(allData).forEach((key) => {
        formData.append(key, allData[key]);
      });
      return formData;
    },
    createNewDashboard(name) {
      const dialogStore = useDialogStore();
      axios
        .post(`/api_server/manager/topic/customized`, {
          index: uuidv4().slice(0, 8),
          name: name,
        })
        .then((rs) => {
          dialogStore.showNotification(
            "success",
            `成功將${name}加入自訂儀表板`
          );
          this.setDashboards();
        })
        .catch((e) => {
          console.log(e);
          dialogStore.showNotification("fail", e.response.status);
        });
    },
    changeCurrentDashboardName(name) {
      const dialogStore = useDialogStore();
      axios
        .put(
          `/api_server/manager/topic/customized/${this.currentDashboard.id}`,
          {
            name: name,
          }
        )
        .then((rs) => {
          dialogStore.showNotification(
            "success",
            `成功將儀表板名稱更改為${name}`
          );
          this.setDashboards();
        })
        .catch((e) => {
          console.log(e);
          dialogStore.showNotification("fail", e.response.status);
        });
    },
    deleteCurrentDashboard() {
      const dialogStore = useDialogStore();
      axios
        .delete(
          `/api_server/manager/topic/customized/${this.currentDashboard.id}`
        )
        .then((rs) => {
          dialogStore.showNotification(
            "success",
            `成功刪除${this.currentDashboard.name}`
          );
          this.setDashboards();
        })
        .catch((e) => {
          dialogStore.showNotification(
            "fail",
            `在刪除${this.currentDashboard.name}時出現失誤`
          );
        });
    },
    deleteComponent(topic_id, component_id, name) {
      const dialogStore = useDialogStore();
      if (!topic_id || !component_id || !name) {
        return;
      }
      axios
        .delete(
          `/api_server/manager/topcomp/customized/${this.currentDashboard.id}/topic/${topic_id}/component/${component_id}`
        )
        .then((rs) => {
          dialogStore.showNotification("success", `成功刪除${name}`);
          this.setDashboards();
        })
        .catch((e) => {
          dialogStore.showNotification("fail", `在刪除${name}時出現失誤`);
        });
    },
    addComponents(component_ids) {
      if (!component_ids) {
        return;
      }
      for (let i = 0; i < component_ids.length; i++) {
        axios
          .post(
            `/api_server/manager/topcomp/customized`,
            {
              admin_topic_id: 63918,
              component_id: +component_ids[i],
              topic_id: +this.currentDashboard.id,
            },
            {
              params: {
                id: 63918,
                type: "customized",
                componentid: +component_ids[i],
              },
            }
          )
          .then((rs) => {
            if (i === component_ids.length - 1) {
              this.setDashboards();
            }
          })
          .catch((e) => console.log(e));
      }
    },
  },
});
