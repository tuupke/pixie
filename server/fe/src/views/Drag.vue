<template>
  <g
     ref="dragGroup" non-existing-prop="12"
     @mousedown.stop="(e: MouseEvent) => { if (editing) {$emit('dragStart', {event: e, coord: coord});}}">
    <slot v-bind="$attrs"></slot>
  </g>
</template>

<script setup lang="ts">

import {computed, inject, onMounted, reactive, ref, watch} from "vue";
import {CoordinateInterface, RotationCoordinateInterface} from "../types.ts";

const emit = defineEmits(['dragStart', 'rotateStart'])
const scale = inject<number>('scale')! ?? 1
const editing = inject<boolean>("editing")! ?? false

const props = withDefaults(defineProps<{
  coord: CoordinateInterface & RotationCoordinateInterface,
  outline?: boolean
  rotate?: boolean
}>(), {
  outline: false,
  rotate: true,
})

const outlinePath = computed(() => `M ${size.x} ${size.y}
l ${size.width} 0
l 0 ${size.height}
l -${size.width} 0
l 0 -${size.height}
`)

const dragGroup = ref()
const size = reactive<{
  width: number
  height: number
  x: number
  y: number
}>({height: 0, width: 0, x: 0, y: 0});

watch(props.coord, rescale)

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

onMounted(rescale)
onMounted(() => {
  window.setTimeout(() => {
    if (!props.rotate) {
      return
    }
    try {
      const els = dragGroup.value.querySelectorAll("#dragHandle")
      if (els.length !== 1) {
        return
      }
      els[0].addEventListener('mouseover', (ee) => {
        console.log('over', ee)
      })
    } catch (e) {
      console.log(e)
    }
  }, 100)
})

function rescale() {
  if (!dragGroup.value) {
    return
  }

  const box = dragGroup.value.getBBox()
  size.width = box.width
  size.height = box.height
  size.x = box.x
  size.y = box.y

  // console.log(toRaw(size))
}
</script>
