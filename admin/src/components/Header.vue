<template>
  <nav class="navbar" role="navigation">
    <div id="navbar-content" class="navbar-menu" style="height: 48px">
      <div class="navbar-brand">
        <a href="/" class="navbar-item">
          <b>osu!Mapping Tournaments</b>
        </a>

        <a
          role="button"
          class="navbar-burger"
          aria-label="menu"
          aria-expanded="false"
          data-target="navbar"
        >
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>

      <div id="navbar" class="navbar-menu">
        <!-- Navbar Main Menu with various navigation, including back to main page -->
        <div class="navbar-start">
          <router-link to="/" class="navbar-item"> Tournaments </router-link>
          <router-link to="/users" class="navbar navbar-item">
            Users
          </router-link>
        </div>

        <!-- Navbar End with Logout button -->
        <div class="navbar-end">
          <a href="/" class="navbar-item">Back to Homepage</a>
          <div class="navbar-divider"></div>
          <a href="/oauth/logout" class="navbar-item is-danger">Log Out</a>
          <!-- Round Username + Avatar To Show who is logged in -->
          <figure class="image is-48x48">
            <img
              :src="'https://a.ripple.moe/' + user?.ripple_id"
              alt="Your Avatar"
            />
          </figure>
        </div>
      </div>
    </div>
  </nav>
</template>

<script lang="ts" setup>
import { storeToRefs } from "pinia";
import { onMounted } from "vue";
import { useUserStore } from "../store/UserStore";

const UserStore = useUserStore();
const { user } = storeToRefs(UserStore);

onMounted(() => {
  console.log("Props User:", user);
});

</script>

<style scoped lang="scss">
@import "~bulma/bulma.sass";

.navbar-item {
  border-bottom: 2px solid transparent;
}

.router-link-active {
  border-bottom: 2px solid $primary;
}
</style>
