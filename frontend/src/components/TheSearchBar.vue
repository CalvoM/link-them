<script setup lang="ts">
import type { ActorAutocompleteDetails } from '@/types'
import {
  AutoComplete,
  type AutoCompleteCompleteEvent,
  type AutoCompleteOptionSelectEvent,
} from 'primevue'
import { ref, type Ref } from 'vue'
const srcActor: Ref<string> = ref('')
const destActor: Ref<string> = ref('')
const props = defineProps<{
  actors: Array<ActorAutocompleteDetails>
}>()
const emit = defineEmits<{
  submit: [srcActor: string, destActor: string]
}>()
const suggestedActors: Ref<Array<string>> = ref([])
const clickedSubmit: Ref<boolean> = ref(false)
const mappedActors: Ref<Array<string>> = ref([])
const mappingDone: Ref<boolean> = ref(false)
const mappedPics: Ref<Array<string>> = ref([])
const baseURL: string = 'https://image.tmdb.org/t/p/original/'
const srcActorPic: Ref<string> = ref('')
const destActorPic: Ref<string> = ref('')

function search() {
  clickedSubmit.value = true
  emit('submit', srcActor.value, destActor.value)
}
function getActors(event: AutoCompleteCompleteEvent) {
  if (event.query.length < 2) {
    return []
  }
  if (props.actors.length === 0 || mappingDone.value === false) {
    mappedActors.value = props.actors.map(
      (actorDetails: ActorAutocompleteDetails) => actorDetails.name,
    )
    mappedPics.value = props.actors.map(
      (actorDetails: ActorAutocompleteDetails) => actorDetails.profile_picture,
    )
    mappingDone.value = true
  }
  suggestedActors.value = mappedActors.value.filter((actor: string) =>
    actor.toLowerCase().includes(event.query.toLowerCase()),
  )
}
function optSelect(event: AutoCompleteOptionSelectEvent) {
  const i: number = mappedActors.value.indexOf(event.value)
  if (event.originalEvent.target?.id.startsWith('srcActor')) {
    srcActorPic.value = baseURL + mappedPics.value[i]
  }
  if (event.originalEvent.target?.id.startsWith('destActor')) {
    destActorPic.value = baseURL + mappedPics.value[i]
  }
}
</script>
<template>
  <div class="flex justify-evenly">
    <FloatLabel variant="on" class="flex-none">
      <AutoComplete
        forceSelection
        class="mb-2"
        inputId="srcActor"
        id="srcActor"
        type="text"
        v-model="srcActor"
        :suggestions="suggestedActors"
        @option-select="optSelect"
        @complete="getActors"
      />
      <label for="srcActor">Starting Actor</label>
    </FloatLabel>
    <Message
      icon="pi pi-times-circle"
      severity="error"
      v-if="!srcActor && clickedSubmit"
      id="msg1"
      >Actor is needed
    </Message>
    <FloatLabel variant="on" class="flex-none">
      <AutoComplete
        forceSelection
        inputId="destActor"
        type="text"
        id="destActor"
        v-model="destActor"
        :suggestions="suggestedActors"
        @option-select="optSelect"
        @complete="getActors"
      />
      <label for="destActor">Final Actor</label>
    </FloatLabel>
    <Message
      icon="pi pi-times-circle"
      severity="error"
      v-if="!destActor && clickedSubmit"
      id="msg2"
      >Actor is needed
    </Message>
    <Button
      label="Find link"
      aria-label="Search"
      icon="pi pi-search"
      severity="success"
      class="flex-none"
      @click="search"
    />
  </div>
  <div class="flex justify-evenly">
    <img :src="srcActorPic" alt="" class="rounded-full w-48 h-auto" />
    <img :src="destActorPic" alt="" class="rounded-full w-48 h-auto" />
    <span></span>
  </div>
</template>
