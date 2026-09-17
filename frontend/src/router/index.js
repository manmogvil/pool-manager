import { createRouter, createWebHistory } from "vue-router";
import { auth } from "../services/api";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("../views/LoginView.vue"),
    meta: { guest: true },
  },
  {
    path: "/",
    name: "Dashboard",
    component: () => import("../views/DashboardView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/users",
    name: "Users",
    component: () => import("../views/UsersView.vue"),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: "/contributions",
    name: "Contributions",
    component: () => import("../views/ContributionsView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/games",
    name: "Games",
    component: () => import("../views/GamesView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/draws",
    name: "Draws",
    component: () => import("../views/DrawsView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/tickets",
    name: "Tickets",
    component: () => import("../views/TicketsView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/",
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const isAuthenticated = auth.isAuthenticated();
  const user = auth.getUser();

  if (to.meta.requiresAuth && !isAuthenticated) {
    next("/login");
  } else if (to.meta.guest && isAuthenticated) {
    next("/");
  } else if (to.meta.requiresAdmin && user?.role !== "admin") {
    next("/");
  } else {
    next();
  }
});

export default router;
