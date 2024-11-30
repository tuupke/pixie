<template>
  <g :transform="'translate('+x+','+y+')'"
     @mouseover="(e) => {if (sequenceKey !== null) $emit('hoverOverElement', sequenceKey)}"
     @mouseout="(e) => {if (sequenceKey !== null) $emit('hoverOutElement', sequenceKey)}"
  >
    <CrossHairs :scale="scale" highlightable v-bind="$attrs" />
  </g>
</template>

<script setup lang="ts">

import {computed, inject} from "vue";
import CrossHairs from "./CrossHairs.vue";
import {CoordinateInterface, QualifiedKey} from "../../types.ts";

defineEmits<{
  hoverOverElement: [QualifiedKey]
  hoverOutElement: [QualifiedKey]
}>()

const props = withDefaults(defineProps<CoordinateInterface & {
  scale?: number,
  sequenceKey?: QualifiedKey,
  scaleMultiply?: number
}>(), {
  scaleMultiply: 1
});
const injectedScale = inject<number>("scale")!
const scale = computed<number>(() => (props.scale != undefined ? props.scale : (1/injectedScale.value)) * props.scaleMultiply)

</script>
