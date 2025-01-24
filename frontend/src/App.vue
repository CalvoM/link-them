<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header>
    <div class="wrapper p-4 bg-green-900 mb-4 text-white flex">
      <nav>
        <RouterLink to="/" class="mx-4 font-extrabold">Link-Them</RouterLink>
        <RouterLink to="#">About</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
  <footer
    class="p-6 mx-auto font-extrabold bg-green-900 text-white text-center bottom-0 w-full"
  >
    All Credit to
    <a href="https://www.themoviedb.org/" class="underline decoration-double"
      >TMDB</a
    >
  </footer>
</template>

<style scoped></style>
