<template>
  <div class="app">
    <div v-if="isLoggedIn()">
      <div v-if="isUserReady()">
        <p>Welcome {{ getUsername() }}! You are logged in!</p>
        <router-view/>
      </div>
    </div>
    <div v-else>
      <LoginVue/>
    </div>

  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import store from '@/store'
import cookies from '@/cookies'
import backend from './backend';

import User from "./models/user"

import LoginVue from './components/Login.vue';

export default defineComponent({
  name: "App",
  components: {
    LoginVue
  },
  data() {
    return {
      store: store
    }
  },
  computed: {
    user(): User {
      return store.state.user
    }
  },
  methods: {
    isLoggedIn: () => cookies.sessionToken() !== "",
    isUserReady() {
      return this.isLoggedIn() && store.state.user?.username !== "ERROR NAME"
    },
    async fetchUser() {
      console.log("Fetching user...")
      backend.GetSelf((user, raw) => {
        let u = user as User
        console.log("Got user:", u)
        console.log("Raw response:", raw)
        store.commit("setUser", u)
      })
    },
    getUsername() {
      return this.user?.ripple_id
    },
  },
  async mounted() {
    if (this.isLoggedIn()) {
      await this.fetchUser();
    }
  }
})
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

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
