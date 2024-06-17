<template>
  <g :transform="'translate('+x+','+y+')'">
    <CrossHairs :scale="scale" highlightable v-bind="$attrs" />
  </g>
</template>

<script setup lang="ts">

import {computed, inject} from "vue";
import CrossHairs from "./CrossHairs.vue";
import {CoordinateInterface} from "../../types.ts";

const props = withDefaults(defineProps<CoordinateInterface & {scale?: number, scaleMultiply?: number}>(), {
  scaleMultiply: 1
});
const injectedScale = inject<number>("scale")!
const scale = computed<number>(() => (props.scale != undefined ? props.scale : (1/injectedScale.value)) * props.scaleMultiply)

</script>
