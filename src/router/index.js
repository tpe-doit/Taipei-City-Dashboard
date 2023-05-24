import { createRouter, createWebHistory } from "vue-router";
import Dashboard from "../views/Dashboard.vue";
import Map from "../views/Map.vue";
import { useContentStore } from "../store/contentStore";
import { useMapStore } from "../store/mapStore";

// Auth: 0 - Not Logged In, 1 - Logged In User, 2 - Logged In Admin

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

// Pass in route info to contentStore if the path starts with /dashboard or /mapview
router.beforeEach((to) => {
  const contentStore = useContentStore();
  const mapStore = useMapStore();
  if (to.path === "/dashboard" || to.path === "/mapview") {
    contentStore.setRouteParams(to.path, to.query.index);
  }
  if (to.path === "/dashboard") {
    mapStore.clearEntireMap();
  } else if (to.path === "/mapview") {
    mapStore.clearOnlyLayers();
  }
});

export default router;
