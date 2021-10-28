<template>
  <div class="tournament">
    <table class="table" v-if="tournaments.length > 0">
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
        <tr class="tr" v-for="tournament in tournaments" :key="tournament.id" :class="getTimeClass(tournament)">
          <td>{{ tournament.name }}</td>
          <td>{{ tournament.description }}</td>
          <td>{{ utils.IsoDateToLocalStr(tournament.start_time) }}</td>
          <td>{{ utils.IsoDateToLocalStr(tournament.end_time) }}</td>
          <td>{{ utils.IsoDateToLocalStr(tournament.registration_start_time) }}</td>
          <td>{{ utils.IsoDateToLocalStr(tournament.registration_end_time) }}</td>
        </tr>
      </tbody>
    </table>
    <span class="tag is-warning is-large" v-else>No tournaments found.</span>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import CustomError from '@/models/CustomError'
import Tournament from '@/models/tournament'
import backend from '@/backend'

import { IsoDateToLocalStr } from '@/utils/iso_to_local'

export default defineComponent({
  name: 'TournamentOverview',
  data() {
    return {
      utils: {
        IsoDateToLocalStr
      },
      tournaments: [] as Tournament[]
    }
  },
  mounted() {
    backend.GetTournaments((tournaments: Tournament[], err: CustomError | null) => {
      if (err !== null) {
        console.error(err)
        return
      }
      this.tournaments = tournaments
    })
  },
  methods: {
    getTimeClass(tournament: Tournament) {
      const now = new Date()
      const start = new Date(tournament.start_time)
      const end = new Date(tournament.end_time)
      if (now < start) {
        return 'is-warning'
      } else if (now < end) {
        return 'is-success'
      } else {
        return 'is-danger'
      }
    }
  }
})

</script>

<style lang="scss" scoped>
</style>