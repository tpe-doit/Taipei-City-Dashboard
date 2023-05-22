import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      "/geo_server": {
        target: "https://geoserver.tuic.gov.taipei/geoserver/",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/geo_server/, ""),
      },
      "/api_server": {
        //測試機: 172.25.201.105 - 重啟大概要15分鐘
        // target: 'http://172.25.201.105:8090/api/v1/',

        //正式機: 172.25.201.116
        target: "http://172.25.201.105:8090/api/v1",

        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api_server/, ""),
      },
    },
  },
});
