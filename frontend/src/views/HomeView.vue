<script setup lang="ts">
import TheSearchBar from '@/components/TheSearchBar.vue'
import { useAuthStore } from '@/stores/auth'
import type { ActorAutocompleteDetails } from '@/types'
import { ref, type Ref } from 'vue'
const auth = useAuthStore()
const actors: Ref<Array<ActorAutocompleteDetails>> = ref([])
// actors.value = new Array<string>('One', 'Two', 'Three')
function submitRequest(srcActor: string, destActor: string) {
  console.log(srcActor)
  console.log(destActor)
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
