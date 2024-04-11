<template>
  <line :x1="start.x" :y1="start.y" :x2="end.x" :y2="end.y" stroke="white" stroke-width="15" stroke-linecap="round" />
  <line :x1="start.x" :y1="start.y" :x2="end.x" :y2="end.y" stroke="gray" stroke-width="5" stroke-linecap="round" />
  <g v-if="editing" v-for="coord in [start, end]" :transform="translate(coord)"
     @mousedown.stop="(e: MouseEvent) => $emit('dragStart', {event: e, coord: coord})" >
    <Crosshairs :scale=scale />
  </g>
</template>

<script lang="ts" setup>

import {CoordinateInterface, PathInterface} from "../../types.ts";
import Crosshairs from "./Crosshairs.vue";

withDefaults(defineProps<PathInterface & {editing?: boolean, scale?: number}>(), {
  editing: false,
  scale: 1
});

defineEmits(['dragStart'])

function translate(c: CoordinateInterface): string {
  return 'translate('+c.x+','+c.y+')';
}

</script>
