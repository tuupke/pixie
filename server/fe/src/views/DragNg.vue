<template>
  <g :transform=transform ref="group">
    <rect v-bind="size" :class="{outlineRect: true, drawOutline: true}"/>
    <g class="dragNg" ref="dragGroup"
       @mousedown.stop="(e: MouseEvent) => $emit('dragStart', {event: e, coord: coord})"
       @mouseover.stop="(e: MouseEvent) => $emit('maybeShow', {event: e, coord: coord})"
       @click.stop="(e: MouseEvent) => $emit('maybeSelect', {event: e, coord: coord})">
      <slot></slot>
    </g>
    <rect v-if="coord.rotation != undefined && rotate" class="rotate" :x=size.width+size.x-50 :y=size.y-50 width="100"
          height="100"
          @mousedown.stop=rotateStart
    />
  </g>
</template>

<style scoped>
.dragNg:hover {
  cursor: move;
}

.outlineRect, .drawOutline {
  fill: none;
}

.drawOutline {
  border: 1px solid blue !important;
}

.rotate:hover {
  cursor: grab;
  fill: orange;
}

</style>

<script setup lang="ts">

import {computed, onMounted, reactive, ref} from "vue";
import {CoordinateInterface, RotationCoordinateInterface} from "../types.ts";

const emit = defineEmits(['dragStart', 'maybeShow', 'maybeSelect', 'rotateStart'])

const props = withDefaults(defineProps<{
  coord: CoordinateInterface & RotationCoordinateInterface,
  outline?: boolean
  rotate?: boolean
  transform?: boolean
}>(), {
  outline: false,
  rotate: true,
  transform: false,
})

const dragGroup = ref()
const size = reactive<{
  width: number
  height: number
  x: number
  y: number
}>({height: 0, width: 0, x: 0, y: 0});

onMounted(rescale)

const transform = computed(() => props.transform
    ? 'translate(' + props.coord.x + ',' + props.coord.y + ')' + (
    props.coord.rotation !== undefined
        ? 'rotate(' + props.coord.rotation + ')'
        : '')
    : '');

const group = ref()

function rotateStart(e: MouseEvent) {
  const box = dragGroup.value.getBBox()

  emit('rotateStart', {
    event: e,
    coord: props.coord,
    width: box.width,
    height: box.height,
    x: box.x,
    y: box.y,
  })
}

function rescale() {
  const box = dragGroup.value.getBBox()
  size.width = box.width
  size.height = box.height
  size.x = box.x
  size.y = box.y
}
</script>
