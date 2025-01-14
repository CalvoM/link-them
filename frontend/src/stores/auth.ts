import { defineStore } from 'pinia'
import { NonTokenHTTP, TokenHTTP } from '@/helpers'
import type { Ref } from 'vue'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const resourceData: Ref<Object> = ref({})

  async function getResource(url: string) {
    const data = await TokenHTTP.get(url)
    resourceData.value = await data.data
    return
  }
  async function postResource(url: string, resource: any) {
    const data = await TokenHTTP.post(url, resource)
    return data
  }
  function syncPostResource(url: string, resource: any) {
    const data = TokenHTTP.post(url, resource)
    return data
  }
  return {
    resourceData,
    getResource,
    postResource,
    syncPostResource,
  }
})
