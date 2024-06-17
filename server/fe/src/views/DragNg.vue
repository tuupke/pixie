<template>
  <g v-if="transform" :transform=coordToTransform>
    <g class="dragNg"
       ref="dragGroup"
       @mousedown="(e: MouseEvent) => $emit('dragStart', {event: e, coord: coord})"
    >
      <slot></slot>
    </g>
  </g>
  <g ref="dragGroup" v-else>
    <slot
          :attrs="$attrs"
          @mousedown="(e: MouseEvent) => $emit('dragStart', {event: e, coord: coord})"
    />
  </g>
</template>

<style scoped>

.outlineRect, .drawOutline {
  fill: none;
}

.drawOutline {
  border: 1px solid blue !important;
}

</style>

<script setup lang="ts">

import {computed, inject, onMounted, onUnmounted, reactive, ref, watch} from "vue";
import {CoordinateInterface, RotationCoordinateInterface} from "../types.ts";
import PlacedCrossHairs from "../components/Layout/PlacedCrossHairs.vue";

const emit = defineEmits(['dragStart', 'rotateStart'])
const scale = inject<number>('scale')!

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

const dragGroup = ref();

const size = reactive<{
  width: number
  height: number
  x: number
  y: number
}>({height: 0, width: 0, x: 0, y: 0});

const coordToTransform = computed(() => props.transform
    ? 'translate(' + props.coord.x + ',' + props.coord.y + ')' + (
    props.coord.rotation !== undefined
        ? 'rotate(' + props.coord.rotation + ')'
        : '')
    : '');


function rotateStart(e: MouseEvent) {
  const box = dragGroup.value.getBBox()
  console.log("start")

  emit('rotateStart', {
    event: e,
    coord: props.coord,
    width: box.width,
    height: box.height,
    x: box.x,
    y: box.y,
  })
}

onMounted(rescale)
onMounted(() => {
  try {
    const el =dragGroup.value.querySelector("#dragEl")
    console.log(el)
    el.addEventListener('mouseover', () => {
      console.log('over')
    })
  } catch (e) {}
})

onUnmounted(() => {
  try {
    const el =dragGroup.value.querySelector("#dragEl")
    console.log(el)
    el.querySelector("#dragEl").removeEventListener('mouseover', rotateStart)
  } catch (e) {}
})

watch(() => props.transform, rescale)

function rescale() {
  if (!props.transform) {
    return
  }

  const box = dragGroup.value.getBBox()
  size.width = box.width
  size.height = box.height
  size.x = box.x
  size.y = box.y
}
</script>
