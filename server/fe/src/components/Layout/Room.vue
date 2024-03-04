<template>
  <DragNg
      v-for="(el) in room.elements"
      v-if="translating"
      :coord=el.base
      :transform="true"
      @dragStart="(e: DragStartEvent) => $emit('dragStart', e)"
      @maybeShow="(e: DragStartEvent) => $emit('maybeShow', e)"
      @rotateStart="(e: RotationStartEvent) => $emit('rotateStart', e)"
  >
    <Sequence
        :x=0
        :y=0
        :rotation=0
        :dir=SequenceDirection.Positive
        :separation=1337
        :radius=1337
        :equivalentSpaced=false
        :axis=SequenceAxis.Horizontal
        :num=1
        :type=SequenceType.Line
        :atRepeats=-1
        :el=el
        :room=room
        :translating="false"
    />
  </DragNg>
  <Sequence
      v-else
      v-for="(el) in room.elements"

      :x=el.base.x
      :y=el.base.y
      :rotation=el.base.rotation
      :dir=SequenceDirection.Positive
      :separation=1337
      :radius=1337
      :equivalentSpaced=false
      :axis=SequenceAxis.Horizontal
      :num=1
      :type=SequenceType.Line
      :atRepeats=-1
      :el=el
      :room=room
      :translating="false"
  />
  <line
      v-for="(coord, i) in room.outline"
      stroke="black"
      :stroke-width="settings.strokeWidth"
      :x1=coord.x
      :y1=coord.y
      :x2=room.outline[(i+1)%room.outline.length].x
      :y2=room.outline[(i+1)%room.outline.length].y
      @click.left.stop="(e: MouseEvent) => {
                  const coord = toInnerCoordinates(e)
                  console.log(coord, e)
                  room.outline.splice(i+1%room.outline.length, 0, coord)
                }"
  />
  <DragNg v-if=translating v-for="(coord, i) in room.outline"
          :coord=coord
          :outline="false"
          :transform="true"
          @dragStart="(e: DragStartEvent) => $emit('dragStart', e)"
          @rotateStart="(e: RotationStartEvent) => $emit('rotateStart', e)"
  >
    <Crosshairs/>
  </DragNg>

</template>

<style scoped>
</style>

<script setup lang="ts">

import Sequence from "./Sequence.vue";
import Crosshairs from "./Crosshairs.vue";
import {
  CoordinateInterface,
  DragStartEvent,
  RoomInterface,
  RotationStartEvent,
  SequenceAxis,
  SequenceDirection,
  SequenceType
} from "../../types.ts";
import {teamareaStore} from "../../stores/teamarea";
import DragNg from "../../views/DragNg.vue";
import {inject, provide} from "vue";

const settings = teamareaStore()

defineEmits(['dragStart', 'maybeShow', 'maybeSelect', 'rotateStart'])

const room = withDefaults(defineProps<RoomInterface & {
  translating?: boolean
}>(), {
  translating: false,
});

const toInnerCoordinates = inject('toInnerCoordinates')
//
// function toInnerCoordinates(e: MouseEvent): Coordinate {
//   let v = e.target
//   while (v && v.nodeName !== "svg") {
//     v = v.parentNode
//   }
//
//   if (!v) {
//     window.alert("Not found")
//     return {x: 0, y: 0}
//   }
//
//   const pt = v.createSVGPoint();
//   pt.x = e.clientX //size.x+size.width/2;
//   pt.y = e.clientY //size.y+size.height/2;
//
//   return {x: 0, y: 0}
// }

</script>
