<template>
  <div>
    <div v-if="isLoggedIn()">
      <Header />
      <div class="container pt-5 pb-5">
        <router-view />
        <loading
          :active="!(isLoggedIn() && isUserReady())"
          :is-full-page="true"
        />
      </div>
      <Footer />
    </div>
    <div v-else>
      <Login />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useUserStore } from "./store/UserStore";

import cookies from "./cookies";
import backend from "./backend";

import Header from "./components/Header.vue";
import Footer from "./components/Footer.vue";
import Login from "./components/Login.vue";

import User from "./models/user";

import "vue-loading-overlay/dist/vue-loading.css";

import { onMounted, ref } from "vue";

const userStore = useUserStore();

function isLoggedIn() {
  return cookies.sessionToken() !== "";
}

function isUserReady() {
  return isLoggedIn() && userStore.user !== null;
}

onMounted(async () => {
  userStore.setUser(await backend.GetSelf())
});
</script>

<style lang="scss">
@charset "utf-8";
@import "~bulma/bulma.sass";

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

.is-success {
  background-color: rgba($success, 0.5);
}

.is-warning {
  background-color: rgba($warning, 0.5);
}

.is-danger {
  background-color: rgba($danger, 0.5);
}
</style>
