// Cleaned

/* authStore */
/*
The contentStore calls APIs to get content info and stores it.
*/

import router from "../router";
import { defineStore } from "pinia";
import axios from "axios";

export const useContentStore = defineStore("content", {
  state: () => ({
    // Stores all dashboards data. Reference the structure in /public/dashboards/all_dashboards.json
    dashboards: [],
    // Stores all components data. Reference the structure in /public/dashboards/all_components.json
    components: {},
    // Picks out the components that are map layers and stores them here
    mapLayers: [],
    // Stores information of the current dashboard
    currentDashboard: {
      // /mapview or /dashboard
      mode: null,
      index: null,
      name: null,
      content: [],
    },
  }),
  getters: {},
  actions: {
    /* Steps in adding content to the application */

    // 1. Check the current path and execute actions based on the current path
    setRouteParams(mode, index) {
      this.currentDashboard.mode = mode;
      // 1-1. Don't do anything if the path is the same
      if (this.currentDashboard.index === index) {
        return;
      }
      this.currentDashboard.index = index;
      // 1-2. If there is no dashboards info, call the setDashboards method (2.)
      if (this.dashboards.length === 0) {
        this.setDashboards();
        return;
      }
      // 1-3. If there is dashboard info but no index is defined, call the setDashboards method (2.)
      if (!index) {
        this.setDashboards();
        return;
      }
      // 1-4. If all info is present, skip steps 2, 3, 4 and call the setCurrentDashboardContent method (5.)
      this.setCurrentDashboardContent();
    },
    // 2. Call an API to get all dashboard info and reroute the user to the first dashboard in the list
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
          // After getting dashboard info, call the setComponents (3.) method to get component info
          this.setComponents();
        })
        .catch((e) => {
          console.log(e);
        });
    },
    // 3. Call and API to get all components info
    setComponents() {
      axios
        .get("/dashboards/all_components.json")
        .then((rs) => {
          this.components = rs.data.data;
          // Step 4.
          this.setMapLayers();
          // Step 5.
          this.setCurrentDashboardContent();
        })
        .catch((e) => console.log(e));
    },
    // 4. Adds components that are map layers into a separate store to be used in mapview
    setMapLayers() {
      const mapLayerInfo = this.dashboards.find(
        (item) => item.index === "map-layers"
      );
      mapLayerInfo.components.forEach((component) => {
        this.mapLayers.push(this.components[component]);
      });
    },
    // 5. Finds the info for the current dashboard based on the index and adds it to "currentDashboard"
    setCurrentDashboardContent() {
      const currentDashboardInfo = this.dashboards.find(
        (item) => item.index === this.currentDashboard.index
      );
      this.currentDashboard.name = currentDashboardInfo.name;
      this.currentDashboard.icon = currentDashboardInfo.icon;
      this.currentDashboard.content = currentDashboardInfo.components.map(
        (item) => {
          return this.components[item];
        }
      );
      // no need to call additional chart data APIs for the map layers dashboard
      if (this.currentDashboard.index === "map-layers") {
        return;
      }
      this.setCurrentDashboardChartData();
    },
    // 6. Call an API for each component to get its chart data and store it
    // Will call an additional API if the component has history data
    setCurrentDashboardChartData() {
      this.currentDashboard.content.forEach((component, index) => {
        axios
          .get(`/chartData/${component.id}.json`)
          .then((rs) => {
            this.currentDashboard.content[index].chart_data = rs.data.data;
          })
          .catch((e) => {
            console.log(e);
          });
        if (this.currentDashboard.content[index].history_data) {
          axios
            .get(`/historyData/${component.id}.json`)
            .then((rs) => {
              this.currentDashboard.content[index].history_data = rs.data.data;
            })
            .catch((e) => {
              console.log(e);
            });
        }
      });
    },

    /* Inactive Functions due to lack of backend */

    // Call this function to create a new dashboard. Pass in the new dashboard name.
    createNewDashboard(name) {},
    // Call this function to change the dashboard name. Pass in the new dashboard name.
    changeCurrentDashboardName(name) {},
    // Call this function to delete the current active dashboard.
    deleteCurrentDashboard() {},
    // Call this function to delete a component. Pass in related info.
    deleteComponent(topic_id, component_id, name) {},
    // Call this function to add components to the current dashboard. Pass in an array of component ids.
    addComponents(component_ids) {},
  },
});
