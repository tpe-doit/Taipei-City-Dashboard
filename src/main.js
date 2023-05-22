// Import the Essentials
import { createApp } from "vue";
import "./assets/styles/globalStyles.css";
import App from "./App.vue";

// Import Core Packages: Vue-Router, Pinia, ElementPlus
import router from "./router";
import { createPinia } from "pinia";

// Import Highcharts Related Packages
import Highcharts from "highcharts";
import HighchartsVue from "highcharts-vue";
import highchartsMore from "highcharts/highcharts-more";
import accessibilityInit from "highcharts/modules/accessibility";
import exportingInit from "highcharts/modules/exporting";
import exportData from "highcharts/modules/export-data";
import loadWordcloud from "highcharts/modules/wordcloud";
import loadTreemap from "highcharts/modules/treemap";

// Import Global Components
import PercentData from "./components/charts/PercentData.vue";
import DistrictData from "./components/charts/DistrictData.vue";
import AreaData from "./components/charts/AreaData.vue";
import TimeData from "./components/charts/TimeData.vue";
import GuageData from "./components/charts/GuageData.vue";

const app = createApp(App);

// Add Core Packages: Vue-Router, Pinia, HighCharts, ElementPlus
app.use(router);
const pinia = createPinia();
app.use(pinia);

// Add Highcharts
highchartsMore(Highcharts);
loadTreemap(Highcharts);
accessibilityInit(Highcharts);
exportingInit(Highcharts);
exportData(Highcharts);
loadWordcloud(Highcharts);
Highcharts.setOptions({
  lang: {
    // Highcharts Language Setting Documentation https://api.highcharts.com/highcharts/lang
    thousandsSep: ",",
    downloadCSV: "另存為CSV",
    downloadXLS: "另存為XLS",
    downloadJPEG: "另存為JPEG",
    downloadPNG: "另存為PNG",
    downloadSVG: "另存為SVG",
    downloadPDF: "另存為PDF",
    hideData: "隱藏資料",
    viewData: "顯示資料",
    viewFullscreen: "全螢幕顯示",
    printChart: "列印",
    exportData: {
      categoryHeader: "",
    },
  },
  colors: ["#2879C7", "#09ABB7", "#71B07E", "#B5C06A", "#7781B6"],
});
app.use(HighchartsVue, { Highcharts });

app.component("PercentData", PercentData);
app.component("DistrictData", DistrictData);
app.component("AreaData", AreaData);
app.component("TimeData", TimeData);
app.component("GuageData", GuageData);

app.mount("#app");
