<template>
  <Sequence
      v-if="repeats.length-1 > atRepeats"
      v-for="i in num"

      :repeats=repeats
      :x=xPos(i)
      :y=yPos(i)
      :rotation=rot(i)

      :atRepeats=atRepeats+1
      :at-num="i"
      :relevant-room-element=relevantRoomElement
      v-bind=repeats[atRepeats+1]
  />
  <g v-else v-for="i in num">
    <TeamTable
        :x=xPos(i)??0
        :y=yPos(i)??0
        :rotation=rot(i)??0

        :team-id="''+(140+i)"
        :relevant-room-element=relevantRoomElement
    />
    <Draggable
        v-if="(i===1 && atRepeats<=repeats.length-1 && atNum===1)"
        :x=xPos(i)??0
        :y=yPos(i)??0
        :rotation=rot(i)??0
        :relevant-room-element=relevantRoomElement
    />
  </g>

</template>

<script setup>

import TeamTable from "@/components/Layout/TeamTable.vue";
import Draggable from "../../views/Draggable.vue";
import {computed} from "vue";

const SequenceType = {
  // Area: "area",
  Line: "Line",
  Circle: "Circle",
};

const props = defineProps({
  'type': {type: String, required: false, default: "Line"},
  'x': {type: Number, required: false, default: 0},
  'y': {type: Number, required: false, default: 0},
  'rotation': {type: Number, required: false, default: 0},
  'radius': {type: Number, required: false, default: 100},
  'num': {type: Number, required: false, default: 1},
  'atNum': {type: Number, required: false, default: 1},
  'axis': {type: Boolean, required: false, default: true},
  'dir': {type: Boolean, required: false, default: true},
  'separation': {type: Number, required: false, default: 50},
  'repeats': {type: Array, required: true},
  'atRepeats': {type: Number, required: true, default: 0},
  'equivalentSpaced': {type: Boolean, required: false, default: true},
  'relevantRoomElement': {type: Array, required: true},
})

function xPos(i) {
  i--;
  if (props.type === SequenceType.Line) {
    return Math.round(props.x + dirVec.value.x * i);
  } else {
    if (i === 0) {
      return props.x;
    }

    const base = props.x + distVec(props.rotation, props.radius, false, props.dir).x
    const rad = (props.rotation + dirInt.value * 90 + axisInt.value * trueSeparation.value * i) * Math.PI / 180;
    const offset = Math.cos(rad) * props.radius;

    return base + offset;
  }
}

function yPos(i) {
  i--;
  if (props.type === SequenceType.Line) {
    return Math.round(props.y + dirVec.value.y * i);
  } else {
    if (i === 0) {
      return props.y;
    }
    const base = props.y + distVec(props.rotation, props.radius, false, props.dir).y

    const rad = (props.rotation + dirInt.value * 90 + axisInt.value * trueSeparation.value * i) * Math.PI / 180;
    const offset = Math.sin(rad) * props.radius;

    return base + offset;
  }
}

function rot(i) {
  if (props.type === SequenceType.Line) {
    return props.rotation
  }

  return props.rotation + axisInt.value * trueSeparation.value * (i-1)
}

function distVec(rotation, sep, axis, dir) {
  let offset = 90;
  if (axis) {
    offset = 0;
  }

  const rad = (rotation + offset) * Math.PI / 180;
  let x = Math.cos(rad)
  let y = Math.sin(rad)

  if (dir) {
    x = -x;
    y = -y;
  }

  const mag = Math.sqrt(x * x + y * y)

  return {x: sep * x / mag, y: sep * y / mag, rot: 0};
}

const dirVec = computed(() => distVec(props.rotation, trueSeparation.value, props.axis, props.dir))
const axisInt = computed(() => props.axis ? 1 : -1);
const dirInt = computed(() => props.dir ? 1 : -1);
const trueSeparation = computed(() => (props.type !== SequenceType.Circle || !props.equivalentSpaced)
    ? props.separation : 360 / Math.max(1, props.num))

</script>

<style scoped>
</style>
