<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import type { ActorAutocompleteDetails } from '@/types'
import {
  AutoComplete,
  type AutoCompleteCompleteEvent,
  type AutoCompleteOptionSelectEvent,
} from 'primevue'
import { ref, type Ref } from 'vue'
const srcActor: Ref<string> = ref('')
const destActor: Ref<string> = ref('')
const emit = defineEmits<{
  submit: [srcActor: object, destActor: object]
}>()
const suggestedActors: Ref<Array<string>> = ref([])
const clickedSubmit: Ref<boolean> = ref(false)
const mappedActors: Ref<Array<string>> = ref([])
const mappedPics: Ref<Array<string>> = ref([])
const baseURL: string = 'https://image.tmdb.org/t/p/original/'
const srcActorPic: Ref<string> = ref('')
const destActorPic: Ref<string> = ref('')
const auth = useAuthStore()
const actors: Ref<Array<ActorAutocompleteDetails>> = ref()
const srcActorObj: Ref<object> = ref()
const destActorObj: Ref<object> = ref()

function search() {
  clickedSubmit.value = true
  emit('submit', srcActorObj.value, destActorObj.value)
}
async function getActors(event: AutoCompleteCompleteEvent) {
  if (event.query.length < 2) {
    return []
  }
  await auth.getResource(`/server/actors?query=${event.query}`)
  actors.value = auth.resourceData as Array<ActorAutocompleteDetails>
  mappedActors.value = actors.value.map(
    (actorDetails: ActorAutocompleteDetails) => actorDetails.name,
  )
  mappedPics.value = actors.value.map(
    (actorDetails: ActorAutocompleteDetails) => actorDetails.profile_picture,
  )
  suggestedActors.value = mappedActors.value.filter((actor: string) =>
    actor.toLowerCase().includes(event.query.toLowerCase()),
  )
}
function optSelect(event: AutoCompleteOptionSelectEvent) {
  if (event.value < 0) return
  const i: number = mappedActors.value.indexOf(event.value)
  if (event.originalEvent.target?.id.startsWith('srcActor')) {
    srcActorPic.value = baseURL + mappedPics.value[i]
    srcActorObj.value = actors.value
      .filter(
        (actorDetails: ActorAutocompleteDetails) =>
          actorDetails.name == srcActor.value,
      )
      .map((actor: ActorAutocompleteDetails) => ({
        srcActor: actor.name,
        srcActorID: actor.id,
      }))[0]
  }
  if (event.originalEvent.target?.id.startsWith('destActor')) {
    destActorPic.value = baseURL + mappedPics.value[i]
    destActorObj.value = actors.value
      .filter(
        (actorDetails: ActorAutocompleteDetails) =>
          actorDetails.name == destActor.value,
      )
      .map((actor: ActorAutocompleteDetails) => ({
        destActor: actor.name,
        destActorID: actor.id,
      }))[0]
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
      label="Find Common Movie"
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
