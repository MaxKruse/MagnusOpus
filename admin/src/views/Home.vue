<template>
  <div class="home">
    <!-- Do cool stuff here -->
    <HomeComponent/>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import HomeComponent from '@/components/Home.vue' // @ is an alias to /src

export default defineComponent({
  name: 'Home',
  components: {
    HomeComponent
  }
})
</script>
