<template>
  <g :transform="'translate(' + (x) + ',' + (y) +') rotate('+rotation+') scale('+1/translator.scale+') '">
  <line
  :x1=0
  :y1=0
  :x2=0
  :y2=-distance
  style="stroke:gray;stroke-width:1"
  />

  <g transform="translate(0, -11) rotate(45)">
    <svg xmlns="http://www.w3.org/2000/svg" @mousedown="dragStart" @click="select" class="translate">
      <circle class="stroked" fill="white" stroke="blue" cx="8" cy="8" r="6"/>
      <path class="stroked" fill="none" stroke="blue" d="M 8 0 L 8 6.5"/>
      <path class="stroked" fill="none" stroke="blue" d="M 0 8 L 6.5 8"/>
      <path class="stroked" fill="none" stroke="blue" d="M 8 9.5 L 8 16"/>
      <path class="stroked" fill="none" stroke="blue" d="M 9.5 8 L 16 8"/>
    </svg>
  </g>

  <circle
      class="rotate"
      r="4"
      fill="green"
      :cx=0
      :cy=-distance
      @mousedown="rotateStart"
  />
  </g>

</template>

<script setup>

import {computed, reactive, ref} from 'vue'
import {roomTranslatorStore} from "../stores/roomTranslator";
const translator = roomTranslatorStore()

const props = defineProps({
  x: Number,
  y: Number,

  rotation: Number,

  relevantRoomElement: Array,
})

const distance = 20
const handle = computed(() => {
  const rot = (props.rotation - 90) * Math.PI / 180;

  return {
    x: props.x + distance * Math.cos(rot),
    y: props.y + distance * Math.sin(rot),
  }
})

function dragStart(e) {
  // Needs to be converted to world coordinates
  translator.offset = [
      props.x*translator.scale - e.clientX,
      props.y*translator.scale - e.clientY,
  ]

  translator.firstClick = [
    e.clientX,
    e.clientY,
  ]
  translator.translatingRoom = props.relevantRoomElement
}

function rotateStart(e) {
  const rot = (props.rotation - 90) * Math.PI / 180;

  translator.offset = [
    props.x + distance * Math.cos(rot) - e.clientX,
    props.y + distance * Math.sin(rot) - e.clientY
  ]

  translator.rotatingRoom = props.relevantRoomElement
}

function select() {
  translator.selectedRoom = props.relevantRoomElement
}

</script>

<style scoped>

svg.translate {
  cursor: move;
}

circle.rotate:hover {
  cursor: crosshair;
}

</style>
