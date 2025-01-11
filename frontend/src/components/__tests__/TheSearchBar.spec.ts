import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TheSearchBar from '../TheSearchBar.vue'
describe('TheSearchBar', () => {
  const actor_suggestions: Array<string> = new Array<string>(
    'Alex Wamatu',
    'Abel Mutua',
    'DJ Shiti',
  )
  it('emits the submit event', () => {
    const wrapper = mount(TheSearchBar, { props: { actors: actor_suggestions } })
    const actorOneAutocomplete = wrapper.find('input#srcActor')
    const actorTwoAutocomplete = wrapper.find('input#destActor')
    actorOneAutocomplete.setValue(actor_suggestions[0])
    actorTwoAutocomplete.setValue(actor_suggestions[1])
    const submitButton = wrapper.find('button')
    submitButton.trigger('click')
    const submitEvent = wrapper.emitted()
    expect(submitEvent).toHaveProperty('submit')
  })
  it('Displays messge when submitted without values', () => {
    const wrapper = mount(TheSearchBar, { props: { actors: actor_suggestions } })
    const submitButton = wrapper.find('button')
    // const msg1 = wrapper.get("div#msg1")
    // const msg2 = wrapper.get("div#msg2")
    expect(wrapper.find('#msg1').isVisible()).toBe(false)
    expect(wrapper.find('#msg2').isVisible()).toBe(false)
    submitButton.trigger('click')
    expect(wrapper.find('#msg1').isVisible()).toBe(true)
    expect(wrapper.find('#msg2').isVisible()).toBe(true)
  })
})
