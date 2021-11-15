<template>
  <footer class="myfooter">
    <div class="has-content has-text-centered">
      <p>
        <b>osu!Mapping Tournaments</b> by <a href="https://github.com/maxkruse/">Maximilian Kruse</a>
      </p>
      <p>
        <b>Ripple</b> by <a href="http://ripple.moe">Ripple Team</a>
      </p>
    </div>
  </footer>
</template>

<script lang="ts">

import { defineComponent } from 'vue'
export default defineComponent({
  name: 'Footer',
})
</script>

<style scoped lang="scss">
.myfooter {
  bottom: 0;
  width: 100%;
  position: sticky;
  padding: 3rem;
  background-color: #fafafa;
  border-top: 1px solid #eaeaea;
}

</style>