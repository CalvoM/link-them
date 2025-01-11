<script setup lang="ts">
import TheSearchBar from '@/components/TheSearchBar.vue'
import { ref, type Ref } from 'vue'
const actors: Ref<Array<string>> = ref([])
actors.value = new Array<string>('One', 'Two', 'Three')
function submitRequest(srcActor: string, destActor: string) {
  console.log(srcActor)
  console.log(destActor)
}
</script>
<template>
  <main>
    <TheSearchBar :actors="actors" @submit="submitRequest" />
  </main>
</template>
