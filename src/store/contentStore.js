import { defineStore } from "pinia";
import axios from "axios";
import router from "../router";
import { v4 as uuidv4 } from "uuid";
import { useDialogStore } from "./dialogStore";

export const useContentStore = defineStore("content", {
  state: () => ({
    dashboards: [],
    components: {},
    currentDashboard: {
      mode: null,
      index: null,
      name: null,
      components: {},
      content: [],
    },
    mapLayers: [],
  }),
  getters: {},
  actions: {
    // Called by the router. Will provide the store with query parameters 'id' and 'type'. If either isn't present, the setDashboards method will be called
    setRouteParams(mode, index) {
      this.currentDashboard.mode = mode;
      if (this.currentDashboard.index === index) {
        return;
      }
      this.currentDashboard.index = index;
      if (this.dashboards.length === 0) {
        this.setDashboards();
        return;
      }
      if (!index) {
        this.setDashboards();
        return;
      }
      this.setCurrentDashboardContent();
    },
    // Will get all of the dashboards that belong to the user. Will also reroute the user to the first dashboard in the list
    setDashboards() {
      axios
        .get("/dashboards/all_dashboards.json")
        .then((rs) => {
          this.dashboards = rs.data.data;
          if (!this.currentDashboard.index) {
            this.currentDashboard.index = this.dashboards[0].index;
            router.replace({
              query: {
                index: this.currentDashboard.index,
              },
            });
          }
          this.setComponents();
        })
        .catch((e) => {
          console.log(e);
        });
    },
    setComponents() {
      axios
        .get("/dashboards/all_components.json")
        .then((rs) => {
          this.components = rs.data.data;
          this.setMapLayers();
          this.setCurrentDashboardContent();
        })
        .catch((e) => console.log(e));
    },
    // Will call an api to get the general info of a dashboard. This function will also filter out null components
    setCurrentDashboardContent() {
      const currentDashboardInfo = this.dashboards.find(
        (item) => item.index === this.currentDashboard.index
      );
      this.currentDashboard.name = currentDashboardInfo.name;
      this.currentDashboard.components = currentDashboardInfo.components;
      this.currentDashboard.icon = currentDashboardInfo.icon;
      this.currentDashboard.content = [];
      this.currentDashboard.components.forEach((component) => {
        this.currentDashboard.content.push(this.components[component]);
      });
      if (this.currentDashboard.index === "map-layers") {
        return;
      }
      this.setCurrentDashboardChartData();
    },
    setMapLayers() {
      const mapLayerInfo = this.dashboards.find(
        (item) => item.index === "map-layers"
      );
      mapLayerInfo.components.forEach((component) => {
        this.mapLayers.push(this.components[component]);
      });
    },
    // Will look through the components and call every requested api (POST) to get chart info for each component.
    // This function will save the chart configs of the components that don't have api paths
    setCurrentDashboardChartData() {
      this.currentDashboard.components.forEach((component, index) => {
        axios
          .get(`/chartData/${component}.json`)
          .then((rs) => {
            this.currentDashboard.content[index].chart_data = rs.data.data;
          })
          .catch((e) => {
            console.log(e);
          });
        if (this.currentDashboard.content[index].history_data) {
          axios
            .get(`/historyData/${component}.json`)
            .then((rs) => {
              this.currentDashboard.content[index].history_data = rs.data.data;
            })
            .catch((e) => {
              console.log(e);
            });
        }
      });
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
