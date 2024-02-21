<template>
  <pattern :style="{fill: fill}" id="diagonalHatch" patternUnits="userSpaceOnUse" :width=4*dist :height=4*dist>
    <rect :width=4*dist :height=4*dist />
    <path :d=path :style="{stroke: color, strokeWidth: stroke}" />
  </pattern>
</template>

<script setup>

import {computed, defineProps} from "vue";

const props = defineProps({
  'color': {type: String, default: '#b7cbff'},
  'fill': {type: String, default: 'orange'},
  'stroke': {type: Number, default: 1},
  'dist': {type: Number, default: 1},
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
