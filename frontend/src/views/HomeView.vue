<script setup lang="ts">
import TheResultsDisplay from '@/components/TheResultsDisplay.vue'
import TheSearchBar from '@/components/TheSearchBar.vue'
import { useAuthStore } from '@/stores/auth'
import type { ActorAutocompleteDetails, MovieResponseDetails } from '@/types'
import { ref, type Ref } from 'vue'
const auth = useAuthStore()
const movies: Ref<Array<MovieResponseDetails>> = ref([])

async function submitRequest(srcActor: object, destActor: object) {
  const moviesResponse = await auth.postResource('/server/actors', {
    ...srcActor,
    ...destActor,
  })
  movies.value = moviesResponse.data
  console.log(movies.value)
}
</script>
<template>
  <main>
    <TheSearchBar @submit="submitRequest" />
    <TheResultsDisplay :movie-details="movies" />
  </main>
</template>
