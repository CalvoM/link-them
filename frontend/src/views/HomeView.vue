<script setup lang="ts">
import TheSearchBar from '@/components/TheSearchBar.vue'
import { useAuthStore } from '@/stores/auth'
import type { ActorAutocompleteDetails } from '@/types'
import { ref, type Ref } from 'vue'
const auth = useAuthStore()
const actors: Ref<Array<ActorAutocompleteDetails>> = ref([])

async function submitRequest(srcActor: string, destActor: string) {
  const srcActorDetails = actors.value
    .filter((actor: ActorAutocompleteDetails) => actor.name === srcActor)
    .map((actor: ActorAutocompleteDetails) => ({
      srcActor: actor.name,
      srcActorID: actor.id,
    }))[0]
  const destActorDetails = actors.value
    .filter((actor: ActorAutocompleteDetails) => actor.name === destActor)
    .map((actor: ActorAutocompleteDetails) => ({
      destActor: actor.name,
      destActorID: actor.id,
    }))[0]
  console.log(srcActorDetails)
  console.log(destActorDetails)
  await auth.postResource('/server/actors', {
    ...srcActorDetails,
    ...destActorDetails,
  })
}
async function getActors() {
  await auth.getResource('/server/actors')
  actors.value = auth.resourceData as Array<ActorAutocompleteDetails>
}
getActors()
</script>
<template>
  <main>
    <TheSearchBar :actors="actors" @submit="submitRequest" />
  </main>
</template>
