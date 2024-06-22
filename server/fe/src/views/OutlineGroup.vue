<template>
  <g ref="dragGroup" :transform="transform"
     @mousedown.stop="(e: MouseEvent) => $emit('dragStart', e)">
    <slot v-bind="$attrs"></slot>
    <path
        v-if="outlinePath"
        :d="outlinePath"
        stroke="black"
        stroke-width="5"
        stroke-linecap="round"
        stroke-linejoin="round"
        fill="none" />
  </g>
</template>

<script setup lang="ts">

import {computed, onMounted, onUnmounted, reactive, ref} from "vue";

const emit = defineEmits(['dragStart'])

withDefaults(defineProps<{
  transform?: string,
}>(), {
  transform: "",
})

const outlinePath = computed(() => `M ${size.x+10} ${size.y+10}
l ${size.width-20} 0
l 0 ${size.height-20}
l -${size.width-20} 0
l 0 -${size.height-20}
`)

const dragGroup = ref()
const size = reactive<{
  width: number
  height: number
  x: number
  y: number
}>({height: 0, width: 0, x: 0, y: 0});

const updater = ref<number>(0)
onMounted(() => {
  updater.value = window.setInterval(rescale, 20)
})

onUnmounted(() => {
  window.clearInterval(updater.value)
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
}
</script>
