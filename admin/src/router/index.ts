import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import HomeView from "../views/Home.vue";
import TournamentOverview from "../components/tournament/TournamentOverview.vue"

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Tournament Overview",
    component: TournamentOverview
  }
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
