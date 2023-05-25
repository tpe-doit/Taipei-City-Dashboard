// Import the Essentials
import { createApp } from "vue";
import "./assets/styles/globalStyles.css";
import App from "./App.vue";

// Import Core Packages: Vue-Router, Pinia, Apexcharts
import router from "./router";
import { createPinia } from "pinia";
import VueApexCharts from "vue3-apexcharts";

// Import Global Components
import DistrictChart from "./components/charts/DistrictChart.vue";
import DonutChart from "./components/charts/DonutChart.vue";
import HorizontalBarChart from "./components/charts/HorizontalBarChart.vue";
import TreemapChart from "./components/charts/TreemapChart.vue";

const app = createApp(App);

// Add Core Packages: Vue-Router, Pinia, Apexcharts
app.use(router);
const pinia = createPinia();
app.use(pinia);
app.use(VueApexCharts);

app.component("DistrictChart", DistrictChart);
app.component("DonutChart", DonutChart);
app.component("HorizontalBarChart", HorizontalBarChart);
app.component("TreemapChart", TreemapChart);

app.mount("#app");
