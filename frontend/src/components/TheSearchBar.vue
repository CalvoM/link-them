<script setup lang="ts">
import { AutoComplete, type AutoCompleteCompleteEvent } from 'primevue';
import { ref, type Ref } from 'vue'
const srcActor: Ref<string> = ref('')
const destActor: Ref<string> = ref('')
const props = defineProps<{
  actors: Array<string>
}>();
const emit = defineEmits<{
  submit: [srcActor: string, destActor: string]
}>();
const suggestedActors: Ref<Array<string>> = ref([])
const clickedSubmit: Ref<boolean> = ref(false)

function search() {
  clickedSubmit.value = true
  emit('submit', srcActor.value, destActor.value)
}
function getActors(event: AutoCompleteCompleteEvent) {
  suggestedActors.value = props.actors.filter((actor: string) => actor.toLowerCase().includes(event.query.toLowerCase()))
}
</script>
<template>
  <div class="flex justify-evenly">
    <FloatLabel variant="on" class="flex-none">
      <AutoComplete class="mb-2" inputId="srcActor" id="srcActor" type="text" v-model="srcActor"
        :suggestions="suggestedActors" @complete="getActors" />
      <label for="srcActor">Starting Actor</label>
    </FloatLabel>
    <Message icon="pi pi-times-circle" severity="error" v-if="!srcActor && clickedSubmit" id="msg1">Actor is
      needed</Message>
    <FloatLabel variant="on" class="flex-none">
      <AutoComplete inputId="destActor" type="text" id="destActor" v-model="destActor" :suggestions="suggestedActors"
        @complete="getActors" />
      <label for="destActor">Final Actor</label>
    </FloatLabel>
    <Message icon="pi pi-times-circle" severity="error" v-if="!destActor && clickedSubmit" id="msg2">Actor is
      needed</Message>
    <Button label="Find link" aria-label="Search" icon="pi pi-search" severity="success" class="flex-none"
      @click="search" />
  </div>
</template>
