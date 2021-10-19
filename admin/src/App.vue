<template>
  <div>
    <div v-if="isLoggedIn() && isUserReady()">
      <Header/>
      <div class="container pt-5">
        <router-view/>
      </div>
      <Footer/>
    </div>
    <div v-else>
      <Login/>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import store from '@/store'
import cookies from '@/cookies'
import backend from './backend';

import User from "./models/user"

import Header from './components/Header.vue';
import Footer from './components/Footer.vue';
import Login from './components/Login.vue';

export default defineComponent({
  name: "App",
  components: {
    Header,
    Footer,
    Login
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
      return this.isLoggedIn() && store.state.user !== null
    },
    async fetchUser() {
      backend.GetSelf((user) => {
        store.commit("setUser", user)
      })
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

</style>
