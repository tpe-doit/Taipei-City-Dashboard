import axios from "axios";
import { defineStore } from "pinia";
import { useContentStore } from "./contentStore";
import {
  TaipeiTown,
  TaipeiVillage,
  TaipeiBuilding,
} from "../assets/configs/mapbox/mapConfig.js";

export const useMapStore = defineStore("map", {
  state: () => ({
    currentLayers: [],
  }),
  getters: {},
  actions: {
    // Called by ComponentMapCharts to add map config data to the store and trigger apis to fetch map data
    // If a layer already exists, this function will only toggle visibility
    // mapLayerId = maptype-index
    addToMapLayerList(mapConfig) {
      const mapType = Object.keys(mapConfig)[0];
      const mapLayerId = `${mapType}-${mapConfig[mapType].index}`;
      const existingItemIndex = this.currentLayers.findIndex(
        (item) => item.mapLayerId === mapLayerId
      );
      if (existingItemIndex >= 0) {
        this.currentLayers[existingItemIndex].visible = true;
        return;
      }
      const mapLayerInfo = {
        mapType: mapType,
        mapLayerId: mapLayerId,
        visible: true,
        ...mapConfig[mapType],
      };
      this.currentLayers.push(mapLayerInfo);
      if (mapType === "raster") {
        this.fetchRaster(mapLayerInfo);
      } else if (mapType === "raster" || mapLayerInfo.url) {
        this.fetchRemoteGeoJson(mapLayerInfo);
      } else {
        this.fetchLocalGeoJson(mapLayerInfo);
      }
    },
    // Called by ComponentMapCharts to turn the visibility off for a layer but not removing it from the store
    turnOffMapLayerVisibility(mapConfig) {
      const mapType = Object.keys(mapConfig)[0];
      const mapLayerId = `${mapType}-${mapConfig[mapType].index}`;
      const existingItemIndex = this.currentLayers.findIndex(
        (item) => item.mapLayerId === mapLayerId
      );
      if (existingItemIndex >= 0) {
        this.currentLayers[existingItemIndex].visible = false;
      }
    },
    // Fetch map layer data for raster maps
    fetchRaster(mapLayerInfo) {
      console.log("raster");
    },
    // Fetch map layer data for geojsons stored locally
    fetchLocalGeoJson(mapLayerInfo) {
      axios
        .get(`/mapData/${mapLayerInfo.index}.geojson`)
        .then((rs) => {
          console.log("local success");
        })
        .catch((e) => console.log(e));
    },
    // Fetch map layer data for remote geoJsons
    fetchRemoteGeoJson(mapLayerInfo) {
      const contentStore = useContentStore();
      const payload = contentStore.parseRequestFormData(mapLayerInfo.form_data);
      axios
        .post(`api_server${mapLayerInfo.url}`, payload)
        .then((rs) => {
          console.log("remote success");
        })
        .catch((e) => console.log(e));
    },

    /* Functions that directly mutate the map */

    // Called when the mapbox instance is first initialized and adds two basic layers to the map
    initializeBasicLayers(map) {
      const layers = map.getStyle().layers;
      fetch(`/mapData/taipei_town.geojson`)
        .then((response) => response.json())
        .then((data) => {
          map
            .addSource("taipei_town", { type: "geojson", data: data })
            .addLayer(TaipeiTown("dark"));
        });
      fetch(`/mapData/taipei_village.geojson`)
        .then((response) => response.json())
        .then((data) => {
          map
            .addSource("taipei_village", { type: "geojson", data: data })
            .addLayer(TaipeiVillage("dark"));
        });
    },
    // Called when the mapbox instance is first initialized and adds 3d renderings of taipei buildings to the map
    initialize3DLayers(map) {
      map
        .addSource("TaipeiBuild", {
          type: "vector",
          scheme: "tms",
          tiles: [
            `${location.origin}/geo_server/gwc/service/tms/1.0.0/taipei_vioc:tp_building_height@EPSG:900913@pbf/{z}/{x}/{y}.pbf`,
          ],
        })
        .addLayer(TaipeiBuilding, "airport-label");
    },

    // Clear the map store
    clearMapStore() {
      this.currentLayers.splice(0, this.currentLayers.length);
    },
  },
});
