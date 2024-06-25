<template>
  <g
     @scrollElement="(e: RotateEvent) => $emit('scrollElement', e)"
     @mousedown.stop="(e: MouseEvent) => $emit('dragStart', {event: e, coord: coord})">
    <slot v-bind="$attrs"></slot>
  </g>
</template>

<script setup lang="ts">

import {CoordinateInterface, DragStartEvent, RotateEvent, RotationCoordinateInterface} from "../types.ts";

defineEmits<{
  dragStart: [DragStartEvent]
  scrollElement: [RotateEvent]
}>()

withDefaults(defineProps<{
  coord: CoordinateInterface | RotationCoordinateInterface,
  outline?: boolean
  rotate?: boolean
}>(), {
  outline: false,
  rotate: true,
})

</script>
