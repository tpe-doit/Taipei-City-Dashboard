import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/Login.vue";
import Dashboard from "../views/Dashboard.vue";
import Map from "../views/Map.vue";
import { useAuthStore } from "../store/authStore";
import { useContentStore } from "../store/contentStore";
import { useMapStore } from "../store/mapStore";

// Auth: 0 - Not Logged In, 1 - Logged In User, 2 - Logged In Admin

const routes = [
  {
    path: "/",
    name: "homeRedirect",
    redirect: "/login",
    meta: {
      auth: 0,
    },
  },
  {
    path: "/login",
    name: "login",
    component: Login,
    meta: {
      auth: 0,
    },
  },
  {
    path: "/Dashboard",
    name: "dashboard",
    component: Dashboard,
    meta: {
      auth: 1,
    },
  },
  {
    path: "/MapView",
    name: "mapView",
    component: Map,
    meta: {
      auth: 1,
    },
  },
  {
    path: "/:pathMatch(.*)*",
    name: "notFoundRedirect",
    redirect: "/Dashboard",
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  base: import.meta.env.BASE_URL,
  routes,
});

// If the route has an auth level higher than 1 and no valid tokens are present, redirect to login
router.beforeEach((to) => {
  const store = useAuthStore();
  if (to.matched.some((record) => record.meta.auth === 1)) {
    if (store.tokens.access_token && store.tokens.refresh_token) {
      return;
    }
    return { name: "login" };
  }
});

// If the route has an auth level of 0 and there are valid tokens, redirect to dashboard
router.beforeEach((to) => {
  const store = useAuthStore();
  if (to.matched.some((record) => record.meta.auth === 0)) {
    if (store.tokens.access_token && store.tokens.refresh_token) {
      return { name: "dashboard" };
    }
    return;
  }
});

// Pass in route info to contentStore if the path starts with /Dashboard or /Mapview
router.beforeEach((to) => {
  const contentStore = useContentStore();
  const mapStore = useMapStore();
  if (to.path === "/Dashboard" || to.path === "/Mapview") {
    contentStore.setRouteParams(to.path, to.query.id, to.query.type);
    mapStore.clearMapStore();
  }
});

export default router;
