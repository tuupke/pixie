<template>
  <g v-if="editing"
     @scrollElement="(e: RotateEvent) => $emit('scrollElement', e)"
     @mousedown.stop="(e: MouseEvent) => { if (editing) {$emit('dragStart', {event: e, coord: coord});}}">
    <slot v-bind="$attrs"></slot>
  </g>
  <slot v-else v-bind="$attrs"></slot>
</template>

<script setup lang="ts">

import {inject} from "vue";
import {CoordinateInterface, RotateEvent, RotationCoordinateInterface} from "../types.ts";

const emit = defineEmits(['dragStart', 'scrollElement'])
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
