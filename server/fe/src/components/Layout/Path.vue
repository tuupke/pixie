<template>
  <line :x1="start.x" :y1="start.y" :x2="end.x" :y2="end.y" stroke="white" stroke-width="15" stroke-linecap="round" />
  <line :x1="start.x" :y1="start.y" :x2="end.x" :y2="end.y" stroke="gray" stroke-width="5" stroke-linecap="round" />
  <PlacedCrossHairs
      v-if="editing"
      v-for="coord in [start, end]"
      v-bind="coord"
      :highlightable="true"
      :scale-multiply="scaleMultiply"
      :color="color" />
</template>

<script lang="ts" setup>

import {CoordinateInterface} from "../../types.ts";
import PlacedCrossHairs from "./PlacedCrossHairs.vue";

withDefaults(defineProps<{
  start: CoordinateInterface,
  end: CoordinateInterface,
  editing?: boolean
  scale?: number
  scaleMultiply?: number
  color?: string
}>(), {
  editing: false,
  scale: 1,
  scaleMultiply: 0.7,
  color: "green"
});

defineEmits(['dragStart'])

</script>
