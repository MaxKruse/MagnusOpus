<template>
  <nav class="navbar" role="navigation">
    <div id="navbar-content" class="navbar-menu">
      <div class="navbar-brand">
        <a href="/" class="navbar-item">
          <b>osu!Mapping Tournaments</b>
        </a>

        <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" data-target="navbar">
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>

      <div class="navbar-menu">
        <!-- Navbar Main Menu with various navigation, including back to main page -->
        <div id="navbar" class="navbar-start">
          <router-link to="/" class="navbar-item">
            Tournaments
          </router-link>
          <router-link to="/users" class="navbar-item">
            Users
          </router-link>
        </div>

          <!-- Navbar End with Logout button -->
        <div class="navbar-end">
          <div class="navbar-item">
            <div class="buttons">
              <a href="/" class="navbar-item button">Back to Homepage</a>
              <div class="navbar-divider"></div>
              <a href="/oauth/logout" class="button is-danger">Log Out</a>
            </div>
          </div>
        </div>
      </div>



    </div>
  </nav>
</template>

<script lang="ts">

import { defineComponent } from 'vue'
export default defineComponent({
  name: 'Header',
  props: {
    activeElement: {
      type: String,
      default: '',
    },
  },
  methods: {
    checkActiveElement(element: string) {
      return this.activeElement === element ? 'is-active' : ''
    },
  }
})
</script>

<style scoped lang="scss">

</style>