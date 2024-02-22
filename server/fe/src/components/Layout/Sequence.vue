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

<script setup lang="ts">

import TeamTable from "./TeamTable.vue";
import Draggable from "../../views/Draggable.vue";
import {computed} from "vue";

const SequenceType = {
  // Area: "area",
  Line: "Line",
  Circle: "Circle",
};

const props = withDefaults(defineProps<{
  type: string
  x: number
  y: number
  rotation: number
  radius: number
  num: number
  atNum: number
  axis: boolean
  dir: boolean
  separation: number
  repeats: any[]
  atRepeats: number
  equivalentSpaced: boolean
  relevantRoomElement: any[]
}>(), {
  type: "Line",
  x: 0,
  y: 0,
  rotation: 0,
  radius: 100,
  num: 1,
  atNum: 1,
  axis: true,
  dir: true,
  separation: 50,
  atRepeats: 0,
  equivalentSpaced: true
});

function xPos(i: number): number {
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

function yPos(i: number): number {
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

function rot(i: number): number {
  if (props.type === SequenceType.Line) {
    return props.rotation
  }

  return props.rotation + axisInt.value * trueSeparation.value * (i-1)
}

function distVec(rotation: number, sep: number, axis: boolean, dir: boolean): {
  x: number
  y: number
  rot: number
} {
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
