<template>
  <g v-if="editing"
     @mousedown.stop="(e: MouseEvent) => { if (editing) {$emit('dragStart', {event: e, coord: coord});}}">
    <slot v-bind="$attrs"></slot>
  </g>
  <slot v-else v-bind="$attrs"></slot>
</template>

<script setup lang="ts">

import {inject} from "vue";
import {CoordinateInterface, RotationCoordinateInterface} from "../types.ts";

const emit = defineEmits(['dragStart'])
const editing = inject<boolean>("editing")! ?? false

const props = withDefaults(defineProps<{
  coord: CoordinateInterface & RotationCoordinateInterface,
  outline?: boolean
  rotate?: boolean
}>(), {
  outline: false,
  rotate: true,
})

</script>
