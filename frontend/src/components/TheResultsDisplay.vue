<script setup lang="ts">
import type { MovieResponseDetails } from '@/types'
const baseURL: string = 'https://image.tmdb.org/t/p/original'

const props = defineProps<{
  movieDetails: Array<MovieResponseDetails>
}>()
</script>
<template>
  <div class="flex flex-row flex-wrap">
    <img
      :src="baseURL + movie.poster_picture"
      :alt="movie.title"
      v-for="(movie, index) in props.movieDetails"
      :key="index"
      class="rounded-lg w-72 h-auto inline-block p-4 mx-auto"
    />
  </div>
</template>
