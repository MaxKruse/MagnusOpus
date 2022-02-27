<template>
  <div class="tournament">
    <table v-if="tournaments.length > 0" class="table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Description</th>
          <th>Start Time</th>
          <th>End Time</th>
          <th>Registration Start</th>
          <th>Registration End</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="tournament in tournaments"
          :key="tournament.id"
          class="tr"
          :class="getTimeClass(tournament)"
        >
          <td>{{ tournament.name }}</td>
          <td>{{ tournament.description }}</td>
          <td>{{ IsoDateToLocalStr(tournament.start_time) }}</td>
          <td>{{ IsoDateToLocalStr(tournament.end_time) }}</td>
          <td>{{ IsoDateToLocalStr(tournament.registration_start_time) }}</td>
          <td>{{ IsoDateToLocalStr(tournament.registration_end_time) }}</td>
        </tr>
      </tbody>
    </table>
    <span v-else class="tag is-warning is-large">No tournaments found.</span>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";

import Tournament from "../../models/tournament";
import backend from "../../backend";

import { IsoDateToLocalStr } from "../../utils/iso_to_local";

const tournaments = ref<Tournament[]>([]);

onMounted(async () => {
  tournaments.value = await backend.GetTournaments();
})

function getTimeClass(tournament: Tournament) {
  const now = new Date();
  const start = new Date(tournament.start_time);
  const end = new Date(tournament.end_time);
  if (now < start) {
    return "is-warning";
  } else if (now < end) {
    return "is-success";
  } else {
    return "is-danger";
  }
}
</script>

<style lang="scss" scoped>
// center .tournament div
.tournament {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
