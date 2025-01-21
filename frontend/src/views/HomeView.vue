<script setup lang="ts">
import TheResultsDisplay from '@/components/TheResultsDisplay.vue'
import TheSearchBar from '@/components/TheSearchBar.vue'
import { useAuthStore } from '@/stores/auth'
import type { ActorAutocompleteDetails } from '@/types'
import { ref, type Ref } from 'vue'
const auth = useAuthStore()
const actors: Ref<Array<ActorAutocompleteDetails>> = ref([])

async function submitRequest(srcActor: object, destActor: object) {
  console.log({ ...srcActor, ...destActor })
  await auth.postResource('/server/actors', {
    ...srcActor,
    ...destActor,
  })
}
async function getActors() {
  await auth.getResource('/server/actors')
  actors.value = auth.resourceData as Array<ActorAutocompleteDetails>
}
// getActors()
</script>
<template>
  <main>
    <TheSearchBar :actors="actors" @submit="submitRequest" />
    <TheResultsDisplay />
  </main>
</template>
