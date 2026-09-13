import { createRouter, createWebHistory } from "vue-router";
import ParticipantsView from "../views/ParticipantsView.vue";

const routes = [
  {
    path: "/",
    name: "Dashboard",
    component: () => import("../views/ParticipantsView.vue"),
  },
  {
    path: "/participants",
    name: "Participants",
    component: () => import("../views/ParticipantsView.vue"),
  },
  {
    path: "/contributions",
    name: "Contributions",
    component: () => import("../views/ContributionsView.vue"),
  },
  {
    path: "/games",
    name: "Games",
    component: () => import("../views/GamesView.vue"),
  },
  {
    path: "/draws",
    name: "Draws",
    component: () => import("../views/DrawsView.vue"),
  },
  {
    path: "/tickets",
    name: "Tickets",
    component: () => import("../views/TicketsView.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
