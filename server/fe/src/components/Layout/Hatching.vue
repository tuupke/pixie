<template>
  <pattern :style="{fill: colorB}" id="diagonalHatch" patternUnits="userSpaceOnUse" :width=4*dist :height=4*dist>
    <rect :width=4*dist :height=4*dist />
    <path :d=path :style="{stroke: colorA, strokeWidth: stroke}" />
  </pattern>
</template>

<script setup lang="ts">

import {computed, defineProps, withDefaults} from "vue";

const props = withDefaults(defineProps<{
  colorA: string
  colorB: string
  stroke: number
  dist: number
}>(), {
  colorA: '#b7cbff',
  colorB: 'orange',
  stroke: 1,
  dist: 1,
})

const path = computed(() => {
  const one = props.dist
  const two = 2*props.dist
  const four = 4*props.dist
  const three = 3*props.dist
  const five = 5*props.dist

  return 'M-'+one+','+one+' l'+two+',-'+two+
         ' M0,'+four+' l'+four+',-'+four+
         ' M'+three+','+five+' l'+two+',-'+two
})

</script>
