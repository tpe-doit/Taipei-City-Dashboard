// Cleaned

import { createRouter, createWebHistory } from "vue-router";
import { useContentStore } from "../store/contentStore";
import { useMapStore } from "../store/mapStore";
import Dashboard from "../views/Dashboard.vue";
import Map from "../views/Map.vue";

const routes = [
  {
    path: "/",
    redirect: "/dashboard",
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: Dashboard,
  },
  {
    path: "/mapview",
    name: "mapview",
    component: Map,
  },
  {
    path: "/:pathMatch(.*)*",
    name: "notFoundRedirect",
    redirect: "/dashboard",
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  base: import.meta.env.BASE_URL,
  routes,
});

router.beforeEach((to) => {
  const contentStore = useContentStore();
  const mapStore = useMapStore();
  // Pass in route info to contentStore if the path starts with /dashboard or /mapview
  if (to.path === "/dashboard" || to.path === "/mapview") {
    contentStore.setRouteParams(to.path, to.query.index);
  }
  // Clear the entire mapStore if the path doesn't start with /mapview
  if (to.path !== "/mapview") {
    mapStore.clearEntireMap();
  }
  // Clear only map layers if the path starts with /mapview
  else if (to.path === "/mapview") {
    mapStore.clearOnlyLayers();
  }
});

export default router;
