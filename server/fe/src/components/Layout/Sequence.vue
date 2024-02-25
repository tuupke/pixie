<template>
  <Sequence
      v-if="el.repeats.length -1 > atRepeats"

      v-bind="el.repeats[atRepeats+1]"
      v-for="i in num"

      :x=xPos(i)
      :y=yPos(i)
      :rotation=rot(i)

      :atRepeats=atRepeats+1
      :el=el
      :room=room
      @dragStart="(e: DragStartEvent) => $emit('dragStart', e)"
      @maybeShow="(e: DragStartEvent) => $emit('maybeShow', e)"
  />
  <g v-else v-for="i in num">
    <DragNg v-if="translating"
        @dragStart="(e: MouseEvent) => $emit('dragStart', {coord: el.base, event: e})"
        @maybeShow="(e: MouseEvent) => $emit('maybeShow', {el: el, room: room, event: e})"
    >
      <TeamTable
          :x=xPos(i)??0
          :y=yPos(i)??0
          :rotation=rot(i)??0

          :team-id="''+(140+i)"
      />
    </DragNg>
    <TeamTable
        v-else
        :x=xPos(i)??0
        :y=yPos(i)??0
        :rotation=rot(i)??0

        :team-id="''+(140+i)"
    />
  </g>

</template>

<script setup lang="ts">

import TeamTable from "./TeamTable.vue";
import {computed} from "vue";
import {
  Coordinate,
  DragStartEvent, ElementInterface, RoomInterface,
  RotationCoordinate,
  SequenceAxis,
  SequenceDirection,
  SequenceInterface,
  SequenceType
} from "../../types.ts";
import DragNg from "../../views/DragNg.vue";

interface SequenceLocal {
  atRepeats: number
  el: ElementInterface,
  room: RoomInterface,
  translating?: boolean,
}

defineEmits(['dragStart', 'maybeShow', 'maybeSelect'])
const props = withDefaults(defineProps<RotationCoordinate & SequenceInterface & SequenceLocal
>(), {
  type: SequenceType.Line,
  x: 0,
  y: 0,
  rotation: 0,
  radius: 100,
  num: 1,
  axis: SequenceAxis.Horizontal,
  dir: SequenceDirection.Positive,
  separation: 50,
  atRepeats: 0,
  equivalentSpaced: true,
  translating: false,
});

function xPos(i: number): number {
  i--;
  if (props.type === SequenceType.Line) {
    return Math.round(props.x + dirVec.value.x * i);
  } else {
    if (i === 0) {
      return props.x;
    }

    const base = props.x + distVec(props.rotation, props.radius, SequenceAxis.Vertical, props.dir).x
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
    const base = props.y + distVec(props.rotation, props.radius, SequenceAxis.Vertical, props.dir).y

    const rad = (props.rotation + dirInt.value * 90 + axisInt.value * trueSeparation.value * i) * Math.PI / 180;
    const offset = Math.sin(rad) * props.radius;

    return base + offset;
  }
}

function rot(i: number): number {
  if (props.type === SequenceType.Line) {
    return props.rotation
  }

  return props.rotation + axisInt.value * trueSeparation.value * (i - 1)
}

function distVec(rotation: number, sep: number, axis: SequenceAxis, dir: SequenceDirection): Coordinate {
  let offset = -90;
  if (axis == SequenceAxis.Horizontal) {
    offset = 0;
  }

  const rad = (rotation + offset) * Math.PI / 180;
  let x = Math.cos(rad)
  let y = Math.sin(rad)

  if (dir === SequenceDirection.Positive) {
    x = -x;
    y = -y;
  }

  const mag = Math.sqrt(x * x + y * y)
  return {x: sep * x / mag, y: sep * y / mag};
}

const dirVec = computed(() => distVec(props.rotation, trueSeparation.value, props.axis, props.dir))
const axisInt = computed(() => props.axis == SequenceAxis.Horizontal ? 1 : -1);
const dirInt = computed(() => props.dir == SequenceDirection.Negative ? 1 : -1);
const trueSeparation = computed(() => (props.type !== SequenceType.Circle || !props.equivalentSpaced)
    ? props.separation : 360 / Math.max(1, props.num))

</script>

<style scoped>
</style>
