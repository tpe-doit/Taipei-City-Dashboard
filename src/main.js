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
import BarChart from "./components/charts/BarChart.vue";
import TreemapChart from "./components/charts/TreemapChart.vue";
import ColumnChart from "./components/charts/ColumnChart.vue";
import BarPercentChart from "./components/charts/BarPercentChart.vue";
import GuageChart from "./components/charts/GuageChart.vue";
import RadarChart from "./components/charts/RadarChart.vue";
import TimelineSeparateChart from "./components/charts/TimelineSeparateChart.vue";
import TimelineStackedChart from "./components/charts/TimelineStackedChart.vue";
import MapLegend from "./components/charts/MapLegend.vue";
import MetroChart from "./components/charts/MetroChart.vue";

const app = createApp(App);

// Add Core Packages: Vue-Router, Pinia, Apexcharts
app.use(router);
const pinia = createPinia();
app.use(pinia);
app.use(VueApexCharts);

app.component("DistrictChart", DistrictChart);
app.component("DonutChart", DonutChart);
app.component("BarChart", BarChart);
app.component("TreemapChart", TreemapChart);
app.component("ColumnChart", ColumnChart);
app.component("BarPercentChart", BarPercentChart);
app.component("GuageChart", GuageChart);
app.component("RadarChart", RadarChart);
app.component("TimelineSeparateChart", TimelineSeparateChart);
app.component("TimelineStackedChart", TimelineStackedChart);
app.component("MapLegend", MapLegend);
app.component("MetroChart", MetroChart);

app.mount("#app");
