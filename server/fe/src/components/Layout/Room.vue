<template>
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
                  // coord.x -=topLeft.x
                  // coord.y -=topLeft.y
                  room.outline.splice(i+1%room.outline.length, 0, coord)
                }"
      />
      <DragNg v-if=translating v-for="(coord, i) in room.outline"
              @dragStart="(e: DragStartEvent) => $emit('dragStart', e)">
        <g :transform="'translate('+(coord.x)+','+(coord.y)+')'">
          <g transform="translate(0, -11) rotate(45)">
            <svg xmlns="http://www.w3.org/2000/svg" class="translate"
                 @click.right.prevent="room.outline.splice(i%room.outline.length, 1)">
              <circle class="stroked" fill="white" stroke="blue" cx="8" cy="8" r="6"/>
              <path class="stroked" fill="none" stroke="blue" d="M 8 0 L 8 6.5"/>
              <path class="stroked" fill="none" stroke="blue" d="M 0 8 L 6.5 8"/>
              <path class="stroked" fill="none" stroke="blue" d="M 8 9.5 L 8 16"/>
              <path class="stroked" fill="none" stroke="blue" d="M 9.5 8 L 16 8"/>
            </svg>
          </g>
        </g>
      </DragNg>

      <Sequence
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

          @dragStart="(e: DragStartEvent) => $emit('dragStart', e)"
          @maybeShow="(e: DragStartEvent) => $emit('maybeShow', e)"
      />
</template>

<style scoped>
</style>

<script setup lang="ts">

import Sequence from "./Sequence.vue";
import DragNg from "../../views/DragNg.vue";
import {
  Coordinate, DragStartEvent,
  RoomInterface,
  SequenceAxis,
  SequenceDirection,
  SequenceType
} from "../../types.ts";
import {teamareaStore} from "../../stores/teamarea";

const settings = teamareaStore()

defineEmits(['dragStart', 'maybeShow', 'maybeSelect'])

const room = withDefaults(defineProps<RoomInterface & {
  translating?: boolean
}>(), {
  translating: false,
});

function toInnerCoordinates(e: MouseEvent): Coordinate {
  return {x: 0, y: 0}
}

</script>
