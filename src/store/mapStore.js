import { createApp, defineComponent, nextTick } from "vue";
import MapPopup from "../components/map/MapPopup.vue";

import axios from "axios";
import { defineStore, createPinia } from "pinia";
import { useContentStore } from "./contentStore";

import mapboxGl from "mapbox-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import mapStyle from "../assets/configs/mapbox/mapStyle.js";
import { MapPositions } from "../assets/configs/mapbox/mapType.js";
import {
  MapObjectConfig,
  TaipeiTown,
  TaipeiVillage,
  TaipeiBuilding,
  maplayerCommonPaint,
  maplayerCommonLayout,
} from "../assets/configs/mapbox/mapConfig.js";

export const useMapStore = defineStore("map", {
  state: () => ({
    currentLayers: [],
    mapConfigs: {},
    currentVisibleLayers: [],
    map: null,
    popup: null,
  }),
  getters: {},
  actions: {
    addToMapLayerList(map_config) {
      map_config.forEach((element) => {
        let mapLayerId = `${element.index}-${element.type}`;
        if (this.currentLayers.find((element) => element === mapLayerId)) {
          this.turnOnMapLayerVisibility(mapLayerId);
          if (
            !this.currentVisibleLayers.find((element) => element === mapLayerId)
          ) {
            this.currentVisibleLayers.push(mapLayerId);
          }
          return;
        }
        let appendLayerId = { ...element };
        appendLayerId.layerId = mapLayerId;
        this.fetchLocalGeoJson(appendLayerId);
      });
    },

    /* Functions that fetch map layers */

    // Raster Maps were used in the official city dashboard due to file size concerns.
    // The open source version doesn't feature any as all data is static and access to our raster maps are restricted to within the Taipei city gov's intranet
    fetchRaster() {},

    // Fetch map layer data for geojsons stored locally
    fetchLocalGeoJson(map_config) {
      axios
        .get(`/mapData/${map_config.index}.geojson`)
        .then((rs) => {
          this.addMapLayerSource(map_config, rs.data);
        })
        .catch((e) => console.log(e));
    },

    // The open source version doesn't feature any as all data is static and access to our remote geoJSONs are restricted to within the Taipei city gov's intranet
    fetchRemoteGeoJson(mapLayerInfo) {},

    /* Functions that directly mutate the map */
    initializeMapBox() {
      const MAPBOXTOKEN = import.meta.env.VITE_MAPBOXTOKEN;
      mapboxGl.accessToken = MAPBOXTOKEN;
      this.map = new mapboxGl.Map({
        ...MapObjectConfig,
        style: mapStyle,
      });
      this.map.addControl(new mapboxGl.NavigationControl());
      this.map.doubleClickZoom.disable();
      this.map
        .on("style.load", () => {
          this.initializeBasicLayers();
          // this.initialize3DLayers();
        })
        .on("click", (event) => {
          if (this.popup) {
            this.popup = null;
          }
          this.addPopup(event);
        });
    },
    // Called when the mapbox instance is first initialized and adds two basic layers to the map
    initializeBasicLayers() {
      fetch(`/mapData/taipei_town.geojson`)
        .then((response) => response.json())
        .then((data) => {
          this.map
            .addSource("taipei_town", { type: "geojson", data: data })
            .addLayer(TaipeiTown("dark"));
        });
      fetch(`/mapData/taipei_village.geojson`)
        .then((response) => response.json())
        .then((data) => {
          this.map
            .addSource("taipei_village", { type: "geojson", data: data })
            .addLayer(TaipeiVillage("dark"));
        });
      this.addSymbolSources();
    },
    // Called when the mapbox instance is first initialized and adds 3d renderings of taipei buildings to the map
    initialize3DLayers() {},

    // Called when the mapbox instance is first initialized and adds the icons that will be used in the map
    addSymbolSources() {
      const images = ["metro", "triangle_green"];
      images.forEach((element) => {
        this.map.loadImage(`/images/${element}.png`, (error, image) => {
          if (error) throw error;
          this.map.addImage(element, image);
        });
      });
    },

    addMapLayerSource(map_config, data) {
      this.map.addSource(`${map_config.layerId}-source`, {
        type: "geojson",
        data: data,
      });
      this.addMapLayer(map_config);
    },

    addMapLayer(map_config) {
      let extra_paint_configs = {};
      let extra_layout_configs = {};
      if (map_config.icon) {
        extra_paint_configs = {
          ...maplayerCommonPaint[`${map_config.type}-${map_config.icon}`],
        };
        extra_layout_configs = {
          ...maplayerCommonLayout[`${map_config.type}-${map_config.icon}`],
        };
      }
      if (map_config.size) {
        extra_paint_configs = {
          ...extra_paint_configs,
          ...maplayerCommonPaint[`${map_config.type}-${map_config.size}`],
        };
        extra_layout_configs = {
          ...extra_layout_configs,
          ...maplayerCommonLayout[`${map_config.type}-${map_config.size}`],
        };
      }
      this.map.addLayer({
        id: map_config.layerId,
        type: map_config.type,
        paint: {
          ...maplayerCommonPaint[`${map_config.type}`],
          ...extra_paint_configs,
          ...map_config.paint,
        },
        layout: {
          ...maplayerCommonLayout[`${map_config.type}`],
          ...extra_layout_configs,
        },
        source: `${map_config.layerId}-source`,
      });
      this.currentLayers.push(map_config.layerId);
      this.mapConfigs[map_config.layerId] = map_config;
      this.currentVisibleLayers.push(map_config.layerId);
    },

    addPopup(event) {
      const clickFeatureDatas = this.map.queryRenderedFeatures(event.point, {
        layers: this.currentVisibleLayers,
      });
      if (!clickFeatureDatas || clickFeatureDatas.length === 0) {
        return;
      }
      const mapConfigs = this.mapConfigs;
      this.popup = new mapboxGl.Popup()
        .setLngLat(event.lngLat)
        .setHTML('<div id="vue-popup-content"></div>')
        .addTo(this.map);
      const PopupComponent = defineComponent({
        extends: MapPopup,
        setup() {
          return {
            popupContent: clickFeatureDatas[0],
            mapConfig: mapConfigs[clickFeatureDatas[0].layer.id],
          };
        },
      });
      nextTick(() => {
        const app = createApp(PopupComponent);
        app.mount("#vue-popup-content");
      });
    },
    removePopup() {
      if (this.popup) {
        this.popup.remove();
      }
      this.popup = null;
    },

    turnOnMapLayerVisibility(mapLayerId) {
      this.map.setLayoutProperty(mapLayerId, "visibility", "visible");
    },

    // Called by ComponentMapCharts to turn the visibility off for a layer but not removing it from the store
    turnOffMapLayerVisibility(map_config) {
      map_config.forEach((element) => {
        let mapLayerId = `${element.index}-${element.type}`;
        if (this.map.getLayer(mapLayerId)) {
          this.map.setLayoutProperty(mapLayerId, "visibility", "none");
        }
        this.currentVisibleLayers = this.currentVisibleLayers.filter(
          (element) => element !== mapLayerId
        );
      });
      this.removePopup();
    },
    // Called when the user is switching between maps
    clearOnlyLayers() {
      this.currentLayers.forEach((element) => {
        this.map.removeLayer(element);
        this.map.removeSource(`${element}-source`);
      });
      this.currentLayers = [];
      this.mapConfigs = {};
      this.currentVisibleLayers = [];
      this.removePopup();
    },
    // Called when user navigates away from the map
    clearEntireMap() {
      this.currentLayers = [];
      this.mapConfigs = {};
      this.map = null;
      this.currentVisibleLayers = [];
      this.removePopup();
    },
  },
});
