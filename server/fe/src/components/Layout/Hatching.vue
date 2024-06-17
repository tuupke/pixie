<template>
  <linearGradient
      :id="id"
      :x2=dist
      gradientUnits="userSpaceOnUse"
      spreadMethod="repeat"
      gradientTransform="rotate(45)">
    <stop v-for="(offset, index) in offsets"
          :offset="offset"
          :stop-color="colors[Math.floor(index/2)]"/>
  </linearGradient>
</template>

<script setup lang="ts">

import {computed, withDefaults} from "vue";

const props = withDefaults(defineProps<{
  id: string
  width?: number
  colors: string[]
}>(), {
  width: 100,
})

// Dist must fill the width, note the gradient is rotated by 45 degrees
const dist = computed(() => (props.width/2/Math.sqrt(2)))

const offsets = computed(() => {
  const stops = [0]
  for (let i = 1; i < props.colors.length; i++) {
    const v = i / props.colors.length
    stops.push(v, v)
  }
  stops.push(1)

  return stops
})

</script>
