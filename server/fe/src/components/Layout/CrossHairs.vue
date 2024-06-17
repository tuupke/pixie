<template>
  <g id="outer" :transform="'scale('+actualScale+')'" @mouseover="hover=true" @mouseout="hover=false">
  <g :transform="'translate(0, '+-11.5+') rotate(45)'">
    <svg xmlns="http://www.w3.org/2000/svg">
      <circle class="stroked" fill="white" :stroke="color" :cx="8" :cy="8" :r="6"/>
      <path class="stroked" fill="none" :stroke="color" d="M 8 0 L 8 6.5"/>
      <path class="stroked" fill="none" :stroke="color" d="M 0 8 L 6.5 8"/>
      <path class="stroked" fill="none" :stroke="color" d="M 8 9.5 L 8 16"/>
      <path class="stroked" fill="none" :stroke="color" d="M 9.5 8 L 16 8"/>
    </svg>
  </g>
  </g>
</template>
<style scoped>
#outer:hover {
  cursor: v-bind('props.mouse');
}

</style>
<script setup lang="ts">

import {computed, ref} from "vue";

const props = withDefaults(defineProps<{
  scale?: number
  color?: string
  highlighted?: boolean
  mouse?: string,
  highlightable?: boolean
}>(), {
  scale: 1,
  color: 'blue',
  highlighted: false,
  mouse: 'default',
  highlightable: false
})

const hover = ref<boolean>(false)
const actualScale = computed<number>(() => (props.scale * (props.highlightable && props.highlighted ? 1.2: 1) * (props.highlightable && hover.value ?1.2: 1)))

</script>
